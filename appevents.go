package cfclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// AppEvent the detailed event type for app events
type AppEvent struct {
	EventResponse
}

// ListAppCreateEvent returns all app creation events
func (c *Client) ListAppCreateEvent() ([]Event, error) {
	requ := c.newRequest("GET", "/v2/events?q=type:audit.app.create")

	resp, err := c.doRequest(requ)
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

	var eventResponse EventResponse
	err = json.Unmarshal(resBody, &eventResponse)
	if err != nil {
		return nil, err
	}

	eventsLen := len(eventResponse.Resources)
	events := make([]Event, eventsLen)
	for i := 0; i < eventsLen; i++ {
		events[i] = eventResponse.Resources[i].Entity
	}
	return events, nil
}
