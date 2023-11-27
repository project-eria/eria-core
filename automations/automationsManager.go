package automations

import (
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/producer"
)

type Automation struct {
	Ref    string  `yaml:"ref"` // thing ref
	Name   string  `yaml:"name" required:"true"`
	Action string  `yaml:"action" required:"true"`
	Groups []Group `yaml:"groups" required:"true"`
}

type Group struct {
	//	Name       string   `yaml:"name" required:"true"`
	Scheduled  string   `yaml:"scheduled"`
	Conditions []string `yaml:"conditions" required:"true"`
}

var (
// _deviceThings  = make(map[string]*consumer.ConsumedThing)
)

/**
 * Start the automations manager, get and schedule jobs
 * @param now
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 */
func Start(now time.Time, automations []Automation, contextsThing consumer.ConsumedThing, exposedThings map[string]producer.ExposedThing, consumedThings map[string]consumer.ConsumedThing) {
	jobs := getJobs(now, automations, contextsThing, exposedThings, consumedThings)
	// run/schedule all jobs
	resetCronScheduler(now.Location())
	scheduleJobs(jobs)
}
