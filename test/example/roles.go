package main

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/client"
	"github.com/cloudfoundry-community/go-cfclient/resource"
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
	conf, err := client.NewConfigFromCFHome()
	if err != nil {
		return err
	}
	conf.SkipSSLValidation(true)
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	err = listSpaceDevsInSpace(cf)
	if err != nil {
		return err
	}
	return listAllSpaceDevelopers(cf)
}

func listSpaceDevsInSpace(cf *client.Client) error {
	// grab the first space
	spaces, _, err := cf.Spaces.List(nil)
	if err != nil {
		return err
	}
	space := spaces[0]

	// list space developer roles and users in the space
	opts := client.NewRoleListIncludeOptions(resource.RoleIncludeUser)
	opts.Space(space.GUID)
	opts.SpaceRoleType(resource.SpaceRoleDeveloper)
	roles, include, err := cf.Roles.ListIncludeAll(opts)
	if err != nil {
		return err
	}
	for _, r := range roles {
		fmt.Printf("%s - %s\n", r.Type, r.GUID)
	}
	for _, u := range include.Users {
		fmt.Printf("%s - %s\n", u.Username, u.GUID)
	}

	return nil
}

func listAllSpaceDevelopers(cf *client.Client) error {
	opts := client.NewRoleListIncludeOptions(resource.RoleIncludeUser)
	opts.SpaceRoleType(resource.SpaceRoleDeveloper)
	_, include, err := cf.Roles.ListIncludeAll(opts)
	if err != nil {
		return err
	}
	for _, u := range include.Users {
		fmt.Printf("%s - %s\n", u.Username, u.GUID)
	}

	return nil
}
