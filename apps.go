package cfclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry/gofileutils/fileutils"
	"github.com/pkg/errors"
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

type AppRequest struct {
	Name      string `json:"name"`
	SpaceGuid string `json:"space_guid"`
}

type App struct {
	Guid                     string                 `json:"guid"`
	CreatedAt                string                 `json:"created_at"`
	UpdatedAt                string                 `json:"updated_at"`
	Name                     string                 `json:"name"`
	Memory                   int                    `json:"memory"`
	Instances                int                    `json:"instances"`
	DiskQuota                int                    `json:"disk_quota"`
	SpaceGuid                string                 `json:"space_guid"`
	StackGuid                string                 `json:"stack_guid"`
	State                    string                 `json:"state"`
	PackageState             string                 `json:"package_state"`
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
	PackageUpdatedAt         string                 `json:"package_updated_at"`
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
	Guid                     string                 `json:"guid"`
	Name                     string                 `json:"name"`
	ServiceCount             int                    `json:"service_count"`
	RunningInstances         int                    `json:"running_instances"`
	SpaceGuid                string                 `json:"space_guid"`
	StackGuid                string                 `json:"stack_guid"`
	Buildpack                string                 `json:"buildpack"`
	DetectedBuildpack        string                 `json:"detected_buildpack"`
	Environment              map[string]interface{} `json:"environment_json"`
	Memory                   int                    `json:"memory"`
	Instances                int                    `json:"instances"`
	DiskQuota                int                    `json:"disk_quota"`
	State                    string                 `json:"state"`
	Command                  string                 `json:"command"`
	PackageState             string                 `json:"package_state"`
	HealthCheckType          string                 `json:"health_check_type"`
	HealthCheckTimeout       int                    `json:"health_check_timeout"`
	StagingFailedReason      string                 `json:"staging_failed_reason"`
	StagingFailedDescription string                 `json:"staging_failed_description"`
	Diego                    bool                   `json:"diego"`
	DockerImage              string                 `json:"docker_image"`
	DetectedStartCommand     string                 `json:"detected_start_command"`
	EnableSSH                bool                   `json:"enable_ssh"`
	DockerCredentials        map[string]interface{} `json:"docker_credentials_json"`
}

type AppEnv struct {
	// These can have arbitrary JSON so need to map to interface{}
	Environment    map[string]interface{} `json:"environment_json"`
	StagingEnv     map[string]interface{} `json:"staging_env_json"`
	RunningEnv     map[string]interface{} `json:"running_env_json"`
	SystemEnv      map[string]interface{} `json:"system_env_json"`
	ApplicationEnv map[string]interface{} `json:"application_env_json"`
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

	possibleFormats := [...]string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05 -0700", "2006-01-02 15:04:05 MST"}

	var value time.Time
	for _, possibleFormat := range possibleFormats {
		if value, err = time.Parse(possibleFormat, timeString); err == nil {
			*s = statTime{value}
			return nil
		}
	}

	return fmt.Errorf("%s was not in any of the expected Date Formats %v", timeString, possibleFormats)
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
		return Space{}, errors.Wrap(err, "Error requesting space")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Space{}, errors.Wrap(err, "Error reading space response")
	}

	err = json.Unmarshal(resBody, &spaceResource)
	if err != nil {
		return Space{}, errors.Wrap(err, "Error unmarshalling body")
	}
	return a.c.mergeSpaceResource(spaceResource), nil
}

// ListAppsByQueryWithLimits queries totalPages app info. When totalPages is
// less and equal than 0, it queries all app info
// When there are no more than totalPages apps on server side, all apps info will be returned
func (c *Client) ListAppsByQueryWithLimits(query url.Values, totalPages int) ([]App, error) {
	return c.listApps("/v2/apps?"+query.Encode(), totalPages)
}

func (c *Client) ListAppsByQuery(query url.Values) ([]App, error) {
	return c.listApps("/v2/apps?"+query.Encode(), -1)
}

// GetAppByGuidNoInlineCall will fetch app info including space and orgs information
// Without using inline-relations-depth=2 call
func (c *Client) GetAppByGuidNoInlineCall(guid string) (App, error) {
	var appResource AppResource
	r := c.NewRequest("GET", "/v2/apps/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return App{}, errors.Wrap(err, "Error requesting apps")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return App{}, errors.Wrap(err, "Error reading app response body")
	}

	err = json.Unmarshal(resBody, &appResource)
	if err != nil {
		return App{}, errors.Wrap(err, "Error unmarshalling app")
	}
	app := c.mergeAppResource(appResource)

	// If no Space Information no need to check org.
	if app.SpaceGuid != "" {
		//Getting Spaces Resource
		space, err := app.Space()
		if err != nil {
			errors.Wrap(err, "Unable to get the Space for the apps "+app.Name)
		} else {
			app.SpaceData.Entity = space

		}

		//Getting orgResource
		org, err := app.SpaceData.Entity.Org()
		if err != nil {
			errors.Wrap(err, "Unable to get the Org for the apps "+app.Name)
		} else {
			app.SpaceData.Entity.OrgData.Entity = org
		}
	}

	return app, nil
}

