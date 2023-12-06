package automations

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type AutomationConfig struct {
	Ref    string `yaml:"ref"` // thing ref
	Name   string `yaml:"name" required:"true"`
	Action string `yaml:"action" required:"true"`
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

type automation struct {
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
	_automations    = make(map[string]*automation)
)

/**
 * Start the automations manager, get and schedule jobs
 * @param location the location of the time
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 */
func Start(location *time.Location, automations []AutomationConfig, contextsThingRef string, exposedThings map[string]producer.ExposedThing, consumedThings map[string]consumer.ConsumedThing) {
	_exposedThings = exposedThings
	_consumedThings = consumedThings
	_location = location
	if _, found := consumedThings[contextsThingRef]; !found {
		zlog.Warn().Str("thing", contextsThingRef).Msg("[automations:Start] Contexts thing not found, contexts will not be used")
	} else {
		_contextsThing = consumedThings[contextsThingRef]
	}
	initCronScheduler()
	scheduleJobs(automations)
	_cronScheduler.StartAsync()
}

var (
	_cronScheduler *gocron.Scheduler
)

func initCronScheduler() {
	if _cronScheduler != nil {
		_cronScheduler.Clear()
	}

	_cronScheduler = gocron.NewScheduler(_location)
	_cronScheduler.TagsUnique()
}
