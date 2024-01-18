package automations

import (
	"errors"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

/**
 * creates and schedules automation jobs
 * @param now
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 * @param exposedThings the exposed things list
 * @param consumedThings the consumed things list
 */
func scheduleJobs(automations []AutomationConfig) {
	now := time.Now().In(_location)
	// 	// Observe context changes to reset schedule
	// 	thing := _contextThings[contextDetails.Url]
	// 	thing.ObserveProperty(contextDetails.Property, func(value interface{}, err error) {
	// 		zlog.Info().Str("context", key).Msg("Context changed, re-schedule")
	// 		// TODO
	// 	})

	for _, automationConfig := range automations {
		for _, ref := range automationConfig.Things {
			if exposedThing, ok := _exposedThings[ref]; ok && exposedThing != nil {
				automation, observables, err := getAutomation(automationConfig, exposedThing)
				if err == nil {
					automation.scheduleJob(now)
					_automations[ref] = automation
					// Set the observables, for monitoring
					for _, context := range observables.contexts {
						if _, found := _contextsAutomations[context]; !found {
							_contextsAutomations[context] = make([]*Automation, 0)
						}
						_contextsAutomations[context] = append(_contextsAutomations[context], automation)
					}
				} else {
					zlog.Warn().Err(err).Str("automation", automationConfig.Name).Str("thing", ref).Msg("[automations:scheduleJobs] Skipped for thing")
				}
			} else {
				zlog.Error().Str("automation", automationConfig.Name).Str("thing", ref).Msg("[automations:scheduleJobs] Exposed thing not found, skipping...")
			}
		}
	}
}

/**
 * Prepares the automation details
 * @param now
 * @param automation the automation details
 * @param exposedThing the exposed thing
 */
func getAutomation(automationConfig AutomationConfig, exposedThing producer.ExposedThing) (*Automation, *Observables, error) {
	zlog.Info().Str("automation", automationConfig.Name).Msg("[automations:getAutomation] Adding automation")

	// Prepare action
	zlog.Trace().Str("automation", automationConfig.Name).Str("action", automationConfig.Action).Msg("[automations:getAutomation] preparing action")
	action, err := getAction(exposedThing, automationConfig.Action)
	if err != nil {
		zlog.Error().Err(err).Str("automation", automationConfig.Name).Str("action", automationConfig.Action).Msg("[automations:getAutomation]")
		return nil, nil, err // Skip this automation
	}
	groups := make([]group, len(automationConfig.Groups))
	observables := &Observables{
		contexts: make([]string, 0),
	}
	for i, groupConfig := range automationConfig.Groups {
		zlog.Trace().Str("automation", automationConfig.Name).Int("group", i).Str("schedule", groupConfig.Schedule).Msg("[automations:getAutomation] preparing schedule")
		s, err := getSchedule(groupConfig.Schedule)
		if err != nil {
			zlog.Error().Err(err).Str("automation", automationConfig.Name).Int("group", i).Str("schedule", groupConfig.Schedule).Msg("[automations:getAutomation] Can't get schedule")
			return nil, nil, err // Skip this automation
		}
		zlog.Trace().Str("automation", automationConfig.Name).Int("group", i).Str("conditions", arrutil.JoinSlice(";", groupConfig.Conditions)).Msg("[automations:getAutomation] preparing conditions")
		c, o, err := getConditions(groupConfig.Conditions)
		if err != nil {
			zlog.Error().Err(err).Str("automation", automationConfig.Name).Int("group", i).Msg("[automations:getAutomation] Can't get conditions")
			return nil, observables, err // Skip this automation
		}
		if o != nil {
			for _, context := range o.contexts {
				// Add the context to the observables list (and remove duplicates)
				if !arrutil.StringsHas(observables.contexts, context) {
					observables.contexts = append(observables.contexts, context)
				}
			}
		}
		groups[i] = group{
			schedule:   s,
			conditions: c,
		}
	}
	return &Automation{
		name:         automationConfig.Name,
		exposedThing: exposedThing,
		action:       action,
		groups:       groups,
	}, observables, nil
}

func (automation *Automation) scheduleJob(now time.Time) {
	zlog.Trace().Str("automation", automation.name).Msg("[automations:scheduleJob] scheduling job")
	j, err := automation.getJob(now)
	if err != nil {
		zlog.Error().Err(err).Str("automation", automation.name).Interface("job", j).Msg("[automations:scheduleJob]")
		automation.status = err.Error()
		automation.job = nil
		return // Skip this automation
	}
	if automation.job != nil {
		if automation.job.equals(j) {
			return // No change
		}
		// Clean/Cancel previous job
		automation.job.cancel()
	}
	automation.job = j // Replace old job with the new one
	err = j.start(automation.action)
	automation.lastScheduled = now
	if err != nil {
		zlog.Error().Err(err).Str("automation", automation.name).Interface("job", j).Msg("[automations:scheduleJob]")
		automation.status = err.Error()
		return // Skip this automation
	}
	zlog.Info().Str("automation", automation.name).Msgf("[automations:scheduleJob] Job scheduled: %s", j.string())
	automation.status = "success"
}

/**
 * Get the job that matches the conditions
 * @param now
 * @param automation the automation details
 * @param contextsThing the thing service to retrive contexts
 * @param exposedThing the exposed thing
 * @param consumedThings the consumed things
 * @return the final job
 */
func (automation *Automation) getJob(now time.Time) (Schedule, error) {
	if len(automation.groups) == 0 {
		return nil, errors.New("missing conditions")
	}
	for _, group := range automation.groups {
		// The first group that matches all conditions, wins
		// Check the group conditions
		ok, err := checkConditions(group.conditions, now)
		if err != nil {
			zlog.Error().Err(err).Str("automation", automation.name).Msg("[automations:getJob]")
			// Skip this group
			continue
		} else if !ok {
			// If one condition is not true, we skip the group
			continue
		}
		// If all conditions are true, we return the association schedule
		err = group.schedule.job()
		if err != nil {
			zlog.Error().Err(err).Str("automation", automation.name).Msg("[automations:getJob]")
			// Skip this group
			continue
		}
		return group.schedule, nil
	}
	return nil, errors.New("no matching conditions")
}
