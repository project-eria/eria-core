package automations

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type GetConditionsTestSuite struct {
	suite.Suite
}

func Test_GetConditionsTestSuite(t *testing.T) {
	suite.Run(t, &GetConditionsTestSuite{})
}

func (ts *GetConditionsTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	newContextCondition = func(_ []string) (*conditionContexts, error) {
		return &conditionContexts{
			{
				context: "test",
			},
		}, nil
	}
	newTimeCondition = func(_ []string) (*conditionTime, error) {
		return &conditionTime{}, nil
	}
}

func (ts *GetConditionsTestSuite) TearDownTest() {
	newContextCondition = NewConditionContext
	newTimeCondition = NewConditionTime
}

func (ts *GetConditionsTestSuite) Test_InvalidCondition() {
	got, obs, err := getConditions([]string{"unknown"})
	ts.Nil(got)
	ts.Nil(obs)
	ts.EqualError(err, "invalid condition type")
}

func (ts *GetConditionsTestSuite) Test_SimpleContextCondition() {
	got, obs, err := getConditions([]string{"context|test"})
	ts.Nil(err)
	ts.Equal([]string{"test"}, obs.contexts)
	ts.Len(got, 1)
	ts.Equal(&conditionContexts{
		{
			context: "test",
		},
	}, got[0])
}

func (ts *GetConditionsTestSuite) Test_ContextConditionWithoutContextThing() {
	newContextCondition = func(_ []string) (*conditionContexts, error) {
		return nil, errors.New("contexts thing not configured")
	}
	got, obs, err := getConditions([]string{"context|test"})
	ts.EqualError(err, "contexts thing not configured")
	ts.Nil(got)
	ts.Nil(obs)
}

func (ts *GetConditionsTestSuite) Test_SimpleTimeCondition() {
	got, obs, err := getConditions([]string{"time|test"})
	ts.Nil(err)
	ts.Equal([]string{}, obs.contexts)
	ts.Len(got, 1)
	ts.Equal(&conditionTime{}, got[0])
}

func (ts *GetConditionsTestSuite) Test_DualCondition() {
	got, obs, err := getConditions([]string{"context|test", "time|test"})
	ts.Nil(err)
	ts.Equal([]string{"test"}, obs.contexts)
	ts.Len(got, 2)
	ts.Equal(&conditionContexts{
		{
			context: "test",
			invert:  false,
		},
	}, got[0])
	ts.Equal(&conditionTime{}, got[1])
}
