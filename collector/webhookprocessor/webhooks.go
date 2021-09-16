package webhookprocessor

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

type Event string

const (
	UnknownEvent                Event = "unknown_event"
	GithubDeploymentStatusEvent Event = "deployment_status"
	PagerDutyActiveIncident     Event = "pagerduty_incident"
)

var (
	ErrInvalidHTTPMethod = errors.New("invalid HTTP Method")
	ErrParsingPayload    = errors.New("error parsing payload")
)

func (h *httpServer) parseWebhook(r *http.Request) (*gabs.Container, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, ErrParsingPayload
	}

	jsonParsed, err := gabs.ParseJSON(payload)
	if err != nil {
		return nil, ErrParsingPayload
	}

	return jsonParsed, nil
}

func (h *httpServer) webhookHandler(events ...Event) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonParsed, err := h.parseWebhook(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "bad request: invalid webhook: %s", err)
			return
		}
		actionType := UnknownEvent

		// TODO: Test for User Agents
		event := r.Header.Get("X-GitHub-Event")
		if event == "deployment_status" {
			deployStatus, ok := jsonParsed.Search("deployment_status", "state").Data().(string)
			deployId, idOk := jsonParsed.Search("deployment_status", "id").Data().(string)

			if ok && idOk {
				actionType = GithubDeploymentStatusEvent
				if deployStatus == "pending" {
					h.addAttribute("github.com.active_deployment", deployId, "")
				} else {
					h.removeAttribute("github.com.active_deployment", "")
				}
			}
		}

		gremlinAttackId, ok := jsonParsed.Search("attackId").Data().(string)
		gremlinAttackStatus, okStatus := jsonParsed.Search("attackStatus").Data().(string)
		if ok && okStatus {
			if gremlinAttackStatus == "RUNNING" {
				h.addAttribute("gremlin.com.active_attack", gremlinAttackId, "")
			} else if gremlinAttackStatus == "FINISHED" {
				h.removeAttribute("gremlin.com.active_attack", "")
			}
		}
		pagerdutyEvent, ok := jsonParsed.Search("messages", "0", "event").Data().(string)
		if ok {
			incident, ok := jsonParsed.Search("messages", "0", "incident", "incident_number").Data().(float64)
			if ok {
				if pagerdutyEvent == "incident.trigger" {
					actionType = PagerDutyActiveIncident
					h.addAttribute("pagerduty.com.has_incident", "true", "")
					h.addAttribute("pagerduty.com.active_incident", fmt.Sprintf("%v", incident), "")
				} else if pagerdutyEvent == "incident.resolve" {
					actionType = PagerDutyActiveIncident
					h.removeAttribute("pagerduty.com.has_incident", "")
					h.removeAttribute("pagerduty.com.active_incident", "")
				}
			}
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, "ok: %v", actionType)
	}
}
