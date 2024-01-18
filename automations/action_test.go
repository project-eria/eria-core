package automations

import (
	"errors"
	"testing"

	"github.com/project-eria/go-wot/interaction"
	"github.com/project-eria/go-wot/mocks"
	"github.com/project-eria/go-wot/producer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ActionTestSuite struct {
	suite.Suite
	exposedThing  *mocks.ExposedThing
	exposedAction producer.ExposedAction
}

func Test_ActionTestSuite(t *testing.T) {
	suite.Run(t, &ActionTestSuite{})
}

func (ts *ActionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.exposedThing = &mocks.ExposedThing{}
	aAction := interaction.NewAction(
		"a",
		"No Input, No Output",
		"",
	)
	ts.exposedAction = producer.NewExposedAction(aAction)
	ts.exposedThing.On("ExposedAction", "on").Return(ts.exposedAction, nil)
	ts.exposedThing.On("ExposedAction", "off").Return(nil, errors.New("exposed action not found"))
}

func (ts *ActionTestSuite) Test_Exists() {

	got, err := getAction(ts.exposedThing, "Name", "on")
	ts.Equal(got, &Action{
		Ref:            "on",
		AutomationName: "Name",
		ExposedThing:   ts.exposedThing,
		ExposedAction:  ts.exposedAction,
		Parameters:     make(map[string]string),
	})
	ts.Nil(err)
}

func (ts *ActionTestSuite) Test_NotExists() {
	got, err := getAction(ts.exposedThing, "Name", "off")
	ts.Nil(got)
	ts.EqualError(err, "exposed action not found")
}

func (ts *ActionTestSuite) Test_Invalid() {
	got, err := getAction(ts.exposedThing, "Name", " ")
	ts.Nil(got)
	ts.EqualError(err, "missing action configuration")
}

func (ts *ActionTestSuite) Test_WithParams() {
	got, err := getAction(ts.exposedThing, "Name", "on|param1=true")
	ts.Equal(got, &Action{
		Ref:            "on",
		AutomationName: "Name",
		ExposedThing:   ts.exposedThing,
		ExposedAction:  ts.exposedAction,
		Parameters: map[string]string{
			"param1": "true",
		},
	})
	ts.Nil(err)
}
