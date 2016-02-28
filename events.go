package cfclient

import (
	"time"
)

// EventInterface an interface that can be implemented by subevents like e.g. AppEvent.
type EventInterface interface {
	GetEventResponse() EventResponse
	SetEventResponse(response EventResponse)
}

// EventResponse is the header of the response e.g.
//  "total_results": 29,
//   "total_pages": 1,
//   "prev_url": null,
//   "next_url": null,
//   "resources": ...
type EventResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextURL   string          `json:"next_url"`
	Resources []EventResource `json:"resources"`
}

// EventResource the event resource e.g.
// "metadata": {...},
// "entity":{ ...}
type EventResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Event `json:"entity"`
}

// Event the actual event
type Event struct {
	EventType string `json:"type"`
	Actor     string `json:"actor"`
	ActorType string `json:"actor_type"`
	ActorName string `json:"actor_name"`
	Actee     string `json:"actee"`
	ActeeType string `json:"actee_type"`
	ActeeName string `json:"actee_name"`
	//Timestamp format "2016-02-26T13:29:44Z"
	Timestamp time.Time `json:"timestamp":`
	// "metadata": {
	// "request": {
	// "memory": 2048,
	// "instances": 1,
	// "name": "email-dev-v4-v4-11-1-b217",
	// "state": "STOPPED",
	// "space_guid": "0baa90dc-3aad-402a-a251-b127c0e1d1a3",
	// "buildpack": "java_buildpack",
	// "console": false,
	// "docker_credentials_json": "PRIVATE DATA HIDDEN",
	// "environment_json": "PRIVATE DATA HIDDEN",
	// "health_check_type": "port",
	// "production": false

}
