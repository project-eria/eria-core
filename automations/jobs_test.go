package automations

import (
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type GetJobTestSuite struct {
	suite.Suite
	now time.Time
}

func Test_GetJobTestSuite(t *testing.T) {
	suite.Run(t, &GetJobTestSuite{})
}

func (ts *GetJobTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)

	// // exposedThing
	// ts.exposedThing = &mocks.ExposedThing{}
	// aAction := interaction.NewAction(
	// 	"a",
	// 	"No Input, No Output",
	// 	"",
	// )
	// exposedAction := producer.NewExposedAction(aAction)
	// ts.exposedThing.On("ExposedAction", "on").Return(exposedAction, nil)
	// ts.exposedThing.On("ExposedAction", "off").Return(nil, errors.New("exposed action not found"))

	// // consumedThings
	// astralThingMock := &mocks.ConsumedThing{}
	// astralThingMock.On("ReadProperty", "timeProperty").Return("2023-11-02T14:30:10Z", nil)
	// astralThingMock.On("ReadProperty", "otherProperty").Return("", errors.New("property otherProperty not found"))
	// _consumedThings = map[string]consumer.ConsumedThing{
	// 	"astral": astralThingMock,
	// }

	// contextsThing
	//	_contextsThing = &mocks.ConsumedThing{}
	_activeContexts = []string{"holiday"}
	// _location = time.UTC
}

// No groups
func (ts *GetJobTestSuite) Test_NoGroups() {
	automation := &Automation{
		groups: []group{},
	}
	j, err := automation.getJob(ts.now)
	ts.Nil(j)
	ts.EqualError(err, "missing conditions")
}

// No matching time condition
func (ts *GetJobTestSuite) Test_NoMatchingTimeCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						after: "13:00",
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Nil(j)
	ts.EqualError(err, "no matching conditions")
}

// Matching time condition
func (ts *GetJobTestSuite) Test_MatchingTimeCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleImmediate{}, j)
	ts.Nil(err)
}

// Group with empty conditions
func (ts *GetJobTestSuite) Test_GroupWithEmptyCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{},
				schedule:   &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleImmediate{}, j)
	ts.Nil(err)
}

// No Matching multiple conditions
func (ts *GetJobTestSuite) Test_NoMatchingMultipleConditions() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
					&conditionContexts{
						{
							context: "away",
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Nil(j)
	ts.EqualError(err, "no matching conditions")
}

// Matching multiple conditions
func (ts *GetJobTestSuite) Test_MatchingMultipleConditions() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
					&conditionContexts{
						{
							context: "away",
							invert:  true,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleImmediate{}, j)

	ts.Nil(err)
}

// No Matching multiple groups",
func (ts *GetJobTestSuite) Test_NoMatchingMultipleGroups() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "away",
							invert:  false,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "holiday",
							invert:  true,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}

	j, err := automation.getJob(ts.now)
	ts.Nil(j)
	ts.EqualError(err, "no matching conditions")
}

// Matching multiple groups #1
func (ts *GetJobTestSuite) Test_MatchingMultipleGroups1() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "away",
							invert:  true,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "holiday",
							invert:  true,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleImmediate{}, j)
	ts.Nil(err)
}

// Matching multiple groups #2
func (ts *GetJobTestSuite) Test_MatchingMultipleGroups2() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "away",
							invert:  false,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContexts{
						{
							context: "holiday",
							invert:  false,
						},
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleImmediate{}, j)
	ts.Nil(err)
}

// Matching condition with at hour create
func (ts *GetJobTestSuite) Test_MatchingConditionWithAtHour() {
	t, _ := time.Parse("15:04", "13:00")
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					scheduledTime: &t,
				},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleAtHour{
		scheduledTime: &t,
	}, j)
	ts.Nil(err)
}

type ScheduleJobTestSuite struct {
	suite.Suite
	now      time.Time
	olderNow time.Time
	//	exposedAction *mocks.ExposedAction
	onAction *Action
}

func Test_ScheduleJobTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleJobTestSuite{})
}

