package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListDomains(t *testing.T) {
	Convey("List domains", t, func() {
		setup(MockRoute{"GET", "/v2/private_domains", listDomainsPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		domains, err := client.ListDomains()
		So(err, ShouldBeNil)

		So(len(domains), ShouldEqual, 4)
		So(domains[0].Guid, ShouldEqual, "b2a35f0c-d5ad-4a59-bea7-461711d96b0d")
		So(domains[0].Name, ShouldEqual, "vcap.me")
		So(domains[0].OwningOrganizationGuid, ShouldEqual, "4cf3bc47-eccd-4662-9322-7833c3bdcded")
		So(domains[0].OwningOrganizationUrl, ShouldEqual, "/v2/organizations/4cf3bc47-eccd-4662-9322-7833c3bdcded")
		So(domains[0].SharedOrganizationsUrl, ShouldEqual, "/v2/private_domains/b2a35f0c-d5ad-4a59-bea7-461711d96b0d/shared_organizations")
	})
}

func TestCreateDomain(t *testing.T) {
	Convey("Create domain", t, func() {
		setup(MockRoute{"POST", "/v2/private_domains", postDomainPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		domain, err := client.CreateDomain("exmaple.com", "8483e4f1-d3a3-43e2-ab8c-b05ea40ef8db")
		So(err, ShouldBeNil)

		So(domain.Guid, ShouldEqual, "b98aeca1-22b9-49f9-8428-3ace9ea2ba11")
	})
}

func TestDeleteDomain(t *testing.T) {
	Convey("Delete domain", t, func() {
		setup(MockRoute{"DELETE", "/v2/private_domains/b2a35f0c-d5ad-4a59-bea7-461711d96b0d", "", ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteDomain("b2a35f0c-d5ad-4a59-bea7-461711d96b0d")
		So(err, ShouldBeNil)
	})
}
