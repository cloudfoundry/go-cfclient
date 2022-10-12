package test

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/client"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	c := createClient(t)
	org := getOrg(t, c)
	fmt.Printf("%v\n", org)
}

func getOrg(t *testing.T, c *client.Client) *resource.Organization {
	opts := client.NewOrgListOptions()
	opts.Names = client.Filter{
		Values: []string{"e2e-test-org"},
	}
	orgs, _, err := c.Organizations.List(opts)
	require.NoError(t, err)

	var org *resource.Organization
	if len(orgs) > 0 {
		org = orgs[0]
	} else {
		oc := &resource.OrganizationCreate{
			Name: "e2e-test-org",
		}
		org, err = c.Organizations.Create(oc)
		require.NoError(t, err)
	}
	require.Equal(t, "e2e-test-org", org.Name)
	require.NotEmpty(t, org.GUID)
	require.NotEmpty(t, org.CreatedAt)
	require.NotEmpty(t, org.UpdatedAt)
	return org
}

func getSpace(t *testing.T, c *client.Client) *resource.Space {
	// TODO find/create space
	return nil
}

func createClient(t *testing.T) *client.Client {
	config, err := client.NewConfigFromCFHome()
	require.NoError(t, err)
	config.SkipSSLValidation(true)
	c, err := client.New(config)
	require.NoError(t, err)
	return c
}