func (ts *ScheduleJobTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	ts.olderNow = time.Date(2000, time.January, 1, 11, 0, 0, 0, time.UTC)
	_cronScheduler, _ = gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
	)
	exposedThing := &mocks.ExposedThing{}
	exposedAction := &mocks.ExposedAction{}
	exposedAction.On("Run", exposedThing, "on", nil, map[string]string{}).Return(nil, nil)
	exposedThing.On("ExposedAction", "on").Return(exposedAction, nil)

	ts.onAction = &Action{
		AutomationName: "",
		Ref:            "on",
		ExposedThings:  map[string]producer.ExposedThing{"": exposedThing},
		Parameters:     make(map[string]string),
	}

}

// Initial Job - Matching condition
func (ts *ScheduleJobTestSuite) Test_InitialMatchingCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
		action: ts.onAction,
	}
	ts.Nil(automation.job)
	ts.Equal(time.Time{}, automation.lastScheduled)
	ts.Empty(automation.status)

	automation.scheduleJob(ts.now)

	ts.Equal(&scheduleImmediate{}, automation.job)
	ts.Equal(ts.now, automation.lastScheduled)
	ts.Equal("success", automation.status)
}

// Initial Job - Not Matching condition
func (ts *ScheduleJobTestSuite) Test_InitialNoMatchingCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						after: "13:00",
					},
				},
				schedule: &scheduleImmediate{},
			},
		},
		action: ts.onAction,
	}
	ts.Nil(automation.job)
	ts.Equal(time.Time{}, automation.lastScheduled)
	ts.Empty(automation.status)

	automation.scheduleJob(ts.now)

	ts.Nil(automation.job)
	ts.Equal(time.Time{}, automation.lastScheduled)
	ts.Equal("no matching conditions", automation.status)
}

// Existing Job - Matchingcondition, replacing job
func (ts *ScheduleJobTestSuite) Test_ExistingMatchingConditionReplacingJob() {
	t, _ := time.Parse("15:04", "13:00")
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					scheduledTime: &t,
				},
			},
		},
		action:        ts.onAction,
		lastScheduled: ts.olderNow,
		status:        "success",
		job:           &scheduleImmediate{},
	}
	ts.Equal(&scheduleImmediate{}, automation.job)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)

	automation.scheduleJob(ts.now)
	ts.NotNil(automation.job)
	ts.Equal("13:00", automation.job.(*scheduleAtHour).scheduledTime.Format("15:04"))
	ts.Equal(ts.now, automation.lastScheduled)
	ts.Equal("success", automation.status)
	// TODO does the old job been canceled?
}

// Existing Job - No Matching condition
func (ts *ScheduleJobTestSuite) Test_ExistingNoMatchingCondition() {
	t, _ := time.Parse("15:04", "13:00")
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						after: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					scheduledTime: &t,
				},
			},
		},
		action:        ts.onAction,
		lastScheduled: ts.olderNow,
		status:        "success",
		job:           &scheduleImmediate{},
	}
	ts.Equal(&scheduleImmediate{}, automation.job)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)

	automation.scheduleJob(ts.now)
	ts.Nil(automation.job)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("no matching conditions", automation.status)
	// TODO does the job been canceled?
}

// Existing Job - Matching condition, indentical job
func (ts *ScheduleJobTestSuite) Test_ExistingMatchingConditionIdenticalJob() {
	cronJob, _ := _cronScheduler.NewJob(
		gocron.DailyJob(1,
			gocron.NewAtTimes(
				gocron.NewAtTime(12, 0, 0),
			),
		),
		gocron.NewTask(
			func() {
				zlog.Info().Msg("Running scheduled job")
			},
		),
		gocron.WithTags("core", "automation", "atHour"),
	)
	_cronScheduler.Start()
	t, _ := time.Parse("15:04", "13:00")
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					scheduledTime: &t,
				},
			},
		},
		action:        ts.onAction,
		lastScheduled: ts.olderNow,
		status:        "success",
		job: &scheduleAtHour{
			cronJob:       cronJob,
			scheduledTime: &t,
		},
	}
	ts.NotNil(automation.job)
	j, _ := cronJob.NextRun()
	ts.Equal("12:00", j.Format("15:04"))

	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)

	automation.scheduleJob(ts.now)
	ts.NotNil(automation.job)
	j, _ = cronJob.NextRun()
	ts.Equal("12:00", j.Format("15:04"))

	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)
	// TODO the old job hasn't been canceled?
}
