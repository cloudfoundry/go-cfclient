package client

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
)

func TestListSecurityGroupsByQuery(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/security_groups", []string{listSecurityGroupsPayload}, "", http.StatusOK, "", nil},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	securityGroups, err := client.ListSecurityGroupsByQuery(nil)
	require.NoError(t, err)
	require.Len(t, securityGroups, 2)

	require.Equal(t, "my-group1", securityGroups[0].Name)
	require.Equal(t, "guid-1", securityGroups[0].GUID)
	require.Equal(t, "my-group2", securityGroups[1].Name)
	require.Equal(t, "guid-2", securityGroups[1].GUID)

	require.Equal(t, true, securityGroups[0].GloballyEnabled.Running)
	require.Equal(t, "tcp", securityGroups[0].Rules[0].Protocol)
	require.Equal(t, "1.2.3.4/10", securityGroups[0].Rules[0].Destination)
	require.Equal(t, "443,80,8080", securityGroups[0].Rules[0].Ports)
	require.Equal(t, 8, *securityGroups[0].Rules[1].Type)
	require.Equal(t, 0, *securityGroups[0].Rules[1].Code)
	require.Equal(t, "test-desc-1", securityGroups[0].Rules[1].Description)
	require.Equal(t, "space-guid-1", securityGroups[0].Relationships["staging_spaces"].Data[0].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-1", securityGroups[0].Links["self"].Href)
	require.Equal(t, true, securityGroups[1].GloballyEnabled.Staging)
	require.Equal(t, "icmp", securityGroups[1].Rules[1].Protocol)
	require.Equal(t, "1.2.3.4/16", securityGroups[1].Rules[1].Destination)
	require.Equal(t, "443,80,8080", securityGroups[1].Rules[0].Ports)
	require.Equal(t, 5, *securityGroups[1].Rules[1].Type)
	require.Equal(t, 0, *securityGroups[1].Rules[1].Code)
	require.Equal(t, "test-desc-2", securityGroups[1].Rules[1].Description)
	require.Equal(t, "space-guid-5", securityGroups[1].Relationships["running_spaces"].Data[0].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-2", securityGroups[1].Links["self"].Href)
}

func TestListSecurityGroupsByQueryWithGroupGUID(t *testing.T) {
	mocks := []MockRoute{
		{"GET", "/v3/security_groups", []string{listSecurityGroupsByGuidPayload}, "", http.StatusOK, "", nil},
	}
	setupMultiple(mocks, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	query := url.Values{}
	query["guids"] = []string{"guid-1"}

	securityGroups, err := client.ListSecurityGroupsByQuery(query)
	require.NoError(t, err)
	require.Len(t, securityGroups, 1)

	require.Equal(t, "my-group1", securityGroups[0].Name)
	require.Equal(t, "guid-1", securityGroups[0].GUID)

	require.Equal(t, true, securityGroups[0].GloballyEnabled.Running)
	require.Equal(t, "tcp", securityGroups[0].Rules[0].Protocol)
	require.Equal(t, "1.2.3.4/10", securityGroups[0].Rules[0].Destination)
	require.Equal(t, "443,80,8080", securityGroups[0].Rules[0].Ports)
	require.Equal(t, 8, *securityGroups[0].Rules[1].Type)
	require.Equal(t, 0, *securityGroups[0].Rules[1].Code)
	require.Equal(t, "test-desc-1", securityGroups[0].Rules[1].Description)
	require.Equal(t, "space-guid-1", securityGroups[0].Relationships["staging_spaces"].Data[0].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-1", securityGroups[0].Links["self"].Href)
}

func TestCreateSecurityGroupWithMinimalParams(t *testing.T) {
	expectedRequestBody := `{"name":"my-sec-group"}`
	expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
	setup(MockRoute{"POST", "/v3/security_groups", []string{expectedResponseBody}, "", http.StatusCreated, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	securityGroups, err := client.CreateSecurityGroup(resource.CreateSecurityGroupRequest{
		Name: "my-sec-group",
	})
	require.NoError(t, err)
	require.NotNil(t, securityGroups)
	require.Equal(t, "guid-1", securityGroups.GUID)
	require.Equal(t, "my-sec-group", securityGroups.Name)
	require.Equal(t, false, securityGroups.GloballyEnabled.Running)
	require.Equal(t, false, securityGroups.GloballyEnabled.Staging)
	require.Equal(t, 0, len(securityGroups.Rules))
}

func TestCreateSecurityGroupWithOptionalParams(t *testing.T) {
	requestBody := `{
			"name": "my-sec-group",
			"globally_enabled": {
			  "running": true
			},
			"rules": [
			  {
				"protocol": "tcp",
				"destination": "10.10.10.0/24",
				"ports": "443,80,8080"
			  },
			  {
				"protocol": "icmp",
				"destination": "10.10.11.0/24",
				"type": 8,
				"code": 0,
				"description": "Allow ping requests to private services"
			  }
			],
			"relationships": {
			  "running_spaces": {
				"data": [
				  {
					"guid": "space-guid-1"
				  },
				  {
					"guid": "space-guid-2"
				  }
				]
			  }
			}
		  }`
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(requestBody))
	require.NoError(t, err)
	expectedRequestBody := buffer.String()
	setup(MockRoute{"POST", "/v3/security_groups", []string{genericSecurityGroupPayload}, "", http.StatusCreated, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)
	icmpType := 8
	icmpCode := 0
	createSecGroupRequest := resource.CreateSecurityGroupRequest{
		Name: "my-sec-group",
		GloballyEnabled: &resource.GloballyEnabled{
			Running: true,
			Staging: false,
		},
		Rules: []*resource.Rule{
			{
				Protocol:    "tcp",
				Destination: "10.10.10.0/24",
				Ports:       "443,80,8080",
			},
			{
				Protocol:    "icmp",
				Destination: "10.10.11.0/24",
				Type:        &icmpType,
				Code:        &icmpCode,
				Description: "Allow ping requests to private services",
			},
		},
		Relationships: map[string]resource.ToManyRelationships{
			"running_spaces": {
				Data: []resource.Relationship{
					{
						GUID: "space-guid-1",
					},
					{
						GUID: "space-guid-2",
					},
				},
			},
		},
	}

	securityGroups, err := client.CreateSecurityGroup(createSecGroupRequest)
	require.NoError(t, err)
	require.NotNil(t, securityGroups)
	require.Equal(t, "guid-1", securityGroups.GUID)
	require.Equal(t, "my-sec-group", securityGroups.Name)
	require.Equal(t, true, securityGroups.GloballyEnabled.Running)
	require.Equal(t, false, securityGroups.GloballyEnabled.Staging)
	require.Equal(t, "tcp", securityGroups.Rules[0].Protocol)
	require.Equal(t, "10.10.10.0/24", securityGroups.Rules[0].Destination)
	require.Equal(t, "443,80,8080", securityGroups.Rules[0].Ports)
	require.Equal(t, "icmp", securityGroups.Rules[1].Protocol)
	require.Equal(t, "10.10.11.0/24", securityGroups.Rules[1].Destination)
	require.Equal(t, 8, *securityGroups.Rules[1].Type)
	require.Equal(t, 0, *securityGroups.Rules[1].Code)
	require.Equal(t, "Allow ping requests to private services", securityGroups.Rules[1].Description)
	require.Equal(t, 0, len(securityGroups.Relationships["staging_spaces"].Data))
	require.Equal(t, "space-guid-1", securityGroups.Relationships["running_spaces"].Data[0].GUID)
	require.Equal(t, "space-guid-2", securityGroups.Relationships["running_spaces"].Data[1].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-1", securityGroups.Links["self"].Href)
}

func TestCreateSecurityGroupWithICMPTypeCodeAndNoName(t *testing.T) {
	expectedRequestBody := `{"name":"","rules":[{"protocol":"icmp","destination":"10.10.11.0/24"}]}`
	expectedResponseBody := `{
			"errors": [
			  {
				"code": 10008,
				"detail": "Rules[0]: code is required for protocols of type ICMP, Rules[0]: code must be an integer between -1 and 255 (inclusive), Rules[0]: type is required for protocols of type ICMP, Rules[0]: type must be an integer between -1 and 255 (inclusive), Name can't be blank, Name must be a string",
				"title": "CF-UnprocessableEntity"
			  }
			]
		  }`
	setup(MockRoute{"POST", "/v3/security_groups", []string{expectedResponseBody}, "", http.StatusUnprocessableEntity, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	_, err = client.CreateSecurityGroup(resource.CreateSecurityGroupRequest{
		Rules: []*resource.Rule{
			{
				Protocol:    "icmp",
				Destination: "10.10.11.0/24",
			},
		},
	})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "code is required for protocols of type ICMP")
	require.Contains(t, err.Error(), "type is required for protocols of type ICMP")
	require.Contains(t, err.Error(), "Name can't be blank, Name must be a string")
}

func TestDeleteSecurityGroup(t *testing.T) {
	setup(MockRoute{"DELETE", "/v3/security_groups/security-group-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	err = client.DeleteSecurityGroup("security-group-guid")
	require.NoError(t, err)
}

func TestUpdateSecurityGroup(t *testing.T) {
	expectedRequestBody := `{"rules":[{"protocol":"icmp","destination":"10.10.11.0/24"}]}`
	expectedResponseBody := `{
			"errors": [
			  {
				"code": 10008,
				"detail": "Rules[0]: code is required for protocols of type ICMP, Rules[0]: code must be an integer between -1 and 255 (inclusive), Rules[0]: type is required for protocols of type ICMP, Rules[0]: type must be an integer between -1 and 255 (inclusive)",
				"title": "CF-UnprocessableEntity"
			  }
			]
		  }`
	setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{expectedResponseBody}, "", http.StatusUnprocessableEntity, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	_, err = client.UpdateSecurityGroup("guid-1", resource.UpdateSecurityGroupRequest{
		Rules: []*resource.Rule{
			{
				Protocol:    "icmp",
				Destination: "10.10.11.0/24",
			},
		},
	})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "code is required for protocols of type ICMP")
	require.Contains(t, err.Error(), "type is required for protocols of type ICMP")
}

func TestUpdateSecurityGroupUpdateName(t *testing.T) {
	expectedRequestBody := `{"name":"my-sec-group"}`
	expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
	setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{expectedResponseBody}, "", http.StatusOK, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	securityGroups, err := client.UpdateSecurityGroup("guid-1", resource.UpdateSecurityGroupRequest{
		Name: "my-sec-group",
	})
	require.NoError(t, err)
	require.NotNil(t, securityGroups)
	require.Equal(t, "guid-1", securityGroups.GUID)
	require.Equal(t, "my-sec-group", securityGroups.Name)
	require.Equal(t, false, securityGroups.GloballyEnabled.Running)
	require.Equal(t, false, securityGroups.GloballyEnabled.Staging)
	require.Equal(t, 0, len(securityGroups.Rules))
}

func TestUpdateSecurityGroupWithOptionalParams(t *testing.T) {
	requestBody := `{
			"name": "my-sec-group",
			"globally_enabled": {
			  "running": true
			},
			"rules": [
			  {
				"protocol": "tcp",
				"destination": "10.10.10.0/24",
				"ports": "443,80,8080"
			  },
			  {
				"protocol": "icmp",
				"destination": "10.10.11.0/24",
				"type": 8,
				"code": 0,
				"description": "Allow ping requests to private services"
			  }
			]
		  }`
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(requestBody))
	require.NoError(t, err)
	expectedRequestBody := buffer.String()
	setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{genericSecurityGroupPayload}, "", http.StatusOK, "", &expectedRequestBody}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)
	icmpType := 8
	icmpCode := 0
	updateSecGroupRequest := resource.UpdateSecurityGroupRequest{
		Name: "my-sec-group",
		GloballyEnabled: &resource.GloballyEnabled{
			Running: true,
			Staging: false,
		},
		Rules: []*resource.Rule{
			{
				Protocol:    "tcp",
				Destination: "10.10.10.0/24",
				Ports:       "443,80,8080",
			},
			{
				Protocol:    "icmp",
				Destination: "10.10.11.0/24",
				Type:        &icmpType,
				Code:        &icmpCode,
				Description: "Allow ping requests to private services",
			},
		},
	}

	securityGroups, err := client.UpdateSecurityGroup("guid-1", updateSecGroupRequest)
	require.NoError(t, err)
	require.NotNil(t, securityGroups)
	require.Equal(t, "guid-1", securityGroups.GUID)
	require.Equal(t, "my-sec-group", securityGroups.Name)
	require.Equal(t, true, securityGroups.GloballyEnabled.Running)
	require.Equal(t, false, securityGroups.GloballyEnabled.Staging)
	require.Equal(t, "tcp", securityGroups.Rules[0].Protocol)
	require.Equal(t, "10.10.10.0/24", securityGroups.Rules[0].Destination)
	require.Equal(t, "443,80,8080", securityGroups.Rules[0].Ports)
	require.Equal(t, "icmp", securityGroups.Rules[1].Protocol)
	require.Equal(t, "10.10.11.0/24", securityGroups.Rules[1].Destination)
	require.Equal(t, 8, *securityGroups.Rules[1].Type)
	require.Equal(t, 0, *securityGroups.Rules[1].Code)
	require.Equal(t, "Allow ping requests to private services", securityGroups.Rules[1].Description)
	require.Equal(t, 0, len(securityGroups.Relationships["staging_spaces"].Data))
	require.Equal(t, "space-guid-1", securityGroups.Relationships["running_spaces"].Data[0].GUID)
	require.Equal(t, "space-guid-2", securityGroups.Relationships["running_spaces"].Data[1].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-1", securityGroups.Links["self"].Href)
}

func TestGetSecurityGroup(t *testing.T) {
	setup(MockRoute{"GET", "/v3/security_groups/guid-1", []string{genericSecurityGroupPayload}, "", http.StatusOK, "", nil}, t)
	defer teardown()

	c, _ := NewTokenConfig(server.URL, "foobar")
	client, err := New(c)
	require.NoError(t, err)

	securityGroup, err := client.GetSecurityGroupByGUID("guid-1")
	require.NoError(t, err)
	require.NotNil(t, securityGroup)
	require.Equal(t, "guid-1", securityGroup.GUID)
	require.Equal(t, "my-sec-group", securityGroup.Name)
	require.Equal(t, true, securityGroup.GloballyEnabled.Running)
	require.Equal(t, false, securityGroup.GloballyEnabled.Staging)
	require.Equal(t, "tcp", securityGroup.Rules[0].Protocol)
	require.Equal(t, "10.10.10.0/24", securityGroup.Rules[0].Destination)
	require.Equal(t, "443,80,8080", securityGroup.Rules[0].Ports)
	require.Equal(t, "icmp", securityGroup.Rules[1].Protocol)
	require.Equal(t, "10.10.11.0/24", securityGroup.Rules[1].Destination)
	require.Equal(t, 8, *securityGroup.Rules[1].Type)
	require.Equal(t, 0, *securityGroup.Rules[1].Code)
	require.Equal(t, "Allow ping requests to private services", securityGroup.Rules[1].Description)
	require.Equal(t, 0, len(securityGroup.Relationships["staging_spaces"].Data))
	require.Equal(t, "space-guid-1", securityGroup.Relationships["running_spaces"].Data[0].GUID)
	require.Equal(t, "space-guid-2", securityGroup.Relationships["running_spaces"].Data[1].GUID)
	require.Equal(t, "https://api.example.org/v3/security_groups/guid-1", securityGroup.Links["self"].Href)
}
