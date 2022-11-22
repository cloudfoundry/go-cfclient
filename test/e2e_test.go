//go:build integration
// +build integration

package test

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	OrgName   = "go-cfclient-e2e"
	SpaceName = "go-cfclient-e2e"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	c := createClient(t)

	// get the org with the access token
	org := getOrg(t, ctx, c)
	fmt.Println(org.Name)

	// try to get the space
	space := getSpace(t, ctx, c, org)
	fmt.Println(space.Name)
}

func getOrg(t *testing.T, ctx context.Context, c *client.Client) *resource.Organization {
	opts := client.NewOrgListOptions()
	opts.Names = client.Filter{
		Values: []string{OrgName},
	}
	orgs, _, err := c.Organizations.List(ctx, opts)
	require.NoError(t, err)

	var org *resource.Organization
	if len(orgs) > 0 {
		org = orgs[0]
	} else {
		oc := &resource.OrganizationCreate{
			Name: OrgName,
		}
		org, err = c.Organizations.Create(ctx, oc)
		require.NoError(t, err)
	}
	require.Equal(t, OrgName, org.Name)
	require.NotEmpty(t, org.GUID)
	require.NotEmpty(t, org.CreatedAt)
	require.NotEmpty(t, org.UpdatedAt)
	return org
}

func getSpace(t *testing.T, ctx context.Context, c *client.Client, org *resource.Organization) *resource.Space {
	opts := client.NewSpaceListOptions()
	opts.Names = client.Filter{
		Values: []string{SpaceName},
	}
	spaces, _, err := c.Spaces.List(ctx, opts)
	require.NoError(t, err)

	var space *resource.Space
	if len(spaces) > 0 {
		space = spaces[0]
	} else {
		sc := resource.NewSpaceCreate(SpaceName, org.GUID)
		space, err = c.Spaces.Create(ctx, sc)
		require.NoError(t, err)
	}
	require.Equal(t, SpaceName, space.Name)
	require.NotEmpty(t, space.GUID)
	require.NotEmpty(t, space.CreatedAt)
	require.NotEmpty(t, space.UpdatedAt)
	return space
}

func createClient(t *testing.T) *client.Client {
	cfg, err := config.NewFromCFHome()
	require.NoError(t, err)
	cfg.WithSkipTLSValidation(true)
	c, err := client.New(cfg)
	require.NoError(t, err)
	return c
}
