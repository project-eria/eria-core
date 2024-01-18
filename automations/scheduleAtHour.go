package automations

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/project-eria/go-wot/consumer"
	zlog "github.com/rs/zerolog/log"
)

/**
 * Automation based on Time
 * The action is planned using gocron, and re-builded every day at 3:00 by default
 * (because of the summer/winter times changes), but that can be changed
 */

// TODO - re-build at specific time
// TODO - thing hour

type scheduleAtHour struct {
	cronJob       *gocron.Job
	fixedHour     string // Keep as string for futur time comparison
	timeThing     consumer.ConsumedThing
	propertyHour  string
	min           string // Keep as string for futur time comparison
	max           string // Keep as string for futur time comparison
	scheduledHour string
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

	// Test fixed hour: `at|hour|<hour in 15:04 format>`
	_, err := time.Parse("15:04", scheduleArray[2])
	if err == nil {
		// Correct time format
		s.fixedHour = scheduleArray[2]
		return s, nil
	}
	// Test fixed hour: `at|hour|<hour in 15:04:05 format>`
	_, err = time.Parse("15:04:05", scheduleArray[2])
	if err == nil {
		// Correct time format
		s.fixedHour = scheduleArray[2]
		return s, nil // Remove the seconds part
	}

	// Test if min/max is set
	if len(scheduleArray) > 3 {
		for i := 3; i < len(scheduleArray); i++ {
			if strings.HasPrefix(scheduleArray[i], "min=") {
				_, err := time.Parse("15:04", scheduleArray[i][4:])
				if err != nil {
					return nil, errors.New("invalid min time: " + err.Error())
				}
				s.min = scheduleArray[i][4:]
			} else if strings.HasPrefix(scheduleArray[i], "max=") {
				_, err := time.Parse("15:04", scheduleArray[i][4:])
				if err != nil {
					return nil, errors.New("invalid max time: " + err.Error())
				}
				s.max = scheduleArray[i][4:]
			}
		}
	}

	// Test thing hour property: `at|hour|<thing device>:<property>`
	timeThingArray := strings.Split(scheduleArray[2], ":")
	if len(timeThingArray) != 2 {
		return nil, errors.New("invalid time")
	}
	// Does the time thing exist?
	timeThing, exists := _consumedThings[timeThingArray[0]]
	if !exists {
		return nil, errors.New("time thing doesn't exist")
	}
	s.timeThing = timeThing

	// Get the property value
	timeThingValue, err := timeThing.ReadProperty(timeThingArray[1], nil)
	if err != nil {
		return nil, errors.New("time thing property not available: " + err.Error())
	}
	s.propertyHour = timeThingArray[1]

	// Does the time thing value is valid?
	// TODO the format can be customised?
	_, err = time.Parse(time.RFC3339, timeThingValue.(string))
	if err != nil {
		return nil, errors.New("invalid time thing value: " + err.Error())
	}

	return s, nil
}

func (s *scheduleAtHour) start(action ActionRunner) error {
	if action == nil {
		return errors.New("missing action")
	}
	if s.scheduledHour == "" {
		return errors.New("missing scheduled hour")
	}
	cronJob, err := _cronScheduler.Every(1).Day().At(s.scheduledHour).Tag("atHour").Do(func() {
		zlog.Info().Str("Automation", action.(*Action).AutomationName).Msg("Running scheduled job")
		err := action.run()
		if err != nil {
			zlog.Error().Err(err).Str("Automation", action.(*Action).AutomationName).Msg("Failed to run scheduled job")
		}
	})
	if err != nil {
		return err
	}
	s.cronJob = cronJob
	return nil
}

func (s *scheduleAtHour) job() error {
	var scheduledHour string

	if s.fixedHour != "" {
		scheduledHour = s.fixedHour
	} else {
		// Get the property value
		timeThingValue, err := s.timeThing.ReadProperty(s.propertyHour, nil)
		if err != nil {
			return errors.New("time thing property not available: " + err.Error())
		}

		// Convert the thing date to a time, without the date/location, for comparison
		timeValue, err := DateToHour(timeThingValue.(string))
		// TODO the format can be customised?
		if err != nil {
			return errors.New("invalid time thing value: " + err.Error())
		}
		scheduledHour = timeValue.Format("15:04")

		if yes, _ := HourIsBefore(scheduledHour, s.min); yes {
			scheduledHour = s.min
		} else if yes, _ := HourIsAfter(scheduledHour, s.max); yes {
			scheduledHour = s.max
		}
	}
	s.scheduledHour = scheduledHour
	return nil
}

func (s *scheduleAtHour) cancel() {
	err := _cronScheduler.RemoveByID(s.cronJob)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed to remove cron job")
	}
}

func (s *scheduleAtHour) equals(other Schedule) bool {
	if other == nil {
		return false
	}
	s2, ok := other.(*scheduleAtHour)
	return ok && s.scheduledHour == s2.scheduledHour
}

func (s *scheduleAtHour) string() string {
	return fmt.Sprintf("every day at %s", s.scheduledHour)
}
