package automations

import (
	"errors"
	"time"

	"testing"

	eriaconsumer "github.com/project-eria/eria-core/consumer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

// invokedAction records one InvokeAction call on fakeNoticesThing.
type invokedAction struct {
	name  string
	input interface{}
}

// fakeNoticesThing is a consumed Thing recording the invoked actions.
type fakeNoticesThing struct {
	calls []invokedAction
	err   error
}

func (f *fakeNoticesThing) Property(string, ...eriaconsumer.PropertyOption) eriaconsumer.Property {
	return nil
}

func (f *fakeNoticesThing) InvokeAction(name string, _ map[string]interface{}, value interface{}) (interface{}, error) {
	f.calls = append(f.calls, invokedAction{name: name, input: value})
	return nil, f.err
}

// fakeNoticesConsumer returns thing for NoticesThingTag, nil otherwise.
type fakeNoticesConsumer struct{ thing eriaconsumer.Thing }

func (f *fakeNoticesConsumer) Thing(string) eriaconsumer.Thing { return nil }
func (f *fakeNoticesConsumer) ThingFromTag(tag string) eriaconsumer.Thing {
	if tag == NoticesThingTag {
		return f.thing
	}
	return nil
}

// switchCondition is a condition whose result is set by the test.
type switchCondition struct{ on bool }

func (c *switchCondition) check(time.Time) (bool, error) { return c.on, nil }

type NotifyActionTestSuite struct {
	suite.Suite
	thing *fakeNoticesThing
}

func Test_NotifyActionTestSuite(t *testing.T) {
	suite.Run(t, &NotifyActionTestSuite{})
}

func (ts *NotifyActionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	ts.thing = &fakeNoticesThing{}
	_consumer = &fakeNoticesConsumer{thing: ts.thing}
}

func (ts *NotifyActionTestSuite) TearDownTest() {
	_consumer = nil
}

func (ts *NotifyActionTestSuite) Test_ParseMinimal() {
	got, err := getNotifyAction(nil, "Name", "notify|key=door|title=Porte ouverte")
	ts.Nil(err)
	ts.Equal(&NotifyAction{
		AutomationName: "Name",
		Key:            "door",
		Severity:       "info",
		Title:          "Porte ouverte",
		Audience:       "all",
	}, got)
}

func (ts *NotifyActionTestSuite) Test_ParseFull() {
	got, err := getNotifyAction([]string{"garage", "other"}, "Name", "notify|key=door|severity=alert|title=Porte|body=a=b|audience=profile")
	ts.Nil(err)
	ts.Equal(&NotifyAction{
		AutomationName: "Name",
		Key:            "door",
		Severity:       "alert",
		Title:          "Porte",
		Body:           "a=b",
		Audience:       "profile",
		Thing:          "garage",
	}, got)
}

func (ts *NotifyActionTestSuite) Test_ParseErrors() {
	cases := map[string]string{
		"notify|title=T":                       "missing notify key",
		"notify|key=k":                         "missing notify title",
		"notify|key=k|title=T|severity=fatal":  "unsupported notify severity",
		"notify|key=k|title=T|audience=nobody": "unsupported notify audience",
		"notify|key=k|title=T|foo=bar":         "unknown notify parameter",
		"notify|key=k|title=T|noequal":         "invalid notify parameter",
	}
	for actionStr, want := range cases {
		got, err := getNotifyAction(nil, "Name", actionStr)
		ts.ErrorContains(err, want, actionStr)
		ts.Nil(got)
	}
}

func (ts *NotifyActionTestSuite) Test_DispatchNotify() {
	got, err := getActionRunner([]string{"garage"}, "Name", "notify|key=k|title=T")
	ts.Nil(err)
	ts.IsType(&NotifyAction{}, got)
}

func (ts *NotifyActionTestSuite) Test_RunRaises() {
	action, _ := getNotifyAction([]string{"garage"}, "Name", "notify|key=door|severity=warning|title=Porte|body=Ouverte")
	ts.Nil(action.run())
	ts.Equal([]invokedAction{{name: "raise", input: map[string]interface{}{
		"key":      "door",
		"source":   "house",
		"severity": "warning",
		"title":    "Porte",
		"body":     "Ouverte",
		"audience": "all",
		"thing":    "garage",
	}}}, ts.thing.calls)
}

func (ts *NotifyActionTestSuite) Test_RunOmitsEmptyBodyAndThing() {
	action, _ := getNotifyAction(nil, "Name", "notify|key=door|title=Porte")
	ts.Nil(action.run())
	ts.Equal(map[string]interface{}{
		"key":      "door",
		"source":   "house",
		"severity": "info",
		"title":    "Porte",
		"audience": "all",
	}, ts.thing.calls[0].input)
}

func (ts *NotifyActionTestSuite) Test_RunInvokeError() {
	ts.thing.err = errors.New("boom")
	action, _ := getNotifyAction(nil, "Name", "notify|key=door|title=Porte")
	ts.EqualError(action.run(), "boom")
}

func (ts *NotifyActionTestSuite) Test_RunThingNotFound() {
	_consumer = &fakeNoticesConsumer{}
	action, _ := getNotifyAction(nil, "Name", "notify|key=door|title=Porte")
	ts.EqualError(action.run(), "notices thing not found")
}

func (ts *NotifyActionTestSuite) Test_Resolve() {
	action, _ := getNotifyAction(nil, "Name", "notify|key=door|title=Porte")
	ts.Nil(action.resolve())
	ts.Equal([]invokedAction{{name: "resolve", input: map[string]interface{}{"key": "door"}}}, ts.thing.calls)
}

// The notice is resolved only when the conditions go from true to false.
func (ts *NotifyActionTestSuite) Test_PlanResolvesOnFalse() {
	action, _ := getNotifyAction(nil, "Name", "notify|key=door|title=Porte")
	condition := &switchCondition{}
	a := &Automation{groups: []actionGroup{{
		action:  action,
		entries: []scheduleEntry{{conditions: []Condition{condition}, schedule: &scheduleImmediate{}}},
	}}}
	now := time.Now()

	a.plan(now) // false at start: nothing to resolve
	ts.Empty(ts.thing.calls)

	condition.on = true
	ts.Len(a.plan(now), 1)
	ts.Empty(ts.thing.calls) // planning doesn't run the action

	a.plan(now) // still true
	ts.Empty(ts.thing.calls)

	condition.on = false
	ts.Empty(a.plan(now))
	ts.Equal([]invokedAction{{name: "resolve", input: map[string]interface{}{"key": "door"}}}, ts.thing.calls)

	a.plan(now) // still false: no second resolve
	ts.Len(ts.thing.calls, 1)
}
