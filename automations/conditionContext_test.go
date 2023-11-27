package automations

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/project-eria/go-wot/mocks"
	"github.com/rs/zerolog"
)

type ContextConditionTestSuite struct {
	suite.Suite
	consumedThingMock *mocks.ConsumedThing
}

func Test_ContextConditionTestSuite(t *testing.T) {
	suite.Run(t, &ContextConditionTestSuite{})
}

func (ts *ContextConditionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	consumedThingMock := mocks.ConsumedThing{}
	ts.consumedThingMock = &consumedThingMock
}

func (ts *ContextConditionTestSuite) Test_contextConditionActive() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(true, nil)
	got, err := contextCondition([]string{"context", "test"}, ts.consumedThingMock)
	ts.True(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_contextConditionNotActive() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(false, nil)
	got, err := contextCondition([]string{"context", "test"}, ts.consumedThingMock)
	ts.False(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_contextConditionTooLong() {
	got, err := contextCondition([]string{"context", "test", "test"}, ts.consumedThingMock)
	ts.False(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid condition length")
}

func (ts *ContextConditionTestSuite) Test_contextConditionNoName() {
	got, err := contextCondition([]string{"context"}, ts.consumedThingMock)
	ts.False(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid condition length")
}

func (ts *ContextConditionTestSuite) Test_contextConditionInvalidChars() {
	got, err := contextCondition([]string{"context", "te!st"}, ts.consumedThingMock)
	ts.False(got)
	ts.consumedThingMock.AssertNotCalled(ts.T(), "ReadProperty", mock.AnythingOfType("string"))
	ts.EqualError(err, "invalid context name")
}

func (ts *ContextConditionTestSuite) Test_contextConditionInverted1() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(false, nil)
	got, err := contextCondition([]string{"context", "!test"}, ts.consumedThingMock)
	ts.True(got)
	ts.Nil(err)
}

func (ts *ContextConditionTestSuite) Test_contextConditionInverted2() {
	ts.consumedThingMock.On("ReadProperty", mock.AnythingOfType("string")).Return(true, nil)
	got, err := contextCondition([]string{"context", "!test"}, ts.consumedThingMock)
	ts.False(got)
	ts.Nil(err)
}
