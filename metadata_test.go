package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateOrgMetadata(t *testing.T) {
	Convey("Set metadata on org", t, func() {
		updateOrgMetadatPayload := `{"metadata":{"annotations":{"hello":"world"},"labels":{"foo":"bar"}}}`
		mocks := []MockRoute{
			{"PATCH", "/v3/organizations/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateOrgMetadatPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata := Metadata{}
		metadata.AddAnnotation("hello", "world")
		metadata.AddLabel("", "foo", "bar")

		err = client.UpdateOrgMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2", metadata)
		So(err, ShouldBeNil)
	})

	Convey("Remove metadata on org", t, func() {
		updateOrgMetadatPayload := `{"metadata":{"annotations":{"hello":null},"labels":{"foo":null}}}`
		mocks := []MockRoute{
			{"PATCH", "/v3/organizations/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateOrgMetadatPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata := Metadata{}
		metadata.RemoveAnnotation("hello")
		metadata.RemoveLabel("", "foo")

		err = client.UpdateOrgMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2", metadata)
		So(err, ShouldBeNil)
	})
}
func TestUpdateSpaceMetadata(t *testing.T) {
	Convey("Set metadata on org", t, func() {
		updateSpaceMetadataPayload := `{"metadata":{"annotations":{"hello":"world"},"labels":{"foo":"bar"}}}`
		mocks := []MockRoute{
			{"PATCH", "/v3/spaces/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateSpaceMetadataPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata := Metadata{}
		metadata.AddAnnotation("hello", "world")
		metadata.AddLabel("", "foo", "bar")

		err = client.UpdateSpaceMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2", metadata)
		So(err, ShouldBeNil)
	})

	Convey("Remove metadata on space", t, func() {
		updateSpaceMetadataPayload := `{"metadata":{"annotations":{"hello":null},"labels":{"foo":null}}}`
		mocks := []MockRoute{
			{"PATCH", "/v3/spaces/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateSpaceMetadataPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata := Metadata{}
		metadata.RemoveAnnotation("hello")
		metadata.RemoveLabel("", "foo")

		err = client.UpdateSpaceMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2", metadata)
		So(err, ShouldBeNil)
	})
}

func TestSpaceMetadata(t *testing.T) {
	Convey("Getting space metadata", t, func() {
		spaceMetadataPayload := `{"guid": "3b6f763f-aae1-4177-9b93-f2de6f2a48f2","name": "space2","metadata":{"annotations":{"hello":"world"},"labels":{"foo":"bar"}}}`
		mocks := []MockRoute{
			{"GET", "/v3/spaces/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{spaceMetadataPayload}, "", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata, err := client.SpaceMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2")
		So(err, ShouldBeNil)
		So(metadata, ShouldNotBeNil)
		So(metadata.Labels, ShouldContainKey, "foo")
		So(metadata.Annotations, ShouldContainKey, "hello")
	})

}

func TestOrgMetadata(t *testing.T) {
	Convey("Getting org metadata", t, func() {
		orgMetadataPayload := `{"guid": "3b6f763f-aae1-4177-9b93-f2de6f2a48f2","name": "space2","metadata":{"annotations":{"hello":"world"},"labels":{"foo":"bar"}}}`
		mocks := []MockRoute{
			{"GET", "/v3/organizations/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{orgMetadataPayload}, "", 200, "", nil},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		metadata, err := client.OrgMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2")
		So(err, ShouldBeNil)
		So(metadata, ShouldNotBeNil)
		So(metadata.Labels, ShouldContainKey, "foo")
		So(metadata.Annotations, ShouldContainKey, "hello")
	})
}

func TestRemoveOrgMetadata(t *testing.T) {
	Convey("Removing org metadata", t, func() {
		orgMetadataPayload := `{"guid": "3b6f763f-aae1-4177-9b93-f2de6f2a48f2","name": "space2","metadata":{"annotations":{"hello":"world","prefix/foo":"bar"},"labels":{"foo":"bar"}}}`
		updateMetadataPayload := `{"metadata":{"annotations":{"foo":null,"hello":null,"prefix/foo":null},"labels":{"foo":null}}}`
		mocks := []MockRoute{
			{"GET", "/v3/organizations/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{orgMetadataPayload}, "", 200, "", nil},
			{"PATCH", "/v3/organizations/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateMetadataPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.RemoveOrgMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2")
		So(err, ShouldBeNil)
	})
}

func TestRemoveSpaceMetadata(t *testing.T) {
	Convey("Removing org metadata", t, func() {
		spaceMetadataPayload := `{"guid": "3b6f763f-aae1-4177-9b93-f2de6f2a48f2","name": "space2","metadata":{"annotations":{"hello":"world","prefix/foo":"bar"},"labels":{"foo":"bar"}}}`
		updateMetadataPayload := `{"metadata":{"annotations":{"foo":null,"hello":null,"prefix/foo":null},"labels":{"foo":null}}}`
		mocks := []MockRoute{
			{"GET", "/v3/spaces/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{spaceMetadataPayload}, "", 200, "", nil},
			{"PATCH", "/v3/spaces/3b6f763f-aae1-4177-9b93-f2de6f2a48f2", []string{""}, "", 200, "", &updateMetadataPayload},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.RemoveSpaceMetadata("3b6f763f-aae1-4177-9b93-f2de6f2a48f2")
		So(err, ShouldBeNil)
	})
}
