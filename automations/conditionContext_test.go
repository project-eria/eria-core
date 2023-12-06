package automations

import (
	"testing"
	"time"

	"github.com/project-eria/go-wot/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContextConditionTestSuite struct {
	suite.Suite
	consumedThingMock *mocks.ConsumedThing
	now               time.Time
}

func Test_ContextConditionTestSuite(t *testing.T) {
	suite.Run(t, &ContextConditionTestSuite{})
}

func (ts *ContextConditionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	consumedThingMock := mocks.ConsumedThing{}
	ts.consumedThingMock = &consumedThingMock
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
}

func (ts *ContextConditionTestSuite) Test_NewNoContextsThing() {
	_contextsThing = nil
	got, err := NewConditionContext([]string{"context", "test"})
	ts.Nil(got)
	ts.EqualError(err, "contexts thing not configured")
}

func (ts *ContextConditionTestSuite) Test_NewExistingContextThing() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(true, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
	}
	got, err := NewConditionContext([]string{"context", "test"})
	ts.Equal(c, got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_NewTooLong() {
	_contextsThing = ts.consumedThingMock
	got, err := NewConditionContext([]string{"context", "test", "test"})
	ts.Nil(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid condition length")
}

func (ts *ContextConditionTestSuite) Test_NewNoName() {
	_contextsThing = ts.consumedThingMock
	got, err := NewConditionContext([]string{"context"})
	ts.Nil(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid condition length")
}

func (ts *ContextConditionTestSuite) Test_NewInvalidChars() {
	_contextsThing = ts.consumedThingMock
	got, err := NewConditionContext([]string{"context", "te!st"})
	ts.Nil(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid context name")
}

func (ts *ContextConditionTestSuite) Test_NewInverted() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(false, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
		invert:  true,
	}
	got, err := NewConditionContext([]string{"context", "!test"})
	ts.Equal(c, got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_CheckContextTrue() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(true, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
		invert:  false,
	}
	got, err := c.check(ts.now)
	ts.True(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_CheckContextFalse() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(false, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
		invert:  false,
	}
	got, err := c.check(ts.now)
	ts.False(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_CheckInvertedContextTrue() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(true, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
		invert:  true,
	}
	got, err := c.check(ts.now)
	ts.False(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_CheckInvertedContextFalse() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(false, nil)
	_contextsThing = ts.consumedThingMock
	c := &conditionContext{
		context: "test",
		invert:  true,
	}
	got, err := c.check(ts.now)
	ts.True(got)
	ts.Nil(err)
}
