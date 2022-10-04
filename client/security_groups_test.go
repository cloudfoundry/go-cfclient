package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	. "github.com/smartystreets/goconvey/convey"
)

func TestListSecurityGroupsByQuery(t *testing.T) {
	Convey("List All  Security Groups", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listSecurityGroupsPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.ListSecurityGroupsByQuery(nil)
		So(err, ShouldBeNil)
		So(securityGroups, ShouldHaveLength, 2)

		So(securityGroups[0].Name, ShouldEqual, "my-group1")
		So(securityGroups[0].GUID, ShouldEqual, "guid-1")
		So(securityGroups[1].Name, ShouldEqual, "my-group2")
		So(securityGroups[1].GUID, ShouldEqual, "guid-2")

		So(securityGroups[0].GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups[0].Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups[0].Rules[0].Destination, ShouldEqual, "1.2.3.4/10")
		So(securityGroups[0].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(*securityGroups[0].Rules[1].Type, ShouldEqual, 8)
		So(*securityGroups[0].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[0].Rules[1].Description, ShouldEqual, "test-desc-1")
		So(securityGroups[0].Relationships["staging_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
		So(securityGroups[1].GloballyEnabled.Staging, ShouldEqual, true)
		So(securityGroups[1].Rules[1].Protocol, ShouldEqual, "icmp")
		So(securityGroups[1].Rules[1].Destination, ShouldEqual, "1.2.3.4/16")
		So(securityGroups[1].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(*securityGroups[1].Rules[1].Type, ShouldEqual, 5)
		So(*securityGroups[1].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[1].Rules[1].Description, ShouldEqual, "test-desc-2")
		So(securityGroups[1].Relationships["running_spaces"].Data[0].GUID, ShouldEqual, "space-guid-5")
		So(securityGroups[1].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-2")
	})

	Convey("List  Security Groups By GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listSecurityGroupsByGuidPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["guids"] = []string{"guid-1"}

		securityGroups, err := client.ListSecurityGroupsByQuery(query)
		So(err, ShouldBeNil)
		So(securityGroups, ShouldHaveLength, 1)

		So(securityGroups[0].Name, ShouldEqual, "my-group1")
		So(securityGroups[0].GUID, ShouldEqual, "guid-1")

		So(securityGroups[0].GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups[0].Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups[0].Rules[0].Destination, ShouldEqual, "1.2.3.4/10")
		So(securityGroups[0].Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(*securityGroups[0].Rules[1].Type, ShouldEqual, 8)
		So(*securityGroups[0].Rules[1].Code, ShouldEqual, 0)
		So(securityGroups[0].Rules[1].Description, ShouldEqual, "test-desc-1")
		So(securityGroups[0].Relationships["staging_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups[0].Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
	})
}

func TestCreateSecurityGroup(t *testing.T) {
	Convey("Create  Security Group With Minimal Parameters", t, func() {
		expectedRequestBody := `{"name":"my-sec-group"}`
		expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
		setup(MockRoute{"POST", "/v3/security_groups", []string{expectedResponseBody}, "", http.StatusCreated, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.CreateSecurityGroup(resource.CreateSecurityGroupRequest{
			Name: "my-sec-group",
		})
		So(err, ShouldBeNil)
		So(securityGroups, ShouldNotBeNil)
		So(securityGroups.GUID, ShouldEqual, "guid-1")
		So(securityGroups.Name, ShouldEqual, "my-sec-group")
		So(securityGroups.GloballyEnabled.Running, ShouldEqual, false)
		So(securityGroups.GloballyEnabled.Staging, ShouldEqual, false)
		So(len(securityGroups.Rules), ShouldEqual, 0)
	})

	Convey("Create  Security Group With Optional Parameters", t, func() {
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
		So(err, ShouldBeNil)
		expectedRequestBody := buffer.String()
		setup(MockRoute{"POST", "/v3/security_groups", []string{genericSecurityGroupPayload}, "", http.StatusCreated, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
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
		So(err, ShouldBeNil)
		So(securityGroups, ShouldNotBeNil)
		So(securityGroups.GUID, ShouldEqual, "guid-1")
		So(securityGroups.Name, ShouldEqual, "my-sec-group")
		So(securityGroups.GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups.GloballyEnabled.Staging, ShouldEqual, false)
		So(securityGroups.Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups.Rules[0].Destination, ShouldEqual, "10.10.10.0/24")
		So(securityGroups.Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroups.Rules[1].Protocol, ShouldEqual, "icmp")
		So(securityGroups.Rules[1].Destination, ShouldEqual, "10.10.11.0/24")
		So(*securityGroups.Rules[1].Type, ShouldEqual, 8)
		So(*securityGroups.Rules[1].Code, ShouldEqual, 0)
		So(securityGroups.Rules[1].Description, ShouldEqual, "Allow ping requests to private services")
		So(len(securityGroups.Relationships["staging_spaces"].Data), ShouldEqual, 0)
		So(securityGroups.Relationships["running_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups.Relationships["running_spaces"].Data[1].GUID, ShouldEqual, "space-guid-2")
		So(securityGroups.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
	})
	Convey("Create  Security Group with empty icmp type and code and no name", t, func() {
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

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		_, err = client.CreateSecurityGroup(resource.CreateSecurityGroupRequest{
			Rules: []*resource.Rule{
				{
					Protocol:    "icmp",
					Destination: "10.10.11.0/24",
				},
			},
		})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "code is required for protocols of type ICMP")
		So(err.Error(), ShouldContainSubstring, "type is required for protocols of type ICMP")
		So(err.Error(), ShouldContainSubstring, "Name can't be blank, Name must be a string")

	})
}

func TestDeleteSecurityGroup(t *testing.T) {
	Convey("Delete  Security Group", t, func() {
		setup(MockRoute{"DELETE", "/v3/security_groups/security-group-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteSecurityGroup("security-group-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateSecurityGroup(t *testing.T) {
	Convey("Update  Security Group with empty type and code", t, func() {
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

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		_, err = client.UpdateSecurityGroup("guid-1", resource.UpdateSecurityGroupRequest{
			Rules: []*resource.Rule{
				{
					Protocol:    "icmp",
					Destination: "10.10.11.0/24",
				},
			},
		})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "code is required for protocols of type ICMP")
		So(err.Error(), ShouldContainSubstring, "type is required for protocols of type ICMP")
	})

	Convey("Update name of  Security Group", t, func() {
		expectedRequestBody := `{"name":"my-sec-group"}`
		expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
		setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{expectedResponseBody}, "", http.StatusOK, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.UpdateSecurityGroup("guid-1", resource.UpdateSecurityGroupRequest{
			Name: "my-sec-group",
		})
		So(err, ShouldBeNil)
		So(securityGroups, ShouldNotBeNil)
		So(securityGroups.GUID, ShouldEqual, "guid-1")
		So(securityGroups.Name, ShouldEqual, "my-sec-group")
		So(securityGroups.GloballyEnabled.Running, ShouldEqual, false)
		So(securityGroups.GloballyEnabled.Staging, ShouldEqual, false)
		So(len(securityGroups.Rules), ShouldEqual, 0)
	})

	Convey("Update  Security Group With Optional Parameters", t, func() {
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
		So(err, ShouldBeNil)
		expectedRequestBody := buffer.String()
		setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{genericSecurityGroupPayload}, "", http.StatusOK, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
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
		So(err, ShouldBeNil)
		So(securityGroups, ShouldNotBeNil)
		So(securityGroups.GUID, ShouldEqual, "guid-1")
		So(securityGroups.Name, ShouldEqual, "my-sec-group")
		So(securityGroups.GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroups.GloballyEnabled.Staging, ShouldEqual, false)
		So(securityGroups.Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroups.Rules[0].Destination, ShouldEqual, "10.10.10.0/24")
		So(securityGroups.Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroups.Rules[1].Protocol, ShouldEqual, "icmp")
		So(securityGroups.Rules[1].Destination, ShouldEqual, "10.10.11.0/24")
		So(*securityGroups.Rules[1].Type, ShouldEqual, 8)
		So(*securityGroups.Rules[1].Code, ShouldEqual, 0)
		So(securityGroups.Rules[1].Description, ShouldEqual, "Allow ping requests to private services")
		So(len(securityGroups.Relationships["staging_spaces"].Data), ShouldEqual, 0)
		So(securityGroups.Relationships["running_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroups.Relationships["running_spaces"].Data[1].GUID, ShouldEqual, "space-guid-2")
		So(securityGroups.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
	})
}

func TestGetSecurityGroup(t *testing.T) {
	Convey("Get  Security Group", t, func() {
		setup(MockRoute{"GET", "/v3/security_groups/guid-1", []string{genericSecurityGroupPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroup, err := client.GetSecurityGroupByGUID("guid-1")
		So(err, ShouldBeNil)
		So(securityGroup, ShouldNotBeNil)
		So(securityGroup.GUID, ShouldEqual, "guid-1")
		So(securityGroup.Name, ShouldEqual, "my-sec-group")
		So(securityGroup.GloballyEnabled.Running, ShouldEqual, true)
		So(securityGroup.GloballyEnabled.Staging, ShouldEqual, false)
		So(securityGroup.Rules[0].Protocol, ShouldEqual, "tcp")
		So(securityGroup.Rules[0].Destination, ShouldEqual, "10.10.10.0/24")
		So(securityGroup.Rules[0].Ports, ShouldEqual, "443,80,8080")
		So(securityGroup.Rules[1].Protocol, ShouldEqual, "icmp")
		So(securityGroup.Rules[1].Destination, ShouldEqual, "10.10.11.0/24")
		So(*securityGroup.Rules[1].Type, ShouldEqual, 8)
		So(*securityGroup.Rules[1].Code, ShouldEqual, 0)
		So(securityGroup.Rules[1].Description, ShouldEqual, "Allow ping requests to private services")
		So(len(securityGroup.Relationships["staging_spaces"].Data), ShouldEqual, 0)
		So(securityGroup.Relationships["running_spaces"].Data[0].GUID, ShouldEqual, "space-guid-1")
		So(securityGroup.Relationships["running_spaces"].Data[1].GUID, ShouldEqual, "space-guid-2")
		So(securityGroup.Links["self"].Href, ShouldEqual, "https://api.example.org/v3/security_groups/guid-1")
	})
}
