package main

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"os"
)

func main() {
	err := execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func execute() error {
	ctx := context.Background()
	conf, err := config.NewFromCFHome()
	if err != nil {
		return err
	}
	conf.WithSkipTLSValidation(true)
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	err = listSpaceDevsInSpace(ctx, cf)
	if err != nil {
		return err
	}
	return listAllSpaceDevelopers(ctx, cf)
}

func listSpaceDevsInSpace(ctx context.Context, cf *client.Client) error {
	// grab the first space
	spaces, _, err := cf.Spaces.List(ctx, nil)
	if err != nil {
		return err
	}
	space := spaces[0]

	// list space developer roles and users in the space
	opts := client.NewRoleListOptions()
	opts.Space(space.GUID)
	opts.SpaceRoleType(resource.SpaceRoleDeveloper)
	roles, users, err := cf.Roles.ListIncludeUsersAll(ctx, opts)
	if err != nil {
		return err
	}
	for _, r := range roles {
		fmt.Printf("%s - %s\n", r.Type, r.GUID)
	}
	for _, u := range users {
		fmt.Printf("%s - %s\n", u.Username, u.GUID)
	}

	return nil
}

func listAllSpaceDevelopers(ctx context.Context, cf *client.Client) error {
	opts := client.NewRoleListOptions()
	opts.SpaceRoleType(resource.SpaceRoleDeveloper)
	_, users, err := cf.Roles.ListIncludeUsersAll(ctx, opts)
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Printf("%s - %s\n", u.Username, u.GUID)
	}

	return nil
}
