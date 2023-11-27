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

func (ts *ActionTestSuite) Test_actionExists() {
	got, err := getAction(ts.exposedThing, []string{"on"})
	ts.Equal(got, Action{
		Ref:        "on",
		Handler:    nil,
		Parameters: make(map[string]interface{}),
	})
	ts.Nil(err)
}

func (ts *ActionTestSuite) Test_actionNotExists() {
	got, err := getAction(ts.exposedThing, []string{"off"})
	ts.Equal(got, Action{})
	ts.EqualError(err, "exposed action not found")
}

func (ts *ActionTestSuite) Test_actionInvalid() {
	got, err := getAction(ts.exposedThing, []string{})
	ts.Equal(got, Action{})
	ts.EqualError(err, "invalid action length")
}

func (ts *ActionTestSuite) Test_actionWithParams() {
	got, err := getAction(ts.exposedThing, []string{"on", "param1=true"})
	ts.Equal(got, Action{
		Ref:     "on",
		Handler: nil,
		Parameters: map[string]interface{}{
			"param1": "true",
		},
	})
	ts.Nil(err)
}
