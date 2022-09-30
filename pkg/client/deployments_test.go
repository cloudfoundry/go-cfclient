package client

import (
	"net/http"
	"testing"

	v3 "github.com/cloudfoundry-community/go-cfclient/pkg/v3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDeployment(t *testing.T) {
	Convey("Get  Deployment", t, func() {
		setup(MockRoute{"GET", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4", []string{getDeploymentPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.GetDeployment("59c3d133-2b83-46f3-960e-7765a129aea4")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})
}

func TestCreateDeployment(t *testing.T) {
	Convey("Create  Deployment without optional parameters", t, func() {
		body := `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateDeployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", nil)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})

	Convey("Create  Deployment with droplet", t, func() {
		body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateDeployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &v3.CreateDeploymentOptionalParameters{
			Droplet: &v3.Relationship{
				GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
			},
		})
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})

	Convey("Create  Deployment with revision", t, func() {
		body := `{"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateDeployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &v3.CreateDeploymentOptionalParameters{
			Revision: &v3.DeploymentRevision{
				GUID:    "56126cba-656a-4eba-a81e-7e9951b2df57",
				Version: 1,
			},
		})
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})

	Convey("Create  Deployment with revision and droplet", t, func() {
		body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getDeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateDeployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &v3.CreateDeploymentOptionalParameters{
			Droplet: &v3.Relationship{
				GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
			},
			Revision: &v3.DeploymentRevision{
				GUID:    "56126cba-656a-4eba-a81e-7e9951b2df57",
				Version: 1,
			},
		})
		So(err, ShouldNotBeNil)
		So(resp, ShouldBeNil)
	})
}

func TestCancelDeployment(t *testing.T) {
	Convey("Cancel  deployment", t, func() {
		setup(MockRoute{"POST", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4/actions/cancel", []string{""}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.CancelDeployment("59c3d133-2b83-46f3-960e-7765a129aea4")
		So(err, ShouldBeNil)
	})
}
