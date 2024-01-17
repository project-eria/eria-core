package automations

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ScheduleNoneTestSuite struct {
	suite.Suite
}

func Test_ScheduleNoneTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleNoneTestSuite{})
}

func (ts *ScheduleNoneTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *ScheduleNoneTestSuite) Test_NewNoParams() {
	s := &scheduleNone{}

	got, err := NewScheduleNone([]string{"none"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleNoneTestSuite) Test_NewTooManyParams() {
	got, err := NewScheduleNone([]string{"none", "x"})
	ts.EqualError(err, "invalid none schedule length")
	ts.Nil(got)
}

func (ts *ScheduleNoneTestSuite) Test_Start() {
	action := &MockedAction{}
	s := &scheduleNone{}

	// set up expectations
	action.On("run", mock.Anything).Return(nil)
	err := s.start(action)

	ts.Nil(err)
	action.AssertNotCalled(ts.T(), "run")
}

func (ts *ScheduleNoneTestSuite) Test_StartMissingAction() {
	s := &scheduleNone{}

	err := s.start(nil)
	ts.EqualError(err, "missing action")
}

func (ts *ScheduleNoneTestSuite) Test_Equal() {
	s1 := &scheduleNone{}
	s2 := &scheduleNone{}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleNoneTestSuite) Test_NotEqualType() {
	s1 := &scheduleNone{}
	s2 := &scheduleAtHour{}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleNoneTestSuite) Test_NotEqualNil() {
	s1 := &scheduleNone{}
	isEqual := s1.equals(nil)
	ts.False(isEqual)
}
