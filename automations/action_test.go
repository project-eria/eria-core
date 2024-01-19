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
	_exposedThings = map[string]producer.ExposedThing{
		"a": ts.exposedThing,
	}
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
	got, err := getAction([]string{"a"}, "Name", "on")
	ts.Nil(err)
	ts.Equal(&Action{
		Ref:            "on",
		AutomationName: "Name",
		ExposedThings:  map[string]producer.ExposedThing{"a": ts.exposedThing},
		Parameters:     make(map[string]string),
	}, got)
}

func (ts *ActionTestSuite) Test_OneThingNotExists() {
	got, err := getAction([]string{"a", "b"}, "Name", "on")
	ts.Nil(err)
	ts.Equal(&Action{
		Ref:            "on",
		AutomationName: "Name",
		ExposedThings:  map[string]producer.ExposedThing{"a": ts.exposedThing},
		Parameters:     make(map[string]string),
	}, got)
}

func (ts *ActionTestSuite) Test_AllThingsNotExist() {
	got, err := getAction([]string{"b"}, "Name", "on")
	ts.EqualError(err, "requested things action not found")
	ts.Nil(got)
}

func (ts *ActionTestSuite) Test_ActionNotExists() {
	got, err := getAction([]string{"a"}, "Name", "off")
	ts.EqualError(err, "requested things action not found")
	ts.Nil(got)
}
func (ts *ActionTestSuite) Test_Invalid() {
	got, err := getAction([]string{"a"}, "Name", " ")
	ts.EqualError(err, "missing action configuration")
	ts.Nil(got)
}

func (ts *ActionTestSuite) Test_WithParams() {
	got, err := getAction([]string{"a"}, "Name", "on|param1=true")
	ts.Nil(err)
	ts.Equal(&Action{
		Ref:            "on",
		AutomationName: "Name",
		ExposedThings:  map[string]producer.ExposedThing{"a": ts.exposedThing},
		Parameters: map[string]string{
			"param1": "true",
		},
	}, got)
}
