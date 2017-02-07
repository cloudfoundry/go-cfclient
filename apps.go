package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"
)

type AppResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []AppResource `json:"resources"`
}

type AppResource struct {
	Meta   Meta `json:"metadata"`
	Entity App  `json:"entity"`
}

type App struct {
	Guid                     string                 `json:"guid"`
	Name                     string                 `json:"name"`
	Memory                   int                    `json:"memory"`
	Instances                int                    `json:"instances"`
	DiskQuota                int                    `json:"disk_quota"`
	SpaceGuid                string                 `json:"space_guid"`
	StackGuid                string                 `json:"stack_guid"`
	State                    string                 `json:"state"`
	Command                  string                 `json:"command"`
	Buildpack                string                 `json:"buildpack"`
	DetectedBuildpack        string                 `json:"detected_buildpack"`
	DetectedBuildpackGuid    string                 `json:"detected_buildpack_guid"`
	HealthCheckHttpEndpoint  string                 `json:"health_check_http_endpoint"`
	HealthCheckType          string                 `json:"health_check_type"`
	HealthCheckTimeout       int                    `json:"health_check_timeout"`
	Diego                    bool                   `json:"diego"`
	EnableSSH                bool                   `json:"enable_ssh"`
	DetectedStartCommand     string                 `json:"detected_start_command"`
	DockerImage              string                 `json:"docker_image"`
	DockerCredentials        map[string]interface{} `json:"docker_credentials_json"`
	Environment              map[string]interface{} `json:"environment_json"`
	StagingFailedReason      string                 `json:"staging_failed_reason"`
	StagingFailedDescription string                 `json:"staging_failed_description"`
	Ports                    []int                  `json:"ports"`
	SpaceURL                 string                 `json:"space_url"`
	SpaceData                SpaceResource          `json:"space"`
	c                        *Client
}

type AppInstance struct {
	State string    `json:"state"`
	Since sinceTime `json:"since"`
}

type AppStats struct {
	State string `json:"state"`
	Stats struct {
		Name      string   `json:"name"`
		Uris      []string `json:"uris"`
		Host      string   `json:"host"`
		Port      int      `json:"port"`
		Uptime    int      `json:"uptime"`
		MemQuota  int      `json:"mem_quota"`
		DiskQuota int      `json:"disk_quota"`
		FdsQuota  int      `json:"fds_quota"`
		Usage     struct {
			Time statTime `json:"time"`
			CPU  float64  `json:"cpu"`
			Mem  int      `json:"mem"`
			Disk int      `json:"disk"`
		} `json:"usage"`
	} `json:"stats"`
}

type AppSummary struct {
	Guid             string `json:"guid"`
	Name             string `json:"name"`
	ServiceCount     int    `json:"service_count"`
	RunningInstances int    `json:"running_instances"`
	Memory           int    `json:"memory"`
	Instances        int    `json:"instances"`
	DiskQuota        int    `json:"disk_quota"`
	State            string `json:"state"`
	Diego            bool   `json:"diego"`
}

// Custom time types to handle non-RFC3339 formatting in API JSON

type sinceTime struct {
	time.Time
}

func (s *sinceTime) UnmarshalJSON(b []byte) (err error) {
	timeFlt, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}
	time := time.Unix(int64(timeFlt), 0)
	*s = sinceTime{time}
	return nil
}

func (s sinceTime) ToTime() time.Time {
	t, _ := time.Parse(time.UnixDate, s.Format(time.UnixDate))
	return t
}

type statTime struct {
	time.Time
}

func (s *statTime) UnmarshalJSON(b []byte) (err error) {
	timeString, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	time, err := time.Parse("2006-01-02 15:04:05 -0700", timeString)
	if err != nil {
		return err
	}
	*s = statTime{time}
	return nil
}

func (s statTime) ToTime() time.Time {
	t, _ := time.Parse(time.UnixDate, s.Format(time.UnixDate))
	return t
}

func (a *App) Space() (Space, error) {
	var spaceResource SpaceResource
	r := a.c.NewRequest("GET", a.SpaceURL)
	resp, err := a.c.DoRequest(r)
	if err != nil {
		return Space{}, fmt.Errorf("Error requesting space: %v", err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request: %v", err)
	}

	err = json.Unmarshal(resBody, &spaceResource)
	if err != nil {
		return Space{}, fmt.Errorf("Error unmarshaling space: %v", err)
	}
	spaceResource.Entity.Guid = spaceResource.Meta.Guid
	spaceResource.Entity.c = a.c
	return spaceResource.Entity, nil
}

func (c *Client) ListAppsByQuery(query url.Values) ([]App, error) {
	var apps []App

	requestUrl := "/v2/apps?" + query.Encode()
	for {
		var appResp AppResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, fmt.Errorf("Error requesting apps %v", err)
		}
		defer resp.Body.Close()
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading app request %v", resBody)
		}

		err = json.Unmarshal(resBody, &appResp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling app %v", err)
		}

		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
			app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
			app.Entity.c = c
			apps = append(apps, app.Entity)
		}

		requestUrl = appResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return apps, nil
}

func (c *Client) ListApps() ([]App, error) {
	q := url.Values{}
	q.Set("inline-relations-depth", "2")
	return c.ListAppsByQuery(q)
}

func (c *Client) GetAppInstances(guid string) (map[string]AppInstance, error) {
	var appInstances map[string]AppInstance

	requestURL := fmt.Sprintf("/v2/apps/%s/instances", guid)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting app instances %v", err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading app instances %v", err)
	}
	err = json.Unmarshal(resBody, &appInstances)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling app instances %v", err)
	}
	return appInstances, nil
}

func (c *Client) GetAppStats(guid string) (map[string]AppStats, error) {
	var appStats map[string]AppStats

	requestURL := fmt.Sprintf("/v2/apps/%s/stats", guid)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting app stats %v", err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading app stats %v", err)
	}
	err = json.Unmarshal(resBody, &appStats)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling app stats %v", err)
	}
	return appStats, nil
}

func (c *Client) KillAppInstance(guid string, index string) error {
	requestURL := fmt.Sprintf("/v2/apps/%s/instances/%s", guid, index)
	r := c.NewRequest("DELETE", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return fmt.Errorf("Error stopping app %s at index %s", guid, index)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return fmt.Errorf("Error stopping app %s at index %s", guid, index)
	}
	return nil
}

func (c *Client) AppByGuid(guid string) (App, error) {
	var appResource AppResource
	r := c.NewRequest("GET", "/v2/apps/"+guid+"?inline-relations-depth=2")
	resp, err := c.DoRequest(r)
	if err != nil {
		return App{}, fmt.Errorf("Error requesting apps: %v", err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading app request %v", resBody)
	}

	err = json.Unmarshal(resBody, &appResource)
	if err != nil {
		return App{}, fmt.Errorf("Error unmarshaling app: %v", err)
	}
	appResource.Entity.Guid = appResource.Meta.Guid
	appResource.Entity.SpaceData.Entity.Guid = appResource.Entity.SpaceData.Meta.Guid
	appResource.Entity.SpaceData.Entity.OrgData.Entity.Guid = appResource.Entity.SpaceData.Entity.OrgData.Meta.Guid
	appResource.Entity.c = c
	return appResource.Entity, nil
}
