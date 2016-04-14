package cfclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

const (
	//AppCrash app.crash event const
	AppCrash = "app.crash"
	//AppStart audit.app.start event const
	AppStart = "audit.app.start"
	//AppStop audit.app.stop event const
	AppStop = "audit.app.stop"
	//AppUpdate audit.app.update event const
	AppUpdate = "audit.app.update"
	//AppCreate audit.app.create event const
	AppCreate = "audit.app.create"
	//AppDelete audit.app.delete-request event const
	AppDelete = "audit.app.delete-request"
	//AppSSHAuth audit.app.ssh-authorized event const
	AppSSHAuth = "audit.app.ssh-authorized"
	//AppSSHUnauth audit.app.ssh-unauthorized event const
	AppSSHUnauth = "audit.app.ssh-unauthorized"
	//FilterTimestamp const for query filter timestamp
	FilterTimestamp = "timestamp"
	//FilterActee const for query filter actee
	FilterActee = "actee"
)

//ValidOperators const for all valid operators in a query
var ValidOperators = []string{":", ">=", "<=", "<", ">", "IN"}

// AppEventResponse the entire response
type AppEventResponse struct {
	Results   int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	PrevURL   string             `json:"prev_url"`
	NextURL   string             `json:"next_url"`
	Resources []AppEventResource `json:"resources"`
}

// AppEventResource the event resources
type AppEventResource struct {
	Meta   Meta           `json:"metadata"`
	Entity AppEventEntity `json:"entity"`
}

//AppEventQuery a struct for defining queries like 'q=filter>value' or 'q=filter IN a,b,c'
type AppEventQuery struct {
	Filter   string
	Operator string
	Value    string
}

// The AppEventEntity the actual app event body
type AppEventEntity struct {
	//EventTypes are app.crash, audit.app.start, audit.app.stop, audit.app.update, audit.app.create, audit.app.delete-request
	EventType string `json:"type"`
	//The GUID of the actor.
	Actor string `json:"actor"`
	//The actor type, user or app
	ActorType string `json:"actor_type"`
	//The name of the actor.
	ActorName string `json:"actor_name"`
	//The GUID of the actee.
	Actee string `json:"actee"`
	//The actee type, space, app or v3-app
	ActeeType string `json:"actee_type"`
	//The name of the actee.
	ActeeName string `json:"actee_name"`
	//Timestamp format "2016-02-26T13:29:44Z". The event creation time.
	Timestamp time.Time `json:"timestamp"`
	MetaData  struct {
		Request struct {
			Name              string  `json:"name,omitempty"`
			Instances         float64 `json:"instances,omitempty"`
			State             string  `json:"state,omitempty"`
			Memory            float64 `json:"memory,omitempty"`
			EnvironmentVars   string  `json:"environment_json,omitempty"`
			DockerCredentials string  `json:"docker_credentials_json,omitempty"`
			//audit.app.create event fields
			Console            bool    `json:"console,omitempty"`
			Buildpack          string  `json:"buildpack,omitempty"`
			Space              string  `json:"space_guid,omitempty"`
			HealthcheckType    string  `json:"health_check_type,omitempty"`
			HealthcheckTimeout float64 `json:"health_check_timeout,omitempty"`
			Production         bool    `json:"production,omitempty"`
			//app.crash event fields
			Index           float64 `json:"index,omitempty"`
			ExitStatus      string  `json:"exit_status,omitempty"`
			ExitDescription string  `json:"exit_description,omitempty"`
			ExitReason      string  `json:"reason,omitempty"`
		} `json:"request"`
	} `json:"metadata"`
}

// ListAppEvents returns all app events based on eventType
func (c *Client) ListAppEvents(eventType string) ([]AppEventEntity, error) {
	return c.ListAppEventsByQuery(eventType, nil)
}

// ListAppEventsByQuery returns all app events based on eventType and queries
func (c *Client) ListAppEventsByQuery(eventType string, queries []AppEventQuery) ([]AppEventEntity, error) {

	if eventType != AppCrash && eventType != AppStart && eventType != AppStop && eventType != AppUpdate && eventType != AppCreate &&
		eventType != AppDelete && eventType != AppSSHAuth && eventType != AppSSHUnauth {
		return nil, errors.New("Unsupported app event type " + eventType)
	}

	var query = "/v2/events?q=type:" + eventType
	//adding the additional queries
	if queries != nil && len(queries) > 0 {
		for _, eventQuery := range queries {
			if eventQuery.Filter != FilterTimestamp && eventQuery.Filter != FilterActee {
				return nil, errors.New("Unsupported query filter type " + eventQuery.Filter)
			}
			if !stringInSlice(eventQuery.Operator, ValidOperators) {
				return nil, errors.New("Unsupported query operator type " + eventQuery.Operator)
			}
			query += "&q=" + eventQuery.Filter + eventQuery.Operator + eventQuery.Value
		}
	}

	requ := c.NewRequest("GET", query)

	resp, err := c.DoRequest(requ)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(string(resBody[:]))
	}

	var eventResponse AppEventResponse
	err = json.Unmarshal(resBody, &eventResponse)
	if err != nil {
		return nil, err
	}

	eventsLen := len(eventResponse.Resources)
	events := make([]AppEventEntity, eventsLen)
	for i := 0; i < eventsLen; i++ {
		events[i] = eventResponse.Resources[i].Entity
	}
	return events, nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
