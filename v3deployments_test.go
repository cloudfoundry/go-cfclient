package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDeployment(t *testing.T) {
	Convey("Get V3 Deployment", t, func() {
		setup(MockRoute{"GET", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4", []string{getV3DeploymentPayload}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.GetV3Deployment("59c3d133-2b83-46f3-960e-7765a129aea4")
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})
}

func TestCreateDeployment(t *testing.T) {
	Convey("Create V3 Deployment without optional parameters", t, func() {
		body := `{"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getV3DeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateV3Deployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", nil)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})

	Convey("Create V3 Deployment with droplet", t, func() {
		body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getV3DeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateV3Deployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &CreateV3DeploymentOptionalParameters{
			Droplet: &V3Relationship{
				GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
			},
		})
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)

		So(resp.GUID, ShouldEqual, "59c3d133-2b83-46f3-960e-7765a129aea4")
		So(resp.Status.Reason, ShouldEqual, "DEPLOYING")
		So(resp.Relationships["app"].Data.GUID, ShouldEqual, "305cea31-5a44-45ca-b51b-e89c7a8ef8b2")
	})

	Convey("Create V3 Deployment with revision", t, func() {
		body := `{"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getV3DeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateV3Deployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &CreateV3DeploymentOptionalParameters{
			Revision: &V3DeploymentRevision{
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

	Convey("Create V3 Deployment with revision and droplet", t, func() {
		body := `{"droplet":{"guid":"44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"},"revision":{"guid":"56126cba-656a-4eba-a81e-7e9951b2df57","version":1},"relationships":{"app":{"data":{"guid":"305cea31-5a44-45ca-b51b-e89c7a8ef8b2"}}}}`
		setup(MockRoute{"POST", "/v3/deployments", []string{getV3DeploymentPayload}, "", http.StatusCreated, "", &body}, t)

		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		resp, err := client.CreateV3Deployment("305cea31-5a44-45ca-b51b-e89c7a8ef8b2", &CreateV3DeploymentOptionalParameters{
			Droplet: &V3Relationship{
				GUID: "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962",
			},
			Revision: &V3DeploymentRevision{
				GUID:    "56126cba-656a-4eba-a81e-7e9951b2df57",
				Version: 1,
			},
		})
		So(err, ShouldNotBeNil)
		So(resp, ShouldBeNil)
	})
}

func TestCancelV3Deployment(t *testing.T) {
	Convey("Cancel V3 deployment", t, func() {
		setup(MockRoute{"POST", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4/actions/cancel", []string{""}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.CancelV3Deployment("59c3d133-2b83-46f3-960e-7765a129aea4")
		So(err, ShouldBeNil)
	})
}
