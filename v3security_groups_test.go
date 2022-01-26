package cfclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListV3SecurityGroupsByQuery(t *testing.T) {
	Convey("List All V3 Security Groups", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listV3SecurityGroupsPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.ListV3SecurityGroupsByQuery(nil)
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

	Convey("List V3 Security Groups By GUID", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/security_groups", []string{listV3SecurityGroupsByGuidPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		query := url.Values{}
		query["guids"] = []string{"guid-1"}

		securityGroups, err := client.ListV3SecurityGroupsByQuery(query)
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

func TestCreateV3SecurityGroup(t *testing.T) {
	Convey("Create V3 Security Group With Minimal Parameters", t, func() {
		expectedRequestBody := `{"name":"my-sec-group"}`
		expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
		setup(MockRoute{"POST", "/v3/security_groups", []string{expectedResponseBody}, "", http.StatusCreated, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.CreateV3SecurityGroup(CreateV3SecurityGroupRequest{
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

	Convey("Create V3 Security Group With Optional Parameters", t, func() {
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
		buffer.Bytes()
		expectedRequestBody := string(buffer.Bytes())
		setup(MockRoute{"POST", "/v3/security_groups", []string{createOrUpdateV3SecurityGroupPayload}, "", http.StatusCreated, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		icmpType := 8
		icmpCode := 0
		createV3SecGroupRequest := CreateV3SecurityGroupRequest{
			Name: "my-sec-group",
			GloballyEnabled: &V3GloballyEnabled{
				Running: true,
				Staging: false,
			},
			Rules: []*V3Rule{
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
			Relationships: map[string]V3ToManyRelationships{
				"running_spaces": {
					Data: []V3Relationship{
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

		securityGroups, err := client.CreateV3SecurityGroup(createV3SecGroupRequest)
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
	Convey("Create V3 Security Group with empty icmp type and code and no name", t, func() {
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

		_, err = client.CreateV3SecurityGroup(CreateV3SecurityGroupRequest{
			Rules: []*V3Rule{
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

func TestDeleteV3SecurityGroup(t *testing.T) {
	Convey("Delete V3 Security Group", t, func() {
		setup(MockRoute{"DELETE", "/v3/security_groups/security-group-guid", []string{""}, "", http.StatusAccepted, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteV3SecurityGroup("security-group-guid")
		So(err, ShouldBeNil)
	})
}

func TestUpdateV3SecurityGroup(t *testing.T) {
	Convey("Update V3 Security Group with empty type and code", t, func() {
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

		_, err = client.UpdateV3SecurityGroup("guid-1", UpdateV3SecurityGroupRequest{
			Rules: []*V3Rule{
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

	Convey("Update name of V3 Security Group", t, func() {
		expectedRequestBody := `{"name":"my-sec-group"}`
		expectedResponseBody := `{"guid":"guid-1", "name":"my-sec-group", "globally_enabled": {"running": false,"staging": false}, "rules": []}`
		setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{expectedResponseBody}, "", http.StatusOK, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		securityGroups, err := client.UpdateV3SecurityGroup("guid-1", UpdateV3SecurityGroupRequest{
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

	Convey("Update V3 Security Group With Optional Parameters", t, func() {
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
		buffer.Bytes()
		expectedRequestBody := string(buffer.Bytes())
		setup(MockRoute{"PATCH", "/v3/security_groups/guid-1", []string{createOrUpdateV3SecurityGroupPayload}, "", http.StatusOK, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		icmpType := 8
		icmpCode := 0
		updateV3SecGroupRequest := UpdateV3SecurityGroupRequest{
			Name: "my-sec-group",
			GloballyEnabled: &V3GloballyEnabled{
				Running: true,
				Staging: false,
			},
			Rules: []*V3Rule{
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

		securityGroups, err := client.UpdateV3SecurityGroup("guid-1", updateV3SecGroupRequest)
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
