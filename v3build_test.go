package cfclient

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateV3Build(t *testing.T) {
	Convey("Get V3 App Environment", t, func() {
		body := `{"metadata":{"labels":{"foo":"bar"}},"package":{"guid":"package-guid"}}`
		setup(MockRoute{"POST", "/v3/builds", []string{createV3BuildPayload}, "", http.StatusCreated, "", &body}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		build, err := client.CreateV3Build("package-guid", nil,
			&V3Metadata{Labels: map[string]string{"foo": "bar"}})
		So(err, ShouldBeNil)
		So(build, ShouldNotBeNil)

		So(build.GUID, ShouldEqual, "585bc3c1-3743-497d-88b0-403ad6b56d16")
		So(build.CreatedBy.Name, ShouldEqual, "bill")
		So(build.Package.GUID, ShouldEqual, "8e4da443-f255-499c-8b47-b3729b5b7432")
	})
}
