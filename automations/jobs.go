package automations

import (
	"errors"
	"strings"
	"time"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/producer"
	zlog "github.com/rs/zerolog/log"
)

type Job struct {
	Key           string
	Name          string
	Action        Action
	ScheduledType string
	Scheduled     string
}

/**
 * Get all jobs that should be executed
 * @param now
 * @param automations the automations list
 * @param contextsThing the thing service to retrive contexts
 * @return the list of jobs
 */
func getJobs(now time.Time, automations []Automation, contextsThing consumer.ConsumedThing, exposedThings map[string]producer.ExposedThing, consumedThings map[string]consumer.ConsumedThing) []Job {
	jobs := []Job{}

	// 	// Observe context changes to reset schedule
	// 	thing := _contextThings[contextDetails.Url]
	// 	thing.ObserveProperty(contextDetails.Property, func(value interface{}, err error) {
	// 		zlog.Info().Str("context", key).Msg("Context changed, re-schedule")
	// 		// TODO
	// 	})

	for _, automation := range automations {
		if exposedThing, ok := exposedThings[automation.Ref]; ok {
			job := getJob(now, automation.Name, automation.Action, automation.Groups, contextsThing, exposedThing, consumedThings)
			if job.ScheduledType != "" {
				jobs = append(jobs, job)
			}
		} else {
			zlog.Error().Str("automation", automation.Name).Str("things", automation.Ref).Msg("[automations:getJobs] Exposed thing not found, skipping...")
		}
	}
	return jobs
}

/**
 * Get the job that matches the conditions
 * @param now
 * @param automation the automation details
 * @param contextsThing the thing service to retrive contexts
 * @return the final job
 */
func getJob(now time.Time, name string, actionStr string, groups []Group, contextsThing consumer.ConsumedThing, exposedThing producer.ExposedThing, consumedThings map[string]consumer.ConsumedThing) Job {
	if exposedThing == nil {
		// Handle the case when exposedThing is nil
		zlog.Fatal().Msg("[automations:getJob] ExposedThing is nil")
	}
	zlog.Info().Str("automation", name).Msg("[automations:getJob] Adding automation")
	for _, group := range groups {
		// The first group that matches all conditions, wins
		var allTrue = true
		// Check all conditions
		for _, condition := range group.Conditions {
			conditionArray := strings.Split(condition, "|")
			var err error
			var ok bool
			switch conditionArray[0] {
			case "context":
				ok, err = contextCondition(conditionArray, contextsThing)
			case "time":
				ok, err = timeCondition(conditionArray, now)
			case "property":
				ok, err = propertyCondition(conditionArray)
			default:
				err = errors.New("invalid condition type")
			}
			if err != nil {
				zlog.Error().Err(err).Str("automation", name).Strs("condition", conditionArray).Msg("[automations:getJob]")
				// Skip this group
				allTrue = false
				break
			} else {
				allTrue = allTrue && ok
			}
		}
		if allTrue {
			// If all conditions are true, we schedule/plan the action
			scheduleArray := strings.Split(group.Scheduled, "|")
			var (
				err          error
				schedule     string
				scheduleType string
			)
			// schedule/run job
			switch scheduleArray[0] {
			case "immediate":
				schedule, err = immediateSchedule(scheduleArray)
				scheduleType = "immediate"
			case "at":
				if scheduleArray[1] == "hour" {
					schedule, err = atHourSchedule(scheduleArray, consumedThings)
					scheduleType = "atHour"
				} else if scheduleArray[1] == "date" {
					//schedule, err = atDateSchedule(scheduleArray, consumedThings)
					scheduleType = "atDate"
				} else {
					err = errors.New("invalid at schedule type")
				}
			// case "every":
			// case: "in":
			default:
				err = errors.New("invalid schedule type")
			}

			if err != nil {
				zlog.Error().Err(err).Str("automation", name).Strs("schedule", scheduleArray).Msg("[automations:getJob]")
				return Job{} // Skip this automation
			}

			// Get action
			actionStr = strings.TrimSpace(actionStr)
			if actionStr == "" {
				zlog.Error().Str("automation", name).Msg("[automations:getJob] Empty action, skipping...")
				return Job{} // Skip this automation
			}
			actionArray := strings.Split(actionStr, "|")
			action, err := getAction(exposedThing, actionArray)
			if err != nil {
				zlog.Error().Err(err).Str("automation", name).Strs("action", actionArray).Msg("[automations:getJob]")
				return Job{} // Skip this automation
			}

			// Generate job with action
			job := Job{
				Name:          name,
				Action:        action,
				ScheduledType: scheduleType,
				Scheduled:     schedule,
			}

			return job // Break the groups loop, to next automation, if all conditions are true
		}
	}

	return Job{} // No matching conditions, return empty job
}
