package client_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
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
	qs := defaultOpts.ToQueryString()
	require.Equal(t, "page=1&per_page=50", qs.Encode())

	// should not include zero values
	opts := newEmptyOpts()
	qs = opts.ToQueryString()
	require.Equal(t, "", qs.Encode())

	// single app by guid
	opts = newEmptyOpts()
	opts.GUIDs.EqualTo("guid-1")
	qs = opts.ToQueryString()
	require.Equal(t, "guids="+url.QueryEscape("guid-1"), qs.Encode())

	// single app by name
	opts = newEmptyOpts()
	opts.Names.EqualTo("app1")
	qs = opts.ToQueryString()
	require.Equal(t, "names="+url.QueryEscape("app1"), qs.Encode())

	// apps by organization ids
	opts = newEmptyOpts()
	opts.OrganizationGUIDs.EqualTo("organization-guid-1", "organization-guid-2")
	qs = opts.ToQueryString()
	require.Equal(t, "organization_guids="+url.QueryEscape("organization-guid-1,organization-guid-2"), qs.Encode())

	// apps by space ids
	opts = newEmptyOpts()
	opts.SpaceGUIDs.EqualTo("space-guid-1")
	qs = opts.ToQueryString()
	require.Equal(t, "space_guids="+url.QueryEscape("space-guid-1"), qs.Encode())

	// apps by stacks
	opts = newEmptyOpts()
	opts.Stacks.EqualTo("cflinuxfs2")
	qs = opts.ToQueryString()
	require.Equal(t, "stacks="+url.QueryEscape("cflinuxfs2"), qs.Encode())

	// multiple apps by name
	opts = newEmptyOpts()
	opts.Names.EqualTo("app1", "app2")
	qs = opts.ToQueryString()
	require.Equal(t, "names="+url.QueryEscape("app1,app2"), qs.Encode())

	// all apps but this one
	opts = newEmptyOpts()
	opts.Names.NotEqualTo("app2")
	qs = opts.ToQueryString()
	require.Equal(t, url.QueryEscape("names[not]")+"="+url.QueryEscape("app2"), qs.Encode())

	// multiple dates
	opts = newEmptyOpts()
	opts.CreateAts.EqualTo(date("2016-03-18T00:00:00Z"), date("2016-10-17T00:00:00Z"))
	qs = opts.ToQueryString()
	require.Equal(t, "created_ats="+url.QueryEscape("2016-03-18T00:00:00Z,2016-10-17T00:00:00Z"), qs.Encode())

	// gt date
	opts = newEmptyOpts()
	opts.CreateAts.After(date("2019-12-31T23:59:59Z"))
	qs = opts.ToQueryString()
	require.Equal(t, url.QueryEscape("created_ats[gt]")+"="+url.QueryEscape("2019-12-31T23:59:59Z"), qs.Encode())

	// lifecycle type
	opts = newEmptyOpts()
	opts.LifecycleType = resource.LifecycleBuildpack
	qs = opts.ToQueryString()
	require.Equal(t, "lifecycle_type="+url.QueryEscape("buildpack"), qs.Encode())

	// app include type
	optsInc := newEmptyOpts()
	optsInc.Include = resource.AppIncludeSpaceOrganization
	qs = optsInc.ToQueryString()
	require.Equal(t, "include="+url.QueryEscape("space.organization"), qs.Encode())
}

func date(v string) time.Time {
	time1, _ := time.Parse(time.RFC3339, v)
	return time1
}
