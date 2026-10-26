package automations

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/gookit/goutil/arrutil"
	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

// AutomationConfig is the YAML declaration of an automation. An automation owns
// a set of things and one or more groups; each group declares an action and a
// first-match list of schedule entries deciding when (or whether) that action
// runs.
type AutomationConfig struct {
	Name   string        `yaml:"name" required:"true"`
	Things []string      `yaml:"things" required:"true"`
	Groups []groupConfig `yaml:"groups" required:"true"`
}

// groupConfig declares one action and the schedule branches that drive it. The
// first branch whose conditions all match wins ("first schedule condition that
// matches wins").
type groupConfig struct {
	Action    string                `yaml:"action" required:"true"`
	Schedules []scheduleEntryConfig `yaml:"schedules" required:"true"`
}

// scheduleEntryConfig is one branch of a group's first-match selection:
//   - conditions: all must match for this branch to be selected (empty = always,
//     i.e. the default branch);
//   - schedule: when the branch is selected, the action is scheduled accordingly.
//     An empty schedule disables the group for that case ("no schedule means
//     this automation is disabled").
type scheduleEntryConfig struct {
	Conditions []string `yaml:"conditions"`
	Schedule   string   `yaml:"schedule"`
}

// --- runtime model ---

// scheduleEntry is the prepared form of scheduleEntryConfig.
type scheduleEntry struct {
	conditions []Condition
	schedule   Schedule // nil => disabled when this entry is the first match
}

// actionGroup is the prepared form of groupConfig: one action selected through
// a first-match list of schedule entries.
type actionGroup struct {
	action  ActionRunner
	entries []scheduleEntry
	active  bool // a scheduled branch matched at the last planning (guarded by _mu)
}

// scheduledAction binds a selected schedule to the action it runs.
type scheduledAction struct {
	schedule Schedule
	action   ActionRunner
}

type Automation struct {
	name          string
	groups        []actionGroup
	jobs          []scheduledAction // the currently planned set (one per active group)
	lastScheduled time.Time
	status        string
}

var (
	_exposedThings       map[string]producer.ExposedThing
	_location            *time.Location
	_contextsAutomations = make(map[string][]*Automation) // list of automations by context

	_activeContexts = []string{} // The currently active contexts
	_consumer       eriaconsumer.Consumer
	_cronScheduler  gocron.Scheduler

	// _automations is the registry of every prepared automation, used by
	// Scheduled() to introspect the live planning state. _mu guards it together
	// with each automation's jobs slice, which is re-planned from event
	// goroutines (context changes, hour re-scheduling).
	_automations []*Automation
	_mu          sync.Mutex
)

/**
 * Start the automations manager, get and schedule jobs
 * @param location the location of the time
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 */
func Start(location *time.Location, automations []AutomationConfig, exposedThings map[string]producer.ExposedThing, consumer *eriaconsumer.EriaConsumer, cronScheduler gocron.Scheduler) {
	_exposedThings = exposedThings
	_consumer = consumer
	_location = location
	_cronScheduler = cronScheduler

	contextsThing := _consumer.ThingFromTag("contexts")
	if contextsThing == nil {
		zlog.Warn().Msg("[automations:Start] Contexts thing not found, contexts will not be used")
	} else {
		rawContexts, err := contextsThing.Property("actives").Value()
		if err == nil {
			// Save the current active contexts
			_activeContexts = arrutil.MustToStrings(rawContexts)
			// Monitor context changes
			contextsThing.Property("actives").Observe(func(value interface{}, err error) {
				if err == nil {
					current := arrutil.MustToStrings(value)
					diff := arrutil.Diff(_activeContexts, current, arrutil.StringEqualsComparer)

					// Save the current active contexts, for the conditions
					_activeContexts = current

					for _, context := range diff {
						zlog.Info().Str("context", context).Msg("[automations:Start] Context changed, re-schedule")
						if _, found := _contextsAutomations[context]; found { // If we have automations for this context
							for _, automation := range _contextsAutomations[context] {
								// Re-schedule jobs
								automation.scheduleJob(time.Now().In(_location))
							}
						}
					}
				}
			})
		} else {
			zlog.Warn().Msg("[automations:Start] Contexts thing not rechable, contexts will not be used")
		}
	}
	scheduleJobs(automations)
}
