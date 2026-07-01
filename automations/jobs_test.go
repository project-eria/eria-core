package automations

import (
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// fakeAction is a minimal ActionRunner for tests.
type fakeAction struct{ ran bool }

func (f *fakeAction) run() error { f.ran = true; return nil }

// ---- plan() ----

type PlanTestSuite struct {
	suite.Suite
	now time.Time
}

func Test_PlanTestSuite(t *testing.T) {
	suite.Run(t, &PlanTestSuite{})
}

func (ts *PlanTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	_activeContexts = []string{"holiday"}
}

// No groups -> nothing planned.
func (ts *PlanTestSuite) Test_NoGroups() {
	a := &Automation{}
	ts.Empty(a.plan(ts.now))
}

// Single branch, condition not matching -> nothing planned.
func (ts *PlanTestSuite) Test_NoMatch() {
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionTime{after: "13:00"}}, schedule: &scheduleImmediate{}},
		},
	}}}
	ts.Empty(a.plan(ts.now))
}

// Single branch, condition matching -> one planned action.
func (ts *PlanTestSuite) Test_Match() {
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionTime{before: "13:00"}}, schedule: &scheduleImmediate{}},
		},
	}}}
	p := a.plan(ts.now)
	ts.Len(p, 1)
	ts.Equal(&scheduleImmediate{}, p[0].schedule)
}

// A branch with no conditions is the default and matches when no earlier branch did.
func (ts *PlanTestSuite) Test_DefaultBranch() {
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionContexts{{context: "away"}}}, schedule: &scheduleImmediate{}}, // no match
			{schedule: &scheduleImmediate{}}, // default
		},
	}}}
	ts.Len(a.plan(ts.now), 1)
}

// The first matching branch wins; later branches are ignored.
func (ts *PlanTestSuite) Test_FirstMatchWins() {
	t8, _ := time.Parse("15:04", "08:00")
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionContexts{{context: "holiday"}}}, schedule: &scheduleAtHour{scheduledTime: &t8}}, // matches
			{schedule: &scheduleImmediate{}}, // default, must NOT be used
		},
	}}}
	p := a.plan(ts.now)
	ts.Len(p, 1)
	_, ok := p[0].schedule.(*scheduleAtHour)
	ts.True(ok)
}

// A matching branch without a schedule disables the group (nothing planned).
func (ts *PlanTestSuite) Test_Disabled() {
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionContexts{{context: "holiday"}}}}, // matches, no schedule -> disabled
			{schedule: &scheduleImmediate{}},                                     // not reached
		},
	}}}
	ts.Empty(a.plan(ts.now))
}

// Each group plans its own action independently.
func (ts *PlanTestSuite) Test_MultipleGroups() {
	a := &Automation{groups: []actionGroup{
		{action: &fakeAction{}, entries: []scheduleEntry{{schedule: &scheduleImmediate{}}}}, // e.g. open
		{action: &fakeAction{}, entries: []scheduleEntry{{schedule: &scheduleImmediate{}}}}, // e.g. close
	}}
	ts.Len(a.plan(ts.now), 2)
}

// ---- scheduleJob() ----

type ScheduleJobTestSuite struct {
	suite.Suite
	now time.Time
}

func Test_ScheduleJobTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleJobTestSuite{})
}

func (ts *ScheduleJobTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	_activeContexts = []string{}
	_cronScheduler, _ = gocron.NewScheduler(gocron.WithLocation(time.UTC))
}

// Matching condition -> job started, action run (immediate), success.
func (ts *ScheduleJobTestSuite) Test_InitialMatching() {
	act := &fakeAction{}
	a := &Automation{groups: []actionGroup{{
		action: act,
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionTime{before: "13:00"}}, schedule: &scheduleImmediate{}},
		},
	}}}
	a.scheduleJob(ts.now)
	ts.Len(a.jobs, 1)
	ts.Equal("success", a.status)
	ts.Equal(ts.now, a.lastScheduled)
	ts.True(act.ran)
}

// No matching condition -> no jobs.
func (ts *ScheduleJobTestSuite) Test_InitialNoMatch() {
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{conditions: []Condition{&conditionTime{after: "13:00"}}, schedule: &scheduleImmediate{}},
		},
	}}}
	a.scheduleJob(ts.now)
	ts.Empty(a.jobs)
	ts.Equal("success", a.status)
}

// Re-planning cancels the previous jobs and re-creates the planned ones.
func (ts *ScheduleJobTestSuite) Test_ReplanReplaces() {
	t13, _ := time.Parse("15:04", "13:00")
	a := &Automation{groups: []actionGroup{{
		action: &fakeAction{},
		entries: []scheduleEntry{
			{schedule: &scheduleAtHour{scheduledTime: &t13}},
		},
	}}}
	a.scheduleJob(ts.now)
	ts.Len(a.jobs, 1)

	a.scheduleJob(ts.now) // re-plan: still one job after cancel + restart
	ts.Len(a.jobs, 1)
	ts.Equal("success", a.status)
}
