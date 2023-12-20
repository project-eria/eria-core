package automations

import (
	"testing"
	"time"

	"github.com/project-eria/go-wot/mocks"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type GetJobTestSuite struct {
	suite.Suite
	now          time.Time
	exposedThing *mocks.ExposedThing
}

var onAction = &action{
	Ref:        "on",
	Handler:    nil,
	Parameters: make(map[string]interface{}),
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
	// 	nil,
	// 	nil,
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
	_contextsThing = &mocks.ConsumedThing{}
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
					&conditionContext{
						context: "away",
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
					&conditionContext{
						context: "away",
						invert:  true,
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
					&conditionContext{
						context: "away",
						invert:  false,
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContext{
						context: "holiday",
						invert:  true,
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
					&conditionContext{
						context: "away",
						invert:  true,
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContext{
						context: "holiday",
						invert:  true,
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
					&conditionContext{
						context: "away",
						invert:  false,
					},
				},
				schedule: &scheduleImmediate{},
			},
			{
				conditions: []Condition{
					&conditionContext{
						context: "holiday",
						invert:  false,
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
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					fixedHour: "13:00",
				},
			},
		},
	}
	j, err := automation.getJob(ts.now)
	ts.Equal(&scheduleAtHour{
		fixedHour:     "13:00",
		scheduledHour: "13:00",
	}, j)
	ts.Nil(err)
}

type ScheduleJobTestSuite struct {
	suite.Suite
	now      time.Time
	olderNow time.Time
}

func Test_ScheduleJobTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleJobTestSuite{})
}

func (ts *ScheduleJobTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	ts.olderNow = time.Date(2000, time.January, 1, 11, 0, 0, 0, time.UTC)
	initCronScheduler()
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
		action: onAction,
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
		action: onAction,
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
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					fixedHour: "13:00",
				},
			},
		},
		action:        onAction,
		lastScheduled: ts.olderNow,
		status:        "success",
		job:           &scheduleImmediate{},
	}
	ts.Equal(&scheduleImmediate{}, automation.job)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)

	automation.scheduleJob(ts.now)
	ts.NotNil(automation.job)
	ts.Equal("13:00", automation.job.(*scheduleAtHour).scheduledHour)
	ts.Equal(ts.now, automation.lastScheduled)
	ts.Equal("success", automation.status)
	// TODO does the old job been canceled?
}

// Existing Job - No Matching condition
func (ts *ScheduleJobTestSuite) Test_ExistingNoMatchingCondition() {
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						after: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					fixedHour: "13:00",
				},
			},
		},
		action:        onAction,
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
	cronJob, _ := _cronScheduler.Every(1).Day().At("12:00").Tag("atHour").Do(func() {
		zlog.Info().Msg("Running scheduled job")
	})
	automation := &Automation{
		groups: []group{
			{
				conditions: []Condition{
					&conditionTime{
						before: "13:00",
					},
				},
				schedule: &scheduleAtHour{
					fixedHour: "13:00",
				},
			},
		},
		action:        onAction,
		lastScheduled: ts.olderNow,
		status:        "success",
		job: &scheduleAtHour{
			cronJob:       cronJob,
			fixedHour:     "12:00", // Not similar to scheduled to detect if replaced
			scheduledHour: "13:00",
		},
	}
	ts.NotNil(automation.job)
	ts.Equal("12:00", automation.job.(*scheduleAtHour).fixedHour)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)

	automation.scheduleJob(ts.now)
	ts.NotNil(automation.job)
	ts.Equal("12:00", automation.job.(*scheduleAtHour).fixedHour)
	ts.Equal(ts.olderNow, automation.lastScheduled)
	ts.Equal("success", automation.status)
	// TODO the old job hasn't been canceled?
}
