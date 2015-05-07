package cfclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	Guid        string                 `json:"guid"`
	Name        string                 `json:"name"`
	Environment map[string]interface{} `json:"environment_json"`
	SpaceURL    string                 `json:"space_url"`
	SpaceData   SpaceResource          `json:"space"`
	c           *Client
}

func (a *App) Space() Space {
	var spaceResource SpaceResource
	r := a.c.newRequest("GET", a.SpaceURL)
	resp, err := a.c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting space %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading space request %v", resBody)
	}

	err = json.Unmarshal(resBody, &spaceResource)
	if err != nil {
		log.Printf("Error unmarshaling space %v", err)
	}
	spaceResource.Entity.Guid = spaceResource.Meta.Guid
	spaceResource.Entity.c = a.c
	return spaceResource.Entity
}

func (c *Client) ListApps() []App {
	var apps []App

	requestUrl := "/v2/apps?inline-relations-depth=2"
	for {
		var appResp AppResponse
		r := c.newRequest("GET", requestUrl)
		resp, err := c.doRequest(r)

		if err != nil {
			log.Printf("Error requesting apps %v", err)
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading app request %v", resBody)
		}

		err = json.Unmarshal(resBody, &appResp)
		if err != nil {
			log.Printf("Error unmarshaling app %v", err)
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
	return apps
}

func (c *Client) AppByGuid(guid string) App {
	var appResource AppResource
	r := c.newRequest("GET", "/v2/apps/"+guid+"?inline-relations-depth=2")
	resp, err := c.doRequest(r)
	if err != nil {
		log.Printf("Error requesting apps %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading app request %v", resBody)
	}

	err = json.Unmarshal(resBody, &appResource)
	if err != nil {
		log.Printf("Error unmarshaling app %v", err)
	}
	appResource.Entity.Guid = appResource.Meta.Guid
	appResource.Entity.SpaceData.Entity.Guid = appResource.Entity.SpaceData.Meta.Guid
	appResource.Entity.SpaceData.Entity.OrgData.Entity.Guid = appResource.Entity.SpaceData.Entity.OrgData.Meta.Guid
	appResource.Entity.c = c
	return appResource.Entity
}
