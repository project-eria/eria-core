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

func (ts *TimeConditionTestSuite) Test_NewInvalidLength() {
	got, err := NewConditionTime([]string{"time"})
	ts.Nil(got)
	ts.EqualError(err, "invalid condition length")
}

func (ts *TimeConditionTestSuite) Test_NewInvalidParamType() {
	got, err := NewConditionTime([]string{"time", "x=13:00"})
	ts.Nil(got)
	ts.EqualError(err, "invalid condition parameter type")
}

func (ts *TimeConditionTestSuite) Test_NewInvalidParam() {
	got, err := NewConditionTime([]string{"time", "x13:00"})
	ts.Nil(got)
	ts.EqualError(err, "invalid condition parameter")
}

func (ts *TimeConditionTestSuite) Test_NewInvalidTime() {
	got, err := NewConditionTime([]string{"time", "after=1300"})
	ts.Nil(got)
	ts.EqualError(err, "invalid condition time")
}

func (ts *TimeConditionTestSuite) Test_NewWithAfter() {
	c := &conditionTime{
		after: "11:00",
	}
	got, err := NewConditionTime([]string{"time", "after=11:00"})
	ts.Equal(c, got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_NewWithBefore() {
	c := &conditionTime{
		before: "13:00",
	}
	got, err := NewConditionTime([]string{"time", "before=13:00"})
	ts.Equal(c, got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_CheckUnexpected() {
	c := &conditionTime{}
	got, err := c.check(ts.now)
	ts.False(got)
	ts.EqualError(err, "unexpected invalid condition time")
}

func (ts *TimeConditionTestSuite) Test_CheckAfterFalse() {
	c := &conditionTime{
		after: "13:00",
	}
	got, err := c.check(ts.now)
	ts.False(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_CheckAfterTrue() {
	c := &conditionTime{
		after: "11:00",
	}
	got, err := c.check(ts.now)
	ts.True(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_CheckBeforeFalse() {
	c := &conditionTime{
		before: "11:00",
	}
	got, err := c.check(ts.now)
	ts.False(got)
	ts.Nil(err)
}

func (ts *TimeConditionTestSuite) Test_CheckBeforeTrue() {
	c := &conditionTime{
		before: "13:00",
	}
	got, err := c.check(ts.now)
	ts.True(got)
	ts.Nil(err)
}
