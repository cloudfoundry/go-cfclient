package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudfoundry-community/go-cfclient/resource"
	. "github.com/smartystreets/goconvey/convey"
)

func TestListPackagesForApp(t *testing.T) {
	Convey("List Package for  Apps", t, func() {
		setup(MockRoute{"GET", "/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages", []string{listPackagesForAppPayloadPage1, listPackagesForAppPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c, _ := NewTokenConfig(server.URL, "foobar")
		client, err := New(c)
		So(err, ShouldBeNil)

		packages, err := client.ListPackagesForApp("f2efe391-2b5b-4836-8518-ad93fa9ebf69", nil)
		So(err, ShouldBeNil)
		So(packages, ShouldHaveLength, 2)

		So(packages[0].Type, ShouldEqual, "bits")
		So(packages[0].State, ShouldEqual, "READY")
		So(packages[0].Links["app"].Href, ShouldEqual, "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69")
		So(packages[0].Links["download"].Href, ShouldEqual, "https://api.example.org/v3/packages/752edab0-2147-4f58-9c25-cd72ad8c3561/download")
		So(packages[1].Type, ShouldEqual, "bits")
		So(packages[1].State, ShouldEqual, "READY")
		So(packages[1].Links["app"].Href, ShouldEqual, "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69")
		So(packages[1].Links["download"].Href, ShouldEqual, "https://api.example.org/v3/packages/2345ab-2147-4f58-9c25-cd72ad8c3561/download")

	})
}

func TestCopyPackage(t *testing.T) {
	Convey("Copy  Package", t, func() {
		expectedBody := `{"relationships":{"app":{"data":{"guid":"app-guid"}}}}`
		setup(MockRoute{"POST", "/v3/packages", []string{copyPackagePayload}, "", http.StatusCreated, "source_guid=package-guid", &expectedBody}, t)
		defer teardown()

		c, _ := NewTokenConfig(server.URL, "foobar")
		client, err := New(c)
		So(err, ShouldBeNil)

		pkg, err := client.CopyPackage("package-guid", "app-guid")
		So(err, ShouldBeNil)

		So(pkg.State, ShouldEqual, "COPYING")
		So(pkg.Type, ShouldEqual, "docker")
		So(pkg.GUID, ShouldEqual, "fec72fc1-e453-4463-a86d-5df426f337a3")

		docker, err := pkg.DockerData()
		So(err, ShouldBeNil)
		So(docker.Image, ShouldEqual, "http://awesome-sauce.example.org")
	})
}

func TestPackageDataDocker(t *testing.T) {
	Convey(" Package Data [type=docker]", t, func() {
		Convey("Errors when type=bits", func() {
			p := resource.Package{Type: "bits"}
			_, err := p.DockerData()
			So(err, ShouldNotBeNil)
		})

		Convey("Unmarshals docker package", func() {
			p := resource.Package{
				Type: "docker",
				Data: json.RawMessage(`{"image":"nginx","username":"admin","password":"password"}`),
			}
			d, err := p.DockerData()
			So(err, ShouldBeNil)
			So(d.Image, ShouldEqual, "nginx")
			So(d.Username, ShouldEqual, "admin")
			So(d.Password, ShouldEqual, "password")
		})
	})
}

func TestPackageDataBits(t *testing.T) {
	Convey(" Package Data [type=bits]", t, func() {
		Convey("Errors when type=docker", func() {
			p := resource.Package{Type: "docker"}
			_, err := p.BitsData()
			So(err, ShouldNotBeNil)
		})

		Convey("Unmarshals docker package", func() {
			p := resource.Package{
				Type: "bits",
				Data: json.RawMessage(`{"error":"None","checksum":{"type":"sha256","value":"foo"}}`),
			}
			b, err := p.BitsData()
			So(err, ShouldBeNil)
			So(b.Error, ShouldEqual, "None")
			So(b.Checksum.Type, ShouldEqual, "sha256")
		})
	})
}
