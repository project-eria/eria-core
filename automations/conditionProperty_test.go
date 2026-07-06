package automations

import (
	"errors"
	"testing"
	"time"

	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ConditionPropertyTestSuite struct {
	suite.Suite
	exposedThing    *mocks.ExposedThing
	exposedProperty *mocks.ExposedProperty
	value           interface{}
}

func Test_ConditionPropertyTestSuite(t *testing.T) {
	suite.Run(t, &ConditionPropertyTestSuite{})
}

func (ts *ConditionPropertyTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.exposedThing = &mocks.ExposedThing{}
	ts.exposedProperty = &mocks.ExposedProperty{}
	_exposedThings = map[string]producer.ExposedThing{
		"a": ts.exposedThing,
	}
	ts.exposedThing.On("ExposedProperty", "triggered").Return(ts.exposedProperty, nil)
	ts.exposedThing.On("ExposedProperty", "missing").Return(nil, errors.New("exposed property not found"))
	ts.exposedProperty.On("Read", ts.exposedThing, "triggered", map[string]string(nil)).Return(
		func(producer.ExposedThing, string, map[string]string) (interface{}, error) {
			return ts.value, nil
		},
	)
}

func (ts *ConditionPropertyTestSuite) Test_New() {
	got, err := NewConditionProperty([]string{"property", "triggered==true"}, []string{"a"})
	ts.Nil(err)
	ts.NotNil(got)
	ts.Equal("triggered", got.property)
	ts.Equal("==", got.operator)
	ts.Equal("true", got.ref)
	ts.Equal(map[string][]string{"a": {"triggered"}}, got.list())
}

func (ts *ConditionPropertyTestSuite) Test_NewInvalidSyntax() {
	got, err := NewConditionProperty([]string{"property", "triggered"}, []string{"a"})
	ts.EqualError(err, "invalid property condition syntax")
	ts.Nil(got)
}

func (ts *ConditionPropertyTestSuite) Test_NewInvalidLength() {
	got, err := NewConditionProperty([]string{"property"}, []string{"a"})
	ts.EqualError(err, "invalid condition length")
	ts.Nil(got)
}

func (ts *ConditionPropertyTestSuite) Test_NewPropertyNotFound() {
	got, err := NewConditionProperty([]string{"property", "missing==true"}, []string{"a"})
	ts.EqualError(err, "requested things property not found")
	ts.Nil(got)
}

func (ts *ConditionPropertyTestSuite) Test_NewThingNotFound() {
	got, err := NewConditionProperty([]string{"property", "triggered==true"}, []string{"b"})
	ts.EqualError(err, "requested things property not found")
	ts.Nil(got)
}

func (ts *ConditionPropertyTestSuite) Test_CheckBoolTrue() {
	condition, _ := NewConditionProperty([]string{"property", "triggered==true"}, []string{"a"})
	ts.value = true
	ok, err := condition.check(time.Now())
	ts.Nil(err)
	ts.True(ok)

	ts.value = false
	ok, err = condition.check(time.Now())
	ts.Nil(err)
	ts.False(ok)
}

func (ts *ConditionPropertyTestSuite) Test_CheckNumeric() {
	condition, _ := NewConditionProperty([]string{"property", "triggered>25"}, []string{"a"})
	ts.value = 27.5
	ok, err := condition.check(time.Now())
	ts.Nil(err)
	ts.True(ok)

	ts.value = 20
	ok, err = condition.check(time.Now())
	ts.Nil(err)
	ts.False(ok)
}

func (ts *ConditionPropertyTestSuite) Test_CheckString() {
	condition, _ := NewConditionProperty([]string{"property", "triggered!=off"}, []string{"a"})
	ts.value = "on"
	ok, err := condition.check(time.Now())
	ts.Nil(err)
	ts.True(ok)
}

func (ts *ConditionPropertyTestSuite) Test_CheckInvalidComparison() {
	condition, _ := NewConditionProperty([]string{"property", "triggered>true"}, []string{"a"})
	ts.value = true
	_, err := condition.check(time.Now())
	ts.ErrorContains(err, `unsupported operator ">" for boolean`)
}

func (ts *ConditionPropertyTestSuite) Test_GetConditionsDispatch() {
	conditions, observables, err := getConditions([]string{"property|triggered==true"}, []string{"a"})
	ts.Nil(err)
	ts.Len(conditions, 1)
	ts.IsType(&conditionProperty{}, conditions[0])
	ts.Equal(map[string][]string{"a": {"triggered"}}, observables.properties)
}
