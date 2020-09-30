package cfclient

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

type Droplet struct {
	Guid      string `json:"guid"`
	State     string `json:"state"`
	Error     string `json:"error"`
	Lifecycle struct {
		Type string   `json:"type"`
		Data struct{} `json:"data"`
	} `json:"lifecycle"`
	ExecutionMetadata string   `json:"execution_metadata"`
	ProcessTypes      struct{} `json:"process_types"`
	Checksum          struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"checksum"`
	Buildpacks []struct {
		Name          string `json:"name"`
		BuildpackName string `json:"buildpack_name"`
		Version       string `json:"version"`
	} `json:"buildpacks"`
	Stack         string    `json:"stack"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Relationships struct {
		App struct {
			Data struct {
				Guid string `json:"guid"`
			} `json:"data"`
		} `json:"app"`
	} `json:"relationships"`
	Links struct {
		Self    Link `json:"self"`
		App     Link `json:"app"`
		Droplet Link `json:"droplet"`
	} `json:"links"`
}

func (c *Client) GetAppsDropletsCurrent(appGuid string) (Droplet, error) {
	var droplet Droplet
	r := c.NewRequest("GET", "/v3/apps/"+appGuid+"/droplets/current")
	resp, err := c.DoRequest(r)
	if err != nil {
		return Droplet{}, errors.Wrap(err, "Error requesting droplets")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Droplet{}, errors.Wrap(err, "Error reading app response body")
	}

	err = json.Unmarshal(resBody, &droplet)
	if err != nil {
		return Droplet{}, errors.Wrap(err, "Error unmarshalling app")
	}
	return droplet, nil
}
