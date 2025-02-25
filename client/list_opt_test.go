package client_test

import (
	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"net/url"
	"testing"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/stretchr/testify/require"
)

func TestListOptions(t *testing.T) {
	newEmptyOpts := func() *client.AppListOptions {
		o := client.NewAppListOptions()
		o.Page = 0
		o.PerPage = 0
		return o
	}

	// defaults
	defaultOpts := client.NewAppListOptions()
	require.Equal(t, "page=1&per_page=50", qs(defaultOpts))

	// should not include zero values
	opts := newEmptyOpts()
	require.Equal(t, "", qs(opts))

	// single app by guid
	opts = newEmptyOpts()
	opts.GUIDs.EqualTo("guid-1")
	require.Equal(t, "guids=guid-1", qs(opts))

	// single app by name
	opts = newEmptyOpts()
	opts.Names.EqualTo("app1")
	require.Equal(t, "names=app1", qs(opts))

	// apps by organization ids
	opts = newEmptyOpts()
	opts.OrganizationGUIDs.EqualTo("organization-guid-1", "organization-guid-2")
	require.Equal(t, "organization_guids=organization-guid-1,organization-guid-2", qs(opts))

	// apps by space ids
	opts = newEmptyOpts()
	opts.SpaceGUIDs.EqualTo("space-guid-1")
	require.Equal(t, "space_guids=space-guid-1", qs(opts))

	// apps by stacks
	opts = newEmptyOpts()
	opts.Stacks.EqualTo("cflinuxfs2")
	require.Equal(t, "stacks=cflinuxfs2", qs(opts))

	// multiple apps by name
	opts = newEmptyOpts()
	opts.Names.EqualTo("app1", "app2")
	require.Equal(t, "names=app1,app2", qs(opts))

	// exclude filter
	auditOpts := client.NewAuditEventListOptions()
	auditOpts.Page = 0
	auditOpts.PerPage = 0
	auditOpts.TargetGUIDs.NotEqualTo("app2")
	require.Equal(t, "target_guids[not]=app2", qsAuditEvents(auditOpts))

	// multiple exact dates
	opts = newEmptyOpts()
	opts.CreateAts.EqualTo(date("2016-03-18T00:00:00Z"), date("2016-10-17T00:00:00Z"))
	require.Equal(t, "created_ats=2016-03-18T00:00:00Z,2016-10-17T00:00:00Z", qs(opts))

	// gt date
	opts = newEmptyOpts()
	opts.CreateAts.After(date("2019-12-31T23:59:59Z"))
	require.Equal(t, "created_ats[gt]=2019-12-31T23:59:59Z", qs(opts))

	// date range
	opts = newEmptyOpts()
	opts.CreateAts.After(date("2025-02-06T13:00:00Z"))
	opts.CreateAts.Before(date("2025-02-06T14:00:00Z"))
	require.Equal(t, "created_ats[gt]=2025-02-06T13:00:00Z&created_ats[lt]=2025-02-06T14:00:00Z", qs(opts))

	// lifecycle type
	opts = newEmptyOpts()
	opts.LifecycleType = resource.LifecycleBuildpack
	require.Equal(t, "lifecycle_type=buildpack", qs(opts))

	// app include type
	opts = newEmptyOpts()
	opts.Include = resource.AppIncludeSpaceOrganization
	require.Equal(t, "include=space.organization", qs(opts))
}

func date(v string) time.Time {
	time1, _ := time.Parse(time.RFC3339, v)
	return time1
}

func qs(opts *client.AppListOptions) string {
	values, _ := opts.ToQueryString()
	u, _ := url.Parse(values.Encode())
	return u.Path
}

func qsAuditEvents(opts *client.AuditEventListOptions) string {
	values, _ := opts.ToQueryString()
	u, _ := url.Parse(values.Encode())
	return u.Path
}
