package automations

import (
	"errors"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/project-eria/eria-core/consumer/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ScheduleAtHourTestSuite struct {
	suite.Suite
	astralThing  *mocks.Thing
	timeProperty *mocks.Property
	observer     eriaconsumer.PropertyObserver
}

func Test_ScheduleAtHourTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleAtHourTestSuite{})
}

func (ts *ScheduleAtHourTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	_consumer = &mocks.Consumer{}
	ts.astralThing = &mocks.Thing{}
	ts.timeProperty = &mocks.Property{}
	ts.observer = nil
	ts.timeProperty.On("Value").Return("2023-11-02T14:30:10Z", nil)
	ts.timeProperty.On("Observe", mock.AnythingOfTypeArgument("eriaconsumer.PropertyObserver")).Return(func(s eriaconsumer.PropertyObserver) uint16 {
		ts.observer = s
		return uint16(0)
	}, nil)
	ts.timeProperty.On("UnObserve", mock.Anything).Return()
	ts.astralThing.On("Property", "timeProperty").Return(ts.timeProperty)
	otherProperty := &mocks.Property{}
	otherProperty.On("Value").Return("", errors.New("property otherProperty not found"))
	ts.astralThing.On("Property", "otherProperty").Return(otherProperty)
	_consumer.(*mocks.Consumer).On("ThingFromTag", "astral").Return(ts.astralThing)
	_consumer.(*mocks.Consumer).On("ThingFromTag", mock.AnythingOfType("string")).Return(nil)
}

