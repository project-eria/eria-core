package automations

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type TimeConditionTestSuite struct {
	suite.Suite
	now time.Time
}

func Test_TimeConditionTestSuite(t *testing.T) {
	suite.Run(t, &TimeConditionTestSuite{})
}

func (ts *TimeConditionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.now = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
}

func (ts *TimeConditionTestSuite) Test_timeConditionInvalidLength() {
	got, err := timeCondition([]string{"time"}, ts.now)
	ts.False(got)
	ts.EqualError(err, "invalid condition length")
}

func (ts *TimeConditionTestSuite) Test_timeConditionInvalidParamType() {
	got, err := timeCondition([]string{"time", "x=13:00"}, ts.now)
	ts.False(got)
	ts.EqualError(err, "invalid condition parameter type")
}

func (ts *TimeConditionTestSuite) Test_timeConditionInvalidParam() {
	got, err := timeCondition([]string{"time", "x13:00"}, ts.now)
	ts.False(got)
	ts.EqualError(err, "invalid condition parameter")
}

func (ts *TimeConditionTestSuite) Test_timeConditionInvalidTime() {
	got, err := timeCondition([]string{"time", "after=1300"}, ts.now)
	ts.False(got)
	ts.EqualError(err, "invalid condition time")
}

func (ts *TimeConditionTestSuite) Test_timeConditionAfterTrue() {
	got, err := timeCondition([]string{"time", "after=11:00"}, ts.now)
	ts.True(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_timeConditionAfterFalse() {
	got, err := timeCondition([]string{"time", "after=13:00"}, ts.now)
	ts.False(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_timeConditionBeforeTrue() {
	got, err := timeCondition([]string{"time", "before=13:00"}, ts.now)
	ts.True(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_timeConditionBeforeFalse() {
	got, err := timeCondition([]string{"time", "before=11:00"}, ts.now)
	ts.False(got)
	ts.Nil(err)
}
