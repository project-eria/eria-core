package automations

import (
	"time"

	"github.com/go-co-op/gocron"
	zlog "github.com/rs/zerolog/log"
)

var (
	_cronScheduler *gocron.Scheduler
)

/**
 * Prepare all the tasks for running/scheduling, in global tables
 * @param jobs the jobs list
 * @param location the location
 */
func scheduleJobs(jobs []Job) {
	// remove all previous every day jobs
	_cronScheduler.RemoveByTag("atHour")

	for _, job := range jobs {
		switch job.ScheduledType {
		case "immediate":
			// TODO: run immediate jobs
			zlog.Info().Str("job", job.Name).Msg("Running immediate job")
		case "atHour":
			_cronScheduler.Every(1).Day().At(job.Scheduled).Tag(job.ScheduledType).Do(func(job Job) {
				zlog.Info().Str("job", job.Name).Msg("Running scheduled job")
			}, job)
		}
	}
}

func resetCronScheduler(location *time.Location) {
	if _cronScheduler != nil {
		_cronScheduler.Clear()
	}

	_cronScheduler = gocron.NewScheduler(location)
	_cronScheduler.TagsUnique()
}
