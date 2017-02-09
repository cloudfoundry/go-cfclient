package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListOrgs(t *testing.T) {
	Convey("List Org", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/organizations", listOrgsPayload, ""},
			{"GET", "/v2/orgsPage2", listOrgsPayloadPage2, ""},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		orgs, err := client.ListOrgs()
		So(err, ShouldBeNil)

		So(len(orgs), ShouldEqual, 4)
		So(orgs[0].Guid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(orgs[0].Name, ShouldEqual, "demo")
	})
}

func TestGetOrgByGuid(t *testing.T) {
	Convey("List Org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/1c0e6074-777f-450e-9abc-c42f39d9b75b", orgByGuidPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		org, err := client.GetOrgByGuid("1c0e6074-777f-450e-9abc-c42f39d9b75b")
		So(err, ShouldBeNil)

		So(org.Guid, ShouldEqual, "1c0e6074-777f-450e-9abc-c42f39d9b75b")
		So(org.Name, ShouldEqual, "name-1716")
	})
}

func TestOrgSpaces(t *testing.T) {
	Convey("Get spaces by org", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/foo/spaces", orgSpacesPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.OrgSpaces("foo")
		So(err, ShouldBeNil)

		So(len(spaces), ShouldEqual, 1)
		So(spaces[0].Guid, ShouldEqual, "b8aff561-175d-45e8-b1e7-67e2aedb03b6")
		So(spaces[0].Name, ShouldEqual, "test")
	})
}

func TestOrgSummary(t *testing.T) {
	Convey("Get org summary", t, func() {
		setup(MockRoute{"GET", "/v2/organizations/06dcedd4-1f24-49a6-adc1-cce9131a1b2c/summary", orgSummaryPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		org := &Org{
			Guid: "06dcedd4-1f24-49a6-adc1-cce9131a1b2c",
			c:    client,
		}
		summary, err := org.Summary()
		So(err, ShouldBeNil)

		So(summary.Guid, ShouldEqual, "06dcedd4-1f24-49a6-adc1-cce9131a1b2c")
		So(summary.Name, ShouldEqual, "system")
		So(summary.Status, ShouldEqual, "active")

		spaces := summary.Spaces
		So(len(spaces), ShouldEqual, 1)
		So(spaces[0].Guid, ShouldEqual, "494d8b64-8181-4183-a6d3-6279db8fec6e")
		So(spaces[0].Name, ShouldEqual, "test")
		So(spaces[0].ServiceCount, ShouldEqual, 1)
		So(spaces[0].AppCount, ShouldEqual, 2)
		So(spaces[0].MemDevTotal, ShouldEqual, 32)
		So(spaces[0].MemProdTotal, ShouldEqual, 64)
	})
}

func TestOrgQuota(t *testing.T) {
	Convey("Get org quota", t, func() {
		setup(MockRoute{"GET", "/v2/quota_definitions/a537761f-9d93-4b30-af17-3d73dbca181b", orgQuotaPayload, ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		org := &Org{
			QuotaDefinitionGuid: "a537761f-9d93-4b30-af17-3d73dbca181b",
			c:                   client,
		}
		orgQuota, err := org.Quota()
		So(err, ShouldBeNil)

		So(orgQuota.Guid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(orgQuota.Name, ShouldEqual, "test-2")
		So(orgQuota.NonBasicServicesAllowed, ShouldEqual, false)
		So(orgQuota.TotalServices, ShouldEqual, 10)
		So(orgQuota.TotalRoutes, ShouldEqual, 20)
		So(orgQuota.TotalPrivateDomains, ShouldEqual, 30)
		So(orgQuota.MemoryLimit, ShouldEqual, 40)
		So(orgQuota.TrialDBAllowed, ShouldEqual, true)
		So(orgQuota.InstanceMemoryLimit, ShouldEqual, 50)
		So(orgQuota.AppInstanceLimit, ShouldEqual, 60)
		So(orgQuota.AppTaskLimit, ShouldEqual, 70)
		So(orgQuota.TotalServiceKeys, ShouldEqual, 80)
		So(orgQuota.TotalReservedRoutePorts, ShouldEqual, 90)
	})
}

func TestDeleteOrg(t *testing.T) {
	Convey("Delete org", t, func() {
		setup(MockRoute{"DELETE", "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b", "", ""}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		err = client.DeleteOrg("a537761f-9d93-4b30-af17-3d73dbca181b")
		So(err, ShouldBeNil)
	})
}
