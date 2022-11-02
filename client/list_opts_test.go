package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

func TestListOptions(t *testing.T) {
	newEmptyOpts := func() *AppListOptions {
		o := NewAppListOptions()
		o.Page = 0
		o.PerPage = 0
		return o
	}

	// defaults
	defaultOpts := NewAppListOptions()
	qs := defaultOpts.ToQueryString()
	require.Equal(t, "page=1&per_page=50", qs.Encode())

	// should not include zero values
	opts := newEmptyOpts()
	qs = opts.ToQueryString()
	require.Equal(t, "", qs.Encode())

	// single app by guid
	opts = newEmptyOpts()
	opts.GUIDs = Filter{
		Values: []string{"guid-1"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "guids="+url.QueryEscape("guid-1"), qs.Encode())

	// single app by name
	opts = newEmptyOpts()
	opts.Names = Filter{
		Values: []string{"app1"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "names="+url.QueryEscape("app1"), qs.Encode())

	// apps by org ids
	opts = newEmptyOpts()
	opts.OrganizationGUIDs = Filter{
		Values: []string{"org-guid-1", "org-guid-2"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "organization_guids="+url.QueryEscape("org-guid-1,org-guid-2"), qs.Encode())

	// apps by space ids
	opts = newEmptyOpts()
	opts.SpaceGUIDs = Filter{
		Values: []string{"space-guid-1"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "space_guids="+url.QueryEscape("space-guid-1"), qs.Encode())

	// apps by stacks
	opts = newEmptyOpts()
	opts.Stacks = Filter{
		Values: []string{"cflinuxfs2"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "stacks="+url.QueryEscape("cflinuxfs2"), qs.Encode())

	// multiple apps by name
	opts = newEmptyOpts()
	opts.Names = Filter{
		Values: []string{"app1", "app2"},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "names="+url.QueryEscape("app1,app2"), qs.Encode())

	// all apps but this one
	opts = newEmptyOpts()
	opts.Names = Filter{
		Values: []string{"app2"},
		Not:    true,
	}
	qs = opts.ToQueryString()
	require.Equal(t, url.QueryEscape("names[not]")+"="+url.QueryEscape("app2"), qs.Encode())

	// multiple dates
	time1, _ := time.Parse(time.RFC3339, "2016-03-18T00:00:00Z")
	time2, _ := time.Parse(time.RFC3339, "2016-10-17T00:00:00Z")
	opts = newEmptyOpts()
	opts.CreateAts = TimestampFilter{
		Timestamp: []time.Time{time1, time2},
	}
	qs = opts.ToQueryString()
	require.Equal(t, "created_ats="+url.QueryEscape("2016-03-18T00:00:00Z,2016-10-17T00:00:00Z"), qs.Encode())

	// gt date
	time1, _ = time.Parse(time.RFC3339, "2019-12-31T23:59:59Z")
	opts = newEmptyOpts()
	opts.CreateAts = TimestampFilter{
		Timestamp: []time.Time{time1},
		Operator:  GreaterThan,
	}
	qs = opts.ToQueryString()
	require.Equal(t, url.QueryEscape("created_ats[gt]")+"="+url.QueryEscape("2019-12-31T23:59:59Z"), qs.Encode())

	// lifecycle type
	opts = newEmptyOpts()
	opts.LifecycleType = resource.LifecycleBuildpack
	qs = opts.ToQueryString()
	require.Equal(t, "lifecycle_type="+url.QueryEscape("buildpack"), qs.Encode())

	// app include type
	optsInc := NewAppListOptions()
	optsInc.Include = resource.AppIncludeSpaceOrganization
	optsInc.Page = 0
	optsInc.PerPage = 0
	qs = optsInc.ToQueryString()
	require.Equal(t, "include="+url.QueryEscape("space.organization"), qs.Encode())
}
