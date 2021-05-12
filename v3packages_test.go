package cfclient

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListPackagesForAppV3(t *testing.T) {
	Convey("List Package for V3 Apps", t, func() {
		setup(MockRoute{"GET", "/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages", []string{listPackagesForV3AppPayloadPage1, listPackagesForV3AppPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		packages, err := client.ListPackagesForAppV3("f2efe391-2b5b-4836-8518-ad93fa9ebf69", nil)
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

func TestCopyPackageV3(t *testing.T) {
	Convey("Copy V3 Package", t, func() {
		expectedBody := `{"relationships":{"app":{"data":{"guid":"app-guid"}}}}`
		setup(MockRoute{"POST", "/v3/packages", []string{copyPackageV3Payload}, "", http.StatusCreated, "source_guid=package-guid", &expectedBody}, t)
		defer teardown()

		c := &Config{ApiAddress: server.URL, Token: "foobar"}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		pkg, err := client.CopyPackageV3("package-guid", "app-guid")
		So(err, ShouldBeNil)

		So(pkg.State, ShouldEqual, "COPYING")
		So(pkg.Type, ShouldEqual, "docker")
		So(pkg.GUID, ShouldEqual, "fec72fc1-e453-4463-a86d-5df426f337a3")

		docker, err := pkg.DockerData()
		So(err, ShouldBeNil)
		So(docker.Image, ShouldEqual, "http://awesome-sauce.example.org")
	})
}

func TestV3PackageDataDocker(t *testing.T) {
	Convey("V3 Package Data [type=docker]", t, func() {
		Convey("Errors when type=bits", func() {
			p := V3Package{Type: "bits"}
			_, err := p.DockerData()
			So(err, ShouldNotBeNil)
		})

		Convey("Unmarshals docker package", func() {
			p := V3Package{
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

func TestV3PackageDataBits(t *testing.T) {
	Convey("V3 Package Data [type=bits]", t, func() {
		Convey("Errors when type=docker", func() {
			p := V3Package{Type: "docker"}
			_, err := p.BitsData()
			So(err, ShouldNotBeNil)
		})

		Convey("Unmarshals docker package", func() {
			p := V3Package{
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
