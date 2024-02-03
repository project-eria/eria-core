package automations

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	eriaconsumer "github.com/project-eria/eria-core/consumer"
	zlog "github.com/rs/zerolog/log"
)

/**
 * Automation based on Time
 */

type scheduleAtHour struct {
	cronJob gocron.Job
	//	fixedHour         *time.Time
	timeThing         eriaconsumer.Thing
	timeThingProperty string
	min               *time.Time
	max               *time.Time
	scheduledTime     *time.Time
	observerRef       uint16
}

/**
 * Handle time scheduling (hour)
 * `hour|hour|<time>`
 *	- fixed hour:
 *		- `at|hour|<hour in 15:04 format>`
 *		- `at|hour|<hour in 15:04:05 format>`
 *	- thing hour property: `at|hour|<thing device>:<property>`
 * 		- min option: `at|hour|<thing device>:<property>|min=<hour in 15:04:05 format>`
 *		- max option: `at|hour|<thing device>:<property>|max=<hour in 15:04:05 format>`
 * @param scheduleArray
 */
func NewScheduleAtHour(scheduleArray []string) (*scheduleAtHour, error) {
	if len(scheduleArray) < 3 || len(scheduleArray) > 5 {
		return nil, errors.New("invalid at scheduling length")
	}

	s := &scheduleAtHour{}

	// Test if hour in 15:04/15:04:05 format`
	t, err := HourToTime(scheduleArray[2])
	if err == nil {
		// Correct time format
		s.scheduledTime = &t
		return s, nil
	}

	// Test if min/max is set
	if len(scheduleArray) > 3 {
		for i := 3; i < len(scheduleArray); i++ {
			if strings.HasPrefix(scheduleArray[i], "min=") {
				t, err := HourToTime(scheduleArray[i][4:])
				if err != nil {
					return nil, errors.New("invalid min time: " + err.Error())
				}
				s.min = &t
			} else if strings.HasPrefix(scheduleArray[i], "max=") {
				t, err := HourToTime(scheduleArray[i][4:])
				if err != nil {
					return nil, errors.New("invalid max time: " + err.Error())
				}
				s.max = &t
			}
		}
	}

	// Test thing hour property: `at|hour|<thing device>:<property>`
	timeThingArray := strings.Split(scheduleArray[2], ":")
	if len(timeThingArray) != 2 {
		return nil, errors.New("invalid time")
	}
	// Does the time thing exist?
	timeThing := _consumer.ThingFromTag(timeThingArray[0])
	if timeThing == nil {
		return nil, errors.New("time thing doesn't exist")
	}
	s.timeThing = timeThing

	// Get the property value
	timeThingValue, err := timeThing.Property(timeThingArray[1]).Value()
	if err != nil {
		return nil, errors.New("time thing property not available: " + err.Error())
	}
	s.timeThingProperty = timeThingArray[1]

	// Does the time thing value is valid?
	// TODO the format can be customised?
	_, err = DateToHour(timeThingValue.(string))
	if err != nil {
		return nil, errors.New("invalid time thing value: " + err.Error())
	}

	return s, nil
}

func (s *scheduleAtHour) start(action ActionRunner) error {
	var automationName string
	if a, ok := action.(*Action); ok {
		// If not in mock
		automationName = a.AutomationName
	}

	if action == nil {
		return errors.New("missing action")
	}
	if s.scheduledTime == nil {
		return errors.New("missing scheduled hour")
	}

	cronJob, err := s.scheduleTask(automationName, action)
	if err != nil {
		return err
	}
	s.cronJob = cronJob

	if s.timeThingProperty != "" {
		// Observe the property, in case of the hour property changes
		s.observerRef, _ = s.timeThing.Property(s.timeThingProperty).Observe(func(value interface{}, err error) {
			if err == nil {
				hour, err := getPropertyHour(value.(string), s.min, s.max)
				if err != nil {
					zlog.Error().Err(err).Msg("[automations:scheduleAtHour:start] Invalid time thing value")
					return
				}
				zlog.Info().Str("Automation", automationName).Msgf("[automations:scheduleAtHour] Schedule hour changed, rescheduling %s -> %s", s.scheduledTime.Format("15:04:05"), hour.Format("15:04:05"))
				s.scheduledTime = hour

				// Cancelling the previous job
				err = _cronScheduler.RemoveJob(s.cronJob.ID())
				if err != nil {
					zlog.Error().Err(err).Msg("[automations:scheduleAtHour:start] Failed to remove cron job")
					return
				}

				// Re-Scheduling the job
				cronJob, err := s.scheduleTask(automationName, action)
				if err != nil {
					zlog.Error().Err(err).Msg("[automations:scheduleAtHour:start] Failed to schedule new cron job")
					return
				}
				s.cronJob = cronJob
			}
		})
	}

	return nil
}

func (s *scheduleAtHour) scheduleTask(automationName string, action ActionRunner) (gocron.Job, error) {
	return _cronScheduler.NewJob(
		gocron.DailyJob(1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(s.scheduledTime.Hour()), uint(s.scheduledTime.Minute()), uint(s.scheduledTime.Second())),
			),
		),
		gocron.NewTask(
			func() {
				zlog.Info().Str("Automation", automationName).Msg("[automations:scheduleAtHour] Running scheduled job")
				err := action.run()
				if err != nil {
					zlog.Error().Err(err).Str("Automation", automationName).Msg("[automations:scheduleAtHour:start] Failed to run scheduled job")
				}
			},
		),
		gocron.WithTags("core", "automation", "atHour"),
	)
}
func (s *scheduleAtHour) job() error {
	if s.timeThingProperty != "" {
		timeThingValue, err := s.timeThing.Property(s.timeThingProperty).Value()
		if err != nil {
			return errors.New("time thing property not available: " + err.Error())
		}
		t, err := getPropertyHour(timeThingValue.(string), s.min, s.max)
		if err != nil {
			return err
		}
		s.scheduledTime = t
	} // else: Fixed hour, time already set
	return nil
}

func getPropertyHour(timeThingValue string, min *time.Time, max *time.Time) (*time.Time, error) {
	// Convert the thing date to a time, without the date/location, for comparison
	timeValue, err := DateToHour(timeThingValue)
	// TODO the format can be customised?
	if err != nil {
		return nil, errors.New("invalid time thing value: " + err.Error())
	}

	if min != nil && timeValue.Before(*min) {
		return min, nil
	}
	if max != nil && timeValue.Before(*max) {
		return max, nil
	}
	return &timeValue, nil
}

func (s *scheduleAtHour) cancel() {
	err := _cronScheduler.RemoveJob(s.cronJob.ID())
	if err != nil {
		zlog.Error().Err(err).Msg("[automations:scheduleAtHour:cancel] Failed to remove cron job")
	}
	s.cronJob = nil
	// TODO s.timeThing.Property(s.propertyHour).UnObserve(s.observerRef)
}

func (s *scheduleAtHour) equals(other Schedule) bool {
	if s == nil || other == nil {
		return false
	}
	s2, ok := other.(*scheduleAtHour)
	if !ok || s2.scheduledTime == nil || s.scheduledTime == nil {
		return false
	}
	return ok && s.scheduledTime.Equal(*s2.scheduledTime)
}

func (s *scheduleAtHour) string() string {
	return fmt.Sprintf("every day at %s", s.scheduledTime.Format("15:04:05"))
}
