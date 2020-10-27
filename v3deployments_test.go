package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDeployment(t *testing.T) {
	Convey("Get V3 Deployment", t, func() {
		setup(MockRoute{"GET", "/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4", getV3DeploymentPayload, "", http.StatusOK, "", nil}, t)
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
