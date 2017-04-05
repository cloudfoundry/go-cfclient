package cfclient

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
)

type ServiceBindingsResponse struct {
	Count     int                      `json:"total_results"`
	Pages     int                      `json:"total_pages"`
	Resources []ServiceBindingResource `json:"resources"`
}

type ServiceBindingResource struct {
	Meta   Meta           `json:"metadata"`
	Entity ServiceBinding `json:"entity"`
}

type ServiceBinding struct {
	Guid                string      `json:"guid"`
	AppGuid             string      `json:"app_guid"`
	ServiceInstanceGuid string      `json:"service_instance_guid"`
	Credentials         interface{} `json:"credentials"`
	BindingOptions      interface{} `json:"binding_options"`
	GatewayData         interface{} `json:"gateway_data"`
	GatewayName         string      `json:"gateway_name"`
	SyslogDrainUrl      string      `json:"syslog_drain_url"`
	VolumeMounts        interface{} `json:"volume_mounts"`
	AppUrl              string      `json:"app_url"`
	ServiceInstanceUrl  string      `json:"service_instance_url"`
	c                   *Client
}

func (c *Client) ListServiceBindingsByQuery(query url.Values) ([]ServiceBinding, error) {
	var serviceBindings []ServiceBinding
	var serviceBindingsResp ServiceBindingsResponse
	r := c.NewRequest("GET", "/v2/service_bindings?"+query.Encode())
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting service bindings")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading service bindings request:")
	}

	err = json.Unmarshal(resBody, &serviceBindingsResp)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling service bindings")
	}
	for _, serviceBinding := range serviceBindingsResp.Resources {
		serviceBinding.Entity.Guid = serviceBinding.Meta.Guid
		serviceBinding.Entity.c = c
		serviceBindings = append(serviceBindings, serviceBinding.Entity)
	}
	return serviceBindings, nil
}

func (c *Client) ListServiceBindings() ([]ServiceBinding, error) {
	return c.ListServiceBindingsByQuery(nil)
}