func (c *Client) ListApps() ([]App, error) {
	q := url.Values{}
	q.Set("inline-relations-depth", "2")
	return c.ListAppsByQuery(q)
}

func (c *Client) ListAppsByRoute(routeGuid string) ([]App, error) {
	return c.listApps(fmt.Sprintf("/v2/routes/%s/apps", routeGuid), -1)
}

func (c *Client) listApps(requestUrl string, totalPages int) ([]App, error) {
	pages := 0
	apps := []App{}
	for {
		var appResp AppResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)

		if err != nil {
			return nil, errors.Wrap(err, "Error requesting apps")
		}
		defer resp.Body.Close()
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading app request")
		}

		err = json.Unmarshal(resBody, &appResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshalling app")
		}
		for _, app := range appResp.Resources {
			apps = append(apps, c.mergeAppResource(app))
		}

		requestUrl = appResp.NextUrl
		if requestUrl == "" {
			break
		}

		pages += 1
		if totalPages > 0 && pages >= totalPages {
			break
		}
	}
	return apps, nil
}

func (c *Client) GetAppInstances(guid string) (map[string]AppInstance, error) {
	var appInstances map[string]AppInstance

	requestURL := fmt.Sprintf("/v2/apps/%s/instances", guid)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting app instances")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading app instances")
	}
	err = json.Unmarshal(resBody, &appInstances)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling app instances")
	}
	return appInstances, nil
}

func (c *Client) GetAppEnv(guid string) (AppEnv, error) {
	var appEnv AppEnv

	requestURL := fmt.Sprintf("/v2/apps/%s/env", guid)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return appEnv, errors.Wrap(err, "Error requesting app env")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return appEnv, errors.Wrap(err, "Error reading app env")
	}
	err = json.Unmarshal(resBody, &appEnv)
	if err != nil {
		return appEnv, errors.Wrap(err, "Error unmarshalling app env")
	}
	return appEnv, nil
}

func (c *Client) GetAppRoutes(guid string) ([]Route, error) {
	return c.fetchRoutes(fmt.Sprintf("/v2/apps/%s/routes", guid))
}

func (c *Client) GetAppStats(guid string) (map[string]AppStats, error) {
	var appStats map[string]AppStats

	requestURL := fmt.Sprintf("/v2/apps/%s/stats", guid)
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting app stats")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading app stats")
	}
	err = json.Unmarshal(resBody, &appStats)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling app stats")
	}
	return appStats, nil
}

func (c *Client) KillAppInstance(guid string, index string) error {
	requestURL := fmt.Sprintf("/v2/apps/%s/instances/%s", guid, index)
	r := c.NewRequest("DELETE", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return errors.Wrapf(err, "Error stopping app %s at index %s", guid, index)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return errors.Wrapf(err, "Error stopping app %s at index %s", guid, index)
	}
	return nil
}

func (c *Client) GetAppByGuid(guid string) (App, error) {
	var appResource AppResource
	r := c.NewRequest("GET", "/v2/apps/"+guid+"?inline-relations-depth=2")
	resp, err := c.DoRequest(r)
	if err != nil {
		return App{}, errors.Wrap(err, "Error requesting apps")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return App{}, errors.Wrap(err, "Error reading app response body")
	}

	err = json.Unmarshal(resBody, &appResource)
	if err != nil {
		return App{}, errors.Wrap(err, "Error unmarshalling app")
	}
	return c.mergeAppResource(appResource), nil
}

func (c *Client) AppByGuid(guid string) (App, error) {
	return c.GetAppByGuid(guid)
}

//AppByName takes an appName, and GUIDs for a space and org, and performs
// the API lookup with those query parameters set to return you the desired
// App object.
func (c *Client) AppByName(appName, spaceGuid, orgGuid string) (app App, err error) {
	query := url.Values{}
	query.Add("q", fmt.Sprintf("organization_guid:%s", orgGuid))
	query.Add("q", fmt.Sprintf("space_guid:%s", spaceGuid))
	query.Add("q", fmt.Sprintf("name:%s", appName))
	apps, err := c.ListAppsByQuery(query)
	if err != nil {
		return
	}
	if len(apps) == 0 {
		err = fmt.Errorf("No app found with name: `%s` in space with GUID `%s` and org with GUID `%s`", appName, spaceGuid, orgGuid)
		return
	}
	app = apps[0]
	return
}

