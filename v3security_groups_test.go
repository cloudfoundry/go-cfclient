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

		securityGroups, err := client.ListV3SecGroupsByQuery(nil)
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

		securityGroups, err := client.ListV3SecGroupsByQuery(query)
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

		securityGroups, err := client.CreateV3SecGroup(CreateV3SecGroupRequest{
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
		setup(MockRoute{"POST", "/v3/security_groups", []string{createV3SecurityGroupPayload}, "", http.StatusCreated, "", &expectedRequestBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)
		icmpType := 8
		icmpCode := 0
		createV3SecGroupRequest := CreateV3SecGroupRequest{
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

		securityGroups, err := client.CreateV3SecGroup(createV3SecGroupRequest)
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
