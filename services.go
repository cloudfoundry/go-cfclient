package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
)

type ServicesResponse struct {
	Count     int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	Resources []ServicesResource `json:"resources"`
}

type ServicesResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Service `json:"entity"`
}

type Service struct {
	Guid  string `json:"guid"`
	Label string `json:"label"`
	c     *Client
}

type ServiceSummary struct {
	Guid          string `json:"guid"`
	Name          string `json:"name"`
	BoundAppCount int    `json:"bound_app_count"`
}

func (c *Client) ListServicesByQuery(query url.Values) ([]Service, error) {
	var services []Service
	if query == nil {
		query = url.Values{}
	}
	done := false
	for currPage := 1; !done; currPage++ {
		var serviceResp ServicesResponse
		query.Set("page", fmt.Sprintf("%d", currPage))
		r := c.NewRequest("GET", "/v2/services?"+query.Encode())
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting services")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading services request:")
		}

		err = json.Unmarshal(resBody, &serviceResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling services")
		}
		done = (serviceResp.Pages == currPage)
		for _, service := range serviceResp.Resources {
			service.Entity.Guid = service.Meta.Guid
			service.Entity.c = c
			services = append(services, service.Entity)
		}
	}
	return services, nil
}

func (c *Client) ListServices() ([]Service, error) {
	return c.ListServicesByQuery(nil)
}
