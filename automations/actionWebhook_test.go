package automations

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type WebhookActionTestSuite struct {
	suite.Suite
}

func Test_WebhookActionTestSuite(t *testing.T) {
	suite.Run(t, &WebhookActionTestSuite{})
}

func (ts *WebhookActionTestSuite) SetupTest() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func (ts *WebhookActionTestSuite) Test_ParseMinimal() {
	got, err := getWebhookAction("Name", "webhook|http://example.com/hook")
	ts.Nil(err)
	ts.Equal(&WebhookAction{
		AutomationName: "Name",
		URL:            "http://example.com/hook",
		Method:         http.MethodPost,
	}, got)
}

func (ts *WebhookActionTestSuite) Test_ParseWithParams() {
	got, err := getWebhookAction("Name", `webhook|https://example.com/hook?a=1|method=get|body={"on":true}`)
	ts.Nil(err)
	ts.Equal(&WebhookAction{
		AutomationName: "Name",
		URL:            "https://example.com/hook?a=1",
		Method:         http.MethodGet,
		Body:           `{"on":true}`,
	}, got)
}

func (ts *WebhookActionTestSuite) Test_ParseMissingURL() {
	got, err := getWebhookAction("Name", "webhook|")
	ts.EqualError(err, "missing webhook url")
	ts.Nil(got)
}

func (ts *WebhookActionTestSuite) Test_ParseInvalidURL() {
	got, err := getWebhookAction("Name", "webhook|not-an-url")
	ts.ErrorContains(err, "invalid webhook url")
	ts.Nil(got)
}

func (ts *WebhookActionTestSuite) Test_ParseInvalidMethod() {
	got, err := getWebhookAction("Name", "webhook|http://example.com|method=DELETE")
	ts.ErrorContains(err, "unsupported webhook method")
	ts.Nil(got)
}

func (ts *WebhookActionTestSuite) Test_ParseUnknownParam() {
	got, err := getWebhookAction("Name", "webhook|http://example.com|foo=bar")
	ts.ErrorContains(err, "unknown webhook parameter")
	ts.Nil(got)
}

func (ts *WebhookActionTestSuite) Test_RunPostWithBody() {
	var gotMethod, gotBody, gotContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		gotContentType = r.Header.Get("Content-Type")
	}))
	defer server.Close()

	action, err := getWebhookAction("Name", "webhook|"+server.URL+`|body={"on":true}`)
	ts.Nil(err)
	ts.Nil(action.run())
	ts.Equal(http.MethodPost, gotMethod)
	ts.Equal(`{"on":true}`, gotBody)
	ts.Equal("application/json", gotContentType)
}

func (ts *WebhookActionTestSuite) Test_RunGet() {
	var gotMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
	}))
	defer server.Close()

	action, err := getWebhookAction("Name", "webhook|"+server.URL+"|method=GET")
	ts.Nil(err)
	ts.Nil(action.run())
	ts.Equal(http.MethodGet, gotMethod)
}

func (ts *WebhookActionTestSuite) Test_RunErrorStatus() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	action, err := getWebhookAction("Name", "webhook|"+server.URL)
	ts.Nil(err)
	ts.ErrorContains(action.run(), "webhook returned status 500")
}

func (ts *WebhookActionTestSuite) Test_RunUnreachable() {
	action, err := getWebhookAction("Name", "webhook|http://127.0.0.1:1/hook")
	ts.Nil(err)
	ts.NotNil(action.run())
}

func (ts *WebhookActionTestSuite) Test_DispatchWebhook() {
	got, err := getActionRunner(nil, "Name", "webhook|http://example.com/hook")
	ts.Nil(err)
	ts.IsType(&WebhookAction{}, got)
}