func (ts *ScheduleAtHourTestSuite) Test_NewNoParams() {
	got, err := NewScheduleAtHour([]string{"at", "hour"})
	ts.EqualError(err, "invalid at scheduling length")
	ts.Nil(got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHourIncorrect() {
	got, err := NewScheduleAtHour([]string{"at", "hour", "1200"})
	ts.EqualError(err, "invalid time")
	ts.Nil(got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHour() {
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}

	got, err := NewScheduleAtHour([]string{"at", "hour", "12:20"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHourWithSeconds() {
	t, _ := time.Parse("15:04:05", "12:20:40")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}

	got, err := NewScheduleAtHour([]string{"at", "hour", "12:20:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHourWithExtraParams() {
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}

	got, err := NewScheduleAtHour([]string{"at", "hour", "12:20", "x"})
	ts.Nil(err)
	ts.Equal(s, got)
}
func (ts *ScheduleAtHourTestSuite) Test_NewThingNotExisting() {
	got, err := NewScheduleAtHour([]string{"at", "hour", "mything:myproperty"})
	ts.EqualError(err, "time thing doesn't exist")
	ts.Nil(got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingNotExistingProperty() {
	//	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:otherProperty"})
	ts.EqualError(err, "time thing property not available: property otherProperty not found")
	ts.Nil(got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingProperty() {
	s := &scheduleAtHour{
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
	}
	//	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMin() {
	min, _ := time.Parse("15:04", "14:20")
	s := &scheduleAtHour{
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
		min:               &min,
	}
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "min=14:20"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMax() {
	max, _ := time.Parse("15:04", "14:40")
	s := &scheduleAtHour{
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
		max:               &max,
	}
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "max=14:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMinMax() {
	min, _ := time.Parse("15:04", "14:20")
	max, _ := time.Parse("15:04", "14:40")
	s := &scheduleAtHour{
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
		min:               &min,
		max:               &max,
	}
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "min=14:20", "max=14:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_Equal() {
	t1, _ := time.Parse("15:04", "12:20")
	t2, _ := time.Parse("15:04", "12:20")
	s1 := &scheduleAtHour{
		scheduledTime: &t1,
	}
	s2 := &scheduleAtHour{
		scheduledTime: &t2,
	}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_EqualExtraAttributes() {
	t1, _ := time.Parse("15:04", "12:20")
	t2, _ := time.Parse("15:04", "12:20")
	s1 := &scheduleAtHour{
		scheduledTime: &t1,
	}
	s2 := &scheduleAtHour{
		scheduledTime: &t2,
	}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualHour() {
	t1, _ := time.Parse("15:04", "12:20")
	t2, _ := time.Parse("15:04", "12:10")
	s1 := &scheduleAtHour{
		scheduledTime: &t1,
	}
	s2 := &scheduleAtHour{
		scheduledTime: &t2,
	}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualType() {
	t1, _ := time.Parse("15:04", "12:20")
	s1 := &scheduleAtHour{
		scheduledTime: &t1,
	}
	s2 := &scheduleImmediate{}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualNil() {
	t1, _ := time.Parse("15:04", "12:20")
	s1 := &scheduleAtHour{
		scheduledTime: &t1,
	}
	isEqual := s1.equals(nil)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_Job() {
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}
	err := s.job()
	ts.Nil(err)
	ts.Equal(&scheduleAtHour{
		scheduledTime: &t,
	}, s)
}

func (ts *ScheduleAtHourTestSuite) Test_Start() {
	action := &MockedAction{}
	_cronScheduler, _ = gocron.NewScheduler()
	_cronScheduler.Start()
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}
	err := s.start(action)
	ts.Nil(err)
	ts.NotNil(s.cronJob)
	ts.Equal("12:20", s.scheduledTime.Format("15:04"))
	j, _ := s.cronJob.NextRun()
	ts.Equal("12:20", j.Format("15:04"))
	ts.Equal(len(_cronScheduler.Jobs()), 1)
	action.AssertNotCalled(ts.T(), "run")
}

func (ts *ScheduleAtHourTestSuite) Test_StartMissingScheduled() {
	action := &MockedAction{}
	s := &scheduleAtHour{}
	err := s.start(action)
	ts.EqualError(err, "missing scheduled hour")
	ts.Nil(s.cronJob)
}

func (ts *ScheduleAtHourTestSuite) Test_StartMissingAction() {
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}
	err := s.start(nil)
	ts.EqualError(err, "missing action")
	ts.Nil(s.cronJob)
}

func (ts *ScheduleAtHourTestSuite) Test_Cancel() {
	action := &MockedAction{}
	_cronScheduler, _ = gocron.NewScheduler()
	_cronScheduler.Start()
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime: &t,
	}
	err := s.start(action)
	ts.Nil(err)
	ts.NotNil(s.cronJob)
	ts.Equal("12:20", s.scheduledTime.Format("15:04"))
	j, _ := s.cronJob.NextRun()
	ts.Equal("12:20", j.Format("15:04"))
	ts.Equal(len(_cronScheduler.Jobs()), 1)
	s.cancel()
	ts.Equal(len(_cronScheduler.Jobs()), 0)
	action.AssertNotCalled(ts.T(), "run")
}

func (ts *ScheduleAtHourTestSuite) Test_StartWithPropertyHour() {
	action := &MockedAction{}
	_cronScheduler, _ = gocron.NewScheduler()
	_cronScheduler.Start()
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime:     &t,
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
	}
	err := s.start(action)
	ts.Nil(err)
	ts.NotNil(s.cronJob)
	ts.Equal("12:20", s.scheduledTime.Format("15:04"))
	j, _ := s.cronJob.NextRun()
	ts.Equal("12:20", j.Format("15:04"))
	ts.Equal(len(_cronScheduler.Jobs()), 1)
	ts.timeProperty.AssertCalled(ts.T(), "Observe", mock.Anything)
	ts.NotNil(ts.observer)
	action.AssertNotCalled(ts.T(), "run")
}

func (ts *ScheduleAtHourTestSuite) Test_StartWithPropertyHourChange() {
	action := &MockedAction{}
	_cronScheduler, _ = gocron.NewScheduler()
	t, _ := time.Parse("15:04", "12:20")
	s := &scheduleAtHour{
		scheduledTime:     &t,
		timeThing:         ts.astralThing,
		timeThingProperty: "timeProperty",
	}
	_cronScheduler.Start()
	err := s.start(action)
	ts.Nil(err)
	var old = s.cronJob
	ts.Equal("12:20", s.scheduledTime.Format("15:04"))
	j, _ := s.cronJob.NextRun()
	ts.Equal("12:20", j.Format("15:04"))
	ts.observer("2000-01-01T15:00:00+01:00", nil)
	ts.Equal("15:00", s.scheduledTime.Format("15:04"))
	j, _ = s.cronJob.NextRun()
	ts.Equal("15:00", j.Format("15:04"))
	ts.NotEqual(old, s.cronJob)
	ts.Equal(len(_cronScheduler.Jobs()), 1)
	ts.timeProperty.AssertCalled(ts.T(), "Observe", mock.Anything)
	ts.NotNil(ts.observer)
	action.AssertNotCalled(ts.T(), "run")
}
