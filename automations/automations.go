package automations

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gookit/goutil/arrutil"
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type AutomationConfig struct {
	Things []string `yaml:"things" required:"true"` // thing ref
	Name   string   `yaml:"name" required:"true"`
	Action string   `yaml:"action" required:"true"`
	Groups []struct {
		//	Name       string   `yaml:"name" required:"true"`
		Schedule   string   `yaml:"schedule" required:"true"`
		Conditions []string `yaml:"conditions"`
	} `yaml:"groups" required:"true"`
}

type group struct {
	schedule   Schedule
	conditions []Condition
}

type Automation struct {
	name          string
	groups        []group
	job           Schedule
	lastScheduled time.Time
	status        string
	action        Action
	exposedThing  producer.ExposedThing
}

var (
	_contextsThing  consumer.ConsumedThing
	_exposedThings  map[string]producer.ExposedThing
	_consumedThings map[string]consumer.ConsumedThing
	_location       *time.Location
	_automations    = make(map[string]*Automation)

	_contextsAutomations = make(map[string][]*Automation) // list of automations by context

	_activeContexts = []string{} // The currently active contexts

	_cronScheduler *gocron.Scheduler
)

/**
 * Start the automations manager, get and schedule jobs
 * @param location the location of the time
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 */
func Start(location *time.Location, automations []AutomationConfig, contextsThingRef string, exposedThings map[string]producer.ExposedThing, consumedThings map[string]consumer.ConsumedThing, cronScheduler *gocron.Scheduler) {
	_exposedThings = exposedThings
	_consumedThings = consumedThings
	_location = location
	_cronScheduler = cronScheduler
	if _, found := consumedThings[contextsThingRef]; !found {
		zlog.Warn().Str("thing", contextsThingRef).Msg("[automations:Start] Contexts thing not found, contexts will not be used")
	} else {
		_contextsThing = consumedThings[contextsThingRef]
		rawContexts, err := _contextsThing.ReadProperty("actives", nil)
		if err == nil {
			// Save the current active contexts
			_activeContexts = arrutil.MustToStrings(rawContexts)
			// Monitor context changes
			_contextsThing.ObserveProperty("actives", nil, func(value interface{}, err error) {
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
			zlog.Warn().Str("thing", contextsThingRef).Msg("[automations:Start] Contexts thing not rechable, contexts will not be used")
		}
	}
	scheduleJobs(automations)
}
