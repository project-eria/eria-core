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

func handler(interface{}, map[string]string) (interface{}, error) {
	return nil, nil
}

type ActionTestSuite struct {
	suite.Suite
	exposedThing *mocks.ExposedThing
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
		nil,
		nil,
	)
	exposedAction := producer.NewExposedAction(aAction)
	ts.exposedThing.On("ExposedAction", "on").Return(exposedAction, nil)
	ts.exposedThing.On("ExposedAction", "off").Return(nil, errors.New("exposed action not found"))
}

func (ts *ActionTestSuite) Test_Exists() {
	got, err := getAction(ts.exposedThing, "on")
	ts.Equal(got, &action{
		Ref:        "on",
		Handler:    nil,
		Parameters: make(map[string]interface{}),
	})
	ts.Nil(err)
}

func (ts *ActionTestSuite) Test_NotExists() {
	got, err := getAction(ts.exposedThing, "off")
	ts.Nil(got)
	ts.EqualError(err, "exposed action not found")
}

func (ts *ActionTestSuite) Test_Invalid() {
	got, err := getAction(ts.exposedThing, " ")
	ts.Nil(got)
	ts.EqualError(err, "missing action configuration")
}

func (ts *ActionTestSuite) Test_WithParams() {
	got, err := getAction(ts.exposedThing, "on|param1=true")
	ts.Equal(got, &action{
		Ref:     "on",
		Handler: nil,
		Parameters: map[string]interface{}{
			"param1": "true",
		},
	})
	ts.Nil(err)
}
