package automations

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"
)

const webhookPrefix = "webhook|"

var _webhookClient = &http.Client{Timeout: 10 * time.Second}

// WebhookAction calls an HTTP(S) URL when the automation fires.
//
// Syntax: `webhook|<url>|method=<GET|POST|PUT>|body=<content>`
//   - method defaults to POST
//   - body is sent as application/json (only for POST/PUT)
type WebhookAction struct {
	AutomationName string
	URL            string
	Method         string
	Body           string
}

func getWebhookAction(automationName string, actionStr string) (*WebhookAction, error) {
	actionArray := strings.Split(actionStr, "|")
	if len(actionArray) < 2 || strings.TrimSpace(actionArray[1]) == "" {
		return nil, errors.New("missing webhook url")
	}
	rawURL := strings.TrimSpace(actionArray[1])
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, fmt.Errorf("invalid webhook url %q", rawURL)
	}

	action := &WebhookAction{
		AutomationName: automationName,
		URL:            rawURL,
		Method:         http.MethodPost,
	}
	for i := 2; i < len(actionArray); i++ {
		// SplitN, not Split: the value may itself contain '=' (JSON body, query...)
		param := strings.SplitN(actionArray[i], "=", 2)
		if len(param) != 2 {
			return nil, fmt.Errorf("invalid webhook parameter %q", actionArray[i])
		}
		switch param[0] {
		case "method":
			method := strings.ToUpper(strings.TrimSpace(param[1]))
			if method != http.MethodGet && method != http.MethodPost && method != http.MethodPut {
				return nil, fmt.Errorf("unsupported webhook method %q", param[1])
			}
			action.Method = method
		case "body":
			action.Body = param[1]
		default:
			return nil, fmt.Errorf("unknown webhook parameter %q", param[0])
		}
	}
	return action, nil
}

func (a *WebhookAction) run() error {
	var body *strings.Reader
	if a.Body != "" && a.Method != http.MethodGet {
		body = strings.NewReader(a.Body)
	} else {
		body = strings.NewReader("")
	}
	req, err := http.NewRequest(a.Method, a.URL, body)
	if err != nil {
		return err
	}
	if a.Body != "" && a.Method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := _webhookClient.Do(req)
	if err != nil {
		zlog.Error().Err(err).Str("automation", a.AutomationName).Str("url", a.URL).Msg("[automations:WebhookAction] Webhook call failed")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		err := fmt.Errorf("webhook returned status %d", resp.StatusCode)
		zlog.Error().Err(err).Str("automation", a.AutomationName).Str("url", a.URL).Msg("[automations:WebhookAction] Webhook call failed")
		return err
	}
	zlog.Info().Str("automation", a.AutomationName).Str("url", a.URL).Int("status", resp.StatusCode).Msg("[automations:WebhookAction] Webhook called")
	return nil
}
