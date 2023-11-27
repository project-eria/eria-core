package automations

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type JobsSchedulingTestSuite struct {
	suite.Suite
}

func Test_JobsSchedulingTestSuite(t *testing.T) {
	suite.Run(t, &JobsSchedulingTestSuite{})
}

func (ts *JobsSchedulingTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	resetCronScheduler(time.UTC)
}

func (ts *JobsSchedulingTestSuite) Test_A() {
	jobs := []Job{
		{
			Name:          "a",
			Action:        Action{},
			ScheduledType: "atHour",
			Scheduled:     "10:00:05",
		},
	}
	scheduleJobs(jobs)
	cronJobs := _cronScheduler.Jobs()
	fmt.Println(cronJobs[0].Tags())
	fmt.Println(cronJobs[0].ScheduledAtTime())
	ts.True(false)
}
