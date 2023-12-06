package automations

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockedAction struct {
	mock.Mock
}

func (m *MockedAction) run() error {
	args := m.Called()
	return args.Error(0)
}

type ScheduleImmediateTestSuite struct {
	suite.Suite
}

func Test_ScheduleImmediateTestSuite(t *testing.T) {
	suite.Run(t, &ScheduleImmediateTestSuite{})
}

func (ts *ScheduleImmediateTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *ScheduleImmediateTestSuite) Test_NewNoParams() {
	s := &scheduleImmediate{}

	got, err := NewScheduleImmediate([]string{"immediate"})
	ts.Nil(err)
	ts.Equal(s, got)
}

func (ts *ScheduleImmediateTestSuite) Test_NewTooManyParams() {
	got, err := NewScheduleImmediate([]string{"immediate", "x"})
	ts.EqualError(err, "invalid immediate schedule length")
	ts.Nil(got)
}

func (ts *ScheduleImmediateTestSuite) Test_Start() {
	action := &MockedAction{}
	s := &scheduleImmediate{}

	// set up expectations
	action.On("run", mock.Anything).Return(nil)
	err := s.start(action)

	ts.Nil(err)
	action.AssertCalled(ts.T(), "run")
}

func (ts *ScheduleImmediateTestSuite) Test_StartMissingAction() {
	s := &scheduleImmediate{}

	err := s.start(nil)
	ts.EqualError(err, "missing action")
}

func (ts *ScheduleImmediateTestSuite) Test_Equal() {
	s1 := &scheduleImmediate{}
	s2 := &scheduleImmediate{}
	isEqual := s1.equals(s2)
	ts.True(isEqual)
}

func (ts *ScheduleImmediateTestSuite) Test_NotEqualType() {
	s1 := &scheduleImmediate{}
	s2 := &scheduleAtHour{}
	isEqual := s1.equals(s2)
	ts.False(isEqual)
}

func (ts *ScheduleImmediateTestSuite) Test_NotEqualNil() {
	s1 := &scheduleImmediate{}
	isEqual := s1.equals(nil)
	ts.False(isEqual)
}

// TODO Test schedule/cancel combinations
