package cfclient

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateIsolationSegment(t *testing.T) {
	Convey("Create Isolation Segment", t, func() {
		mocks := []MockRoute{
			{"POST", "/v3/isolation_segments", []string{createIsolationSegmentPayload}, "", http.StatusCreated, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		name := "TheKittenIsTheShark"

		isolationsegment, err := client.CreateIsolationSegment(name)
		So(err, ShouldBeNil)

		So(isolationsegment.Name, ShouldEqual, name)
		So(isolationsegment.GUID, ShouldEqual, "323f211e-fea3-4161-9bd1-615392327913")
		So(isolationsegment.CreatedAt.String(), ShouldEqual, time.Date(2016, 10, 19, 20, 25, 04, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment.UpdatedAt.String(), ShouldEqual, time.Date(2016, 11, 8, 16, 41, 26, 0, time.FixedZone("UTC", 0)).String())
	})
}

func TestGetIsolationSegmentByGUID(t *testing.T) {
	Convey("Request existing Isolation Segment", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/isolation_segments/323f211e-fea3-4161-9bd1-615392327913", []string{createIsolationSegmentPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		name := "TheKittenIsTheShark"

		isolationsegment, err := client.GetIsolationSegmentByGUID("323f211e-fea3-4161-9bd1-615392327913")
		So(err, ShouldBeNil)

		So(isolationsegment.Name, ShouldEqual, name)
		So(isolationsegment.GUID, ShouldEqual, "323f211e-fea3-4161-9bd1-615392327913")
		So(isolationsegment.CreatedAt.String(), ShouldEqual, time.Date(2016, 10, 19, 20, 25, 04, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment.UpdatedAt.String(), ShouldEqual, time.Date(2016, 11, 8, 16, 41, 26, 0, time.FixedZone("UTC", 0)).String())
	})

	Convey("Request non-existing Isolation Segment", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/isolation_segments/323f211e-fea3-4161--9bd1-615392327913", []string{createIsolationSegmentPayload}, "", http.StatusOK, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		isolationsegment, err := client.GetIsolationSegmentByGUID("does not exit")
		So(err, ShouldNotBeNil)
		So(isolationsegment, ShouldBeNil)
	})
}

func TestListIsolationSegments(t *testing.T) {
	Convey("Request list of all Isolation Segments", t, func() {
		setup(MockRoute{"GET", "/v3/isolation_segments", []string{listIsolationSegmentsPayloadPage1, listIsolationSegmentsPayloadPage2}, "", http.StatusOK, "", nil}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		isolationsegment, err := client.ListIsolationSegments()
		So(err, ShouldBeNil)
		So(isolationsegment, ShouldNotBeNil)
		So(len(isolationsegment), ShouldEqual, 4)

		So(isolationsegment[0].Name, ShouldEqual, "shared")
		So(isolationsegment[0].GUID, ShouldEqual, "033b4c58-12bb-499a-b05d-4b6fc9e2993b")
		So(isolationsegment[0].CreatedAt.String(), ShouldEqual, time.Date(2017, 4, 2, 11, 22, 4, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment[0].UpdatedAt.String(), ShouldEqual, time.Date(2017, 4, 2, 11, 22, 4, 0, time.FixedZone("UTC", 0)).String())

		So(isolationsegment[1].Name, ShouldEqual, "my_segment")
		So(isolationsegment[1].GUID, ShouldEqual, "23d0baf4-9d3c-44d8-b2dc-1767bcdad1e0")
		So(isolationsegment[1].CreatedAt.String(), ShouldEqual, time.Date(2017, 4, 7, 11, 20, 16, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment[1].UpdatedAt.String(), ShouldEqual, time.Date(2017, 4, 7, 11, 20, 16, 0, time.FixedZone("UTC", 0)).String())

		So(isolationsegment[2].Name, ShouldEqual, "shared1")
		So(isolationsegment[2].GUID, ShouldEqual, "abcdefg12-12bb-499a-b05d-4b6fc9e2993b")
		So(isolationsegment[2].CreatedAt.String(), ShouldEqual, time.Date(2017, 5, 2, 11, 22, 4, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment[2].UpdatedAt.String(), ShouldEqual, time.Date(2017, 5, 2, 11, 22, 4, 0, time.FixedZone("UTC", 0)).String())

		So(isolationsegment[3].Name, ShouldEqual, "my_segment1")
		So(isolationsegment[3].GUID, ShouldEqual, "abcdef123-9d3c-44d8-b2dc-1767bcdad1e0")
		So(isolationsegment[3].CreatedAt.String(), ShouldEqual, time.Date(2017, 5, 7, 11, 20, 16, 0, time.FixedZone("UTC", 0)).String())
		So(isolationsegment[3].UpdatedAt.String(), ShouldEqual, time.Date(2017, 5, 7, 11, 20, 16, 0, time.FixedZone("UTC", 0)).String())
	})
}

func TestDeleteIsolationSegmentByGUID(t *testing.T) {
	Convey("Delete an Isolation Segment by GUID", t, func() {
		mocks := []MockRoute{
			{"DELETE", "/v3/isolation_segments/033b4c58-12bb-499a-b05d-4b6fc9e2993b", []string{""}, "", http.StatusNoContent, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteIsolationSegmentByGUID("033b4c58-12bb-499a-b05d-4b6fc9e2993b")
		So(err, ShouldBeNil)

		err = client.DeleteIsolationSegmentByGUID("theKittenIsTheShark")
		So(err, ShouldNotBeNil)
	})
}

func TestIsolationSegmentMethods(t *testing.T) {

	postData := `{"data":[{"guid":"theKittenIsTheShark"}]}`

	Convey("Request list of all Isolation Segments", t, func() {
		mocks := []MockRoute{
			{"GET", "/v3/isolation_segments", []string{listIsolationSegmentsPayload}, "", http.StatusOK, "", nil},
			{"DELETE", "/v3/isolation_segments/033b4c58-12bb-499a-b05d-4b6fc9e2993b", []string{""}, "", http.StatusNoContent, "", nil},
			{"POST", "/v3/isolation_segments/033b4c58-12bb-499a-b05d-4b6fc9e2993b/relationships/organizations", []string{""}, "", http.StatusOK, "", &postData},
			{"DELETE", "/v3/isolation_segments/033b4c58-12bb-499a-b05d-4b6fc9e2993b/relationships/organizations/theKittenIsTheShark", []string{""}, "", http.StatusNoContent, "", nil},
			{"PUT", "/v2/spaces/theKittenIsTheShark", []string{""}, "", http.StatusCreated, "", nil},
			{"DELETE", "/v2/spaces/theKittenIsTheShark/isolation_segment", []string{""}, "", http.StatusNoContent, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		isolationsegment, err := client.ListIsolationSegments()

		So(err, ShouldBeNil)
		So(isolationsegment, ShouldNotBeNil)
		So(len(isolationsegment), ShouldEqual, 2)

		errAddOrg := isolationsegment[0].AddOrg("theKittenIsTheShark")
		So(errAddOrg, ShouldBeNil)

		errRemOrg := isolationsegment[0].RemoveOrg("theKittenIsTheShark")
		So(errRemOrg, ShouldBeNil)

		errAddSpace := isolationsegment[0].AddSpace("theKittenIsTheShark")
		So(errAddSpace, ShouldBeNil)

		errRemSpace := isolationsegment[0].RemoveSpace("theKittenIsTheShark")
		So(errRemSpace, ShouldBeNil)

		errDel := isolationsegment[0].Delete()
		So(errDel, ShouldBeNil)
	})
}