//Upload app bits
func (c *Client) Upload(file io.Reader, guid string) error {
	var capturedErr error
	fileutils.TempFile("requests", func(requestFile *os.File, err error) {
		if err != nil {
			log.Fatal(err)
		}

		writer := multipart.NewWriter(requestFile)
		// err = writer.WriteField("async", "true")
		// if err != nil {
		//     capturedErr = fmt.Errorf("Error writing resource field %s", err)
		// }

		err = writer.WriteField("resources", "[]")
		if err != nil {
			capturedErr = fmt.Errorf("Error writing resource field %s", err)
		}

		part, err := writer.CreateFormFile("application", "application.zip")
		if err != nil {
			_ = writer.Close()
			capturedErr = err
			log.Fatal(err)

		}

		_, err = io.Copy(part, file)
		if err != nil {
			capturedErr = fmt.Errorf("Error creating upload: %s", err.Error())
			log.Fatal(err)

		}

		err = writer.Close()
		if err != nil {
			capturedErr = err
			log.Fatal(err)

		}

		requestFile.Seek(0, 0)
		fileStats, err := requestFile.Stat()
		if err != nil {
			capturedErr = fmt.Errorf("Error getting file info: %s", err)
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/apps/%s/bits", c.Config.ApiAddress, guid), requestFile)
		if err != nil {
			capturedErr = err
			log.Fatal(err)
		}

		req.ContentLength = fileStats.Size()
		contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())
		req.Header.Set("Content-Type", contentType)

		resp, err := c.Do(req) //client.Do() handles the HTTP status code checking for us
		if err != nil {
			capturedErr = err
			log.Fatal(err)
		}
		defer resp.Body.Close()
		log.Printf("Response %s", resp.Status)
	})

	return errors.Wrap(capturedErr, "Error uploading application bits:")
}

// GetAppBits downloads the application's bits as a tar file
func (c *Client) GetAppBits(guid string) (io.ReadCloser, error) {
	requestURL := fmt.Sprintf("/v2/apps/%s/download", guid)
	httpRequest, _ := c.NewRequest("GET", requestURL).toHTTP()
	httpRequest.Header.Set("Content-type", "application/zip")

	// temporarily don't follow redirects
	prevCheckRedirect := c.Config.HttpClient.CheckRedirect
	c.Config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		c.Config.HttpClient.CheckRedirect = prevCheckRedirect
	}()

	resp, err := c.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Error downloading app bits")
	}
	if isRedirect(resp) {
		// directly download the bits using a non cloud controller transport
		blobStoreLocation := resp.Header.Get("Location")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Config.SkipSslValidation},
		}
		client := &http.Client{Transport: tr}
		resp, err = client.Get(blobStoreLocation)
		if err != nil {
			return nil, errors.Wrap(err, "Error downloading app bits from blobstore")
		}
	} else {
		return nil, errors.Wrap(err, "Error downloading app bits, expected a redirect to a blobstore")
	}
	return resp.Body, nil
}

func (c *Client) CreateApp(name, spaceGuid string) (App, error) {
	var appResp AppResource
	app := &AppRequest{
		Name:      name,
		SpaceGuid: spaceGuid,
	}
	appBytes, err := json.Marshal(app)
	if err != nil {
		return App{}, err
	}
	r := c.NewRequestWithBody("POST", "/v2/apps", bytes.NewReader(appBytes))
	resp, err := c.DoRequest(r)
	if err != nil {
		return App{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return App{}, errors.Wrapf(err, "Error creating app %s, response code: %d", name, resp.StatusCode)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return App{}, errors.Wrap(err, "Error creating app")
	}
	err = json.Unmarshal(resBody, &appResp)
	if err != nil {
		return App{}, errors.Wrap(err, "Error unmarshalling app")
	}
	return c.mergeAppResource(appResp), nil
}

func (c *Client) DeleteApp(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/apps/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting app %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) mergeAppResource(app AppResource) App {
	app.Entity.Guid = app.Meta.Guid
	app.Entity.CreatedAt = app.Meta.CreatedAt
	app.Entity.UpdatedAt = app.Meta.UpdatedAt
	app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
	app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
	app.Entity.c = c
	return app.Entity
}

func isRedirect(res *http.Response) bool {
	switch res.StatusCode {
	case http.StatusTemporaryRedirect, http.StatusPermanentRedirect, http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther:
		return true
	}
	return false
}
