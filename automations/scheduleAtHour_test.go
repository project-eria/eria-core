package automations

import (
	"errors"
	"testing"

	"github.com/project-eria/go-wot/consumer"
	"github.com/project-eria/go-wot/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ScheduleAtHourTestSuite struct {
	suite.Suite
	consumedThings map[string]consumer.ConsumedThing
}

func Test_ScheduleAtHourTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleAtHourTestSuite{})
}

func (ts *ScheduleAtHourTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	consumedThingMock := &mocks.ConsumedThing{}
	consumedThingMock.On("ReadProperty", "timeProperty", mock.Anything).Return("2023-11-02T14:30:10Z", nil)
	consumedThingMock.On("ReadProperty", "otherProperty", mock.Anything).Return("", errors.New("property otherProperty not found"))
	ts.consumedThings = map[string]consumer.ConsumedThing{
		"astral": consumedThingMock,
	}
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
	s := &scheduleAtHour{
		fixedHour: "12:20",
	}

	got, err := NewScheduleAtHour([]string{"at", "hour", "12:20"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHourWithSeconds() {
	s := &scheduleAtHour{
		fixedHour: "12:20:40",
	}

	got, err := NewScheduleAtHour([]string{"at", "hour", "12:20:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewFixedHourWithExtraParams() {
	s := &scheduleAtHour{
		fixedHour: "12:20",
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
	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:otherProperty"})
	ts.EqualError(err, "time thing property not available: property otherProperty not found")
	ts.Nil(got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingProperty() {
	s := &scheduleAtHour{
		timeThing:    ts.consumedThings["astral"],
		propertyHour: "timeProperty",
	}
	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMin() {
	s := &scheduleAtHour{
		timeThing:    ts.consumedThings["astral"],
		propertyHour: "timeProperty",
		min:          "14:20",
	}
	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "min=14:20"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMax() {
	s := &scheduleAtHour{
		timeThing:    ts.consumedThings["astral"],
		propertyHour: "timeProperty",
		max:          "14:40",
	}
	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "max=14:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_NewThingExistingPropertyWithMinMax() {
	s := &scheduleAtHour{
		timeThing:    ts.consumedThings["astral"],
		propertyHour: "timeProperty",
		min:          "14:20",
		max:          "14:40",
	}
	_consumedThings = ts.consumedThings
	got, err := NewScheduleAtHour([]string{"at", "hour", "astral:timeProperty", "min=14:20", "max=14:40"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleAtHourTestSuite) Test_Equal() {
	s1 := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	s2 := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_EqualExtraAttributes() {
	s1 := &scheduleAtHour{
		fixedHour:     "13:00", // This should not append, but that just for test case
		scheduledHour: "12:20",
	}
	s2 := &scheduleAtHour{
		fixedHour:     "11:00", // This should not append, but that just for test case
		scheduledHour: "12:20",
	}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualHour() {
	s1 := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	s2 := &scheduleAtHour{
		scheduledHour: "12:10",
	}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualType() {
	s1 := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	s2 := &scheduleImmediate{}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_NotEqualNil() {
	s1 := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	isEqual := s1.equals(nil)
	ts.False(isEqual)
}

func (ts *ScheduleAtHourTestSuite) Test_Job() {
	s := &scheduleAtHour{
		fixedHour: "12:20",
	}
	err := s.job()
	ts.Nil(err)
	ts.Equal(&scheduleAtHour{
		fixedHour:     "12:20",
		scheduledHour: "12:20",
	}, s)
}

func (ts *ScheduleAtHourTestSuite) Test_Start() {
	action := &MockedAction{}
	initCronScheduler()
	s := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	err := s.start(action)
	ts.Nil(err)
	ts.NotNil(s.cronJob)
	ts.Equal("12:20", s.scheduledHour)
	ts.Equal("12:20", s.cronJob.ScheduledAtTime())
	ts.Equal(_cronScheduler.Len(), 1)
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
	s := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	err := s.start(nil)
	ts.EqualError(err, "missing action")
	ts.Nil(s.cronJob)
}

func (ts *ScheduleAtHourTestSuite) Test_Cancel() {
	action := &MockedAction{}
	initCronScheduler()
	s := &scheduleAtHour{
		scheduledHour: "12:20",
	}
	err := s.start(action)
	ts.Nil(err)
	ts.NotNil(s.cronJob)
	ts.Equal("12:20", s.scheduledHour)
	ts.Equal("12:20", s.cronJob.ScheduledAtTime())
	ts.Equal(_cronScheduler.Len(), 1)
	s.cancel()
	ts.Equal(_cronScheduler.Len(), 0)
	action.AssertNotCalled(ts.T(), "run")
}
