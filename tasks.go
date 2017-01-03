package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

// TaskListResponse is the JSON response from the API.
type TaskListResponse struct {
	Pagination struct {
		TotalResults int `json:"total_results"`
		TotalPages   int `json:"total_pages"`
		First        struct {
			Href string `json:"href"`
		} `json:"first"`
		Last struct {
			Href string `json:"href"`
		} `json:"last"`
		Next     interface{} `json:"next"`
		Previous interface{} `json:"previous"`
	} `json:"pagination"`
	Tasks []Task `json:"resources"`
}

// Task is a description of a task element.
type Task struct {
	GUID       string `json:"guid"`
	SequenceID int    `json:"sequence_id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	State      string `json:"state"`
	MemoryInMb int    `json:"memory_in_mb"`
	DiskInMb   int    `json:"disk_in_mb"`
	Result     struct {
		FailureReason string `json:"failure_reason"`
	} `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DropletGUID string    `json:"droplet_guid"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		App struct {
			Href string `json:"href"`
		} `json:"app"`
		Droplet struct {
			Href string `json:"href"`
		} `json:"droplet"`
	} `json:"links"`
}

// TaskRequest is a v3 JSON object as described in:
// http://v3-apidocs.cloudfoundry.org/version/3.0.0/index.html#create-a-task
type TaskRequest struct {
	Command          string `json:"command"`
	Name             string `json:"name"`
	MemoryInMegabyte int    `json:"memory_in_mb"`
	DiskInMegabyte   int    `json:"disk_in_mb"`
	DropletGUID      string `json:"droplet_guid"`
}

func (c *Client) makeTaskListRequest() ([]byte, error) {
	req := c.NewRequest("GET", "/v3/tasks")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("Error requesting tasks %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error requesting tasks: status code not 200, it was %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func parseTaskListRespones(answer []byte) (TaskListResponse, error) {
	var response TaskListResponse
	errUnmarshal := json.Unmarshal(answer, &response)
	if errUnmarshal != nil {
		return response, fmt.Errorf("Error unmarshaling response %v", errUnmarshal)
	}
	return response, nil
}

// ListTasks returns all tasks the user has access to.
func (c *Client) ListTasks() ([]Task, error) {
	body, err := c.makeTaskListRequest()
	if err != nil {
		return nil, fmt.Errorf("Error requesting tasks %v", err)
	}
	response, errParse := parseTaskListRespones(body)
	if errParse != nil {
		return nil, fmt.Errorf("Error reading tasks %v", response)
	}
	return response.Tasks, nil
}

func createReader(tr TaskRequest) (io.Reader, error) {
	rmap := make(map[string]string)
	rmap["command"] = tr.Command
	if tr.Name != "" {
		rmap["name"] = tr.Name
	}
	// setting droplet GUID causing issues
	if tr.MemoryInMegabyte != 0 {
		rmap["memory_in_mb"] = fmt.Sprintf("%d", tr.MemoryInMegabyte)
	}
	if tr.DiskInMegabyte != 0 {
		rmap["disk_in_mb"] = fmt.Sprintf("%d", tr.DiskInMegabyte)
	}

	bodyReader := bytes.NewBuffer(nil)
	enc := json.NewEncoder(bodyReader)
	if err := enc.Encode(rmap); err != nil {
		return nil, fmt.Errorf("Error during encoding task request %v", err)
	}
	return bodyReader, nil
}

// CreateTask creates a new task in CF system and returns its structure.
func (c *Client) CreateTask(tr TaskRequest) (task Task, err error) {
	bodyReader, err := createReader(tr)
	if err != nil {
		return task, err
	}

	request := fmt.Sprintf("/v3/apps/%s/tasks", tr.DropletGUID)
	req := c.NewRequestWithBody("POST", request, bodyReader)

	resp, errReq := c.DoRequest(req)
	if errReq != nil {
		return task, fmt.Errorf("Error creating task %v", errReq)
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		return task, fmt.Errorf("Error reading task after creation %v", errBody)
	}

	errUnmarshal := json.Unmarshal(body, &task)
	if errUnmarshal != nil {
		return task, fmt.Errorf("Error unmarshaling task %v", errUnmarshal)
	}
	return task, errUnmarshal
}

// TaskByGuid returns a task structure by requesting it with the tasks GUID.
func (c *Client) TaskByGuid(guid string) (task Task, err error) {
	request := fmt.Sprintf("/v3/tasks/%s", guid)
	req := c.NewRequest("GET", request)

	resp, errReq := c.DoRequest(req)
	if errReq != nil {
		return task, fmt.Errorf("Error requesting task %v", errReq)
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		return task, fmt.Errorf("Error reading task %v", errBody)
	}

	errUnmarshal := json.Unmarshal(body, &task)
	if errUnmarshal != nil {
		return task, fmt.Errorf("Error unmarshaling task %v", errUnmarshal)
	}
	return task, errUnmarshal
}

// TasksByApp retuns task structures which aligned to an app identified by the given guid.
func (c *Client) TasksByApp(guid string) ([]Task, error) {
	request := fmt.Sprintf("/v3/apps/%s/tasks", guid)
	req := c.NewRequest("GET", request)

	resp, errReq := c.DoRequest(req)
	if errReq != nil {
		return nil, fmt.Errorf("Error requesting task %v", errReq)
	}
	defer resp.Body.Close()

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		return nil, fmt.Errorf("Error reading tasks %v", errBody)
	}

	response, errParse := parseTaskListRespones(body)
	if errParse != nil {
		return nil, fmt.Errorf("Error parsing tasks %v", response)
	}
	return response.Tasks, nil
}

// TerminateTask cancels a task identified by its GUID.
func (c *Client) TerminateTask(guid string) error {
	req := c.NewRequest("PUT", fmt.Sprintf("/v3/tasks/%s/cancel", guid))
	resp, err := c.DoRequest(req)
	if err != nil {
		return fmt.Errorf("Error terminating task %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("Failed terminating task, response status code %d", resp.StatusCode)
	}
	return nil
}
