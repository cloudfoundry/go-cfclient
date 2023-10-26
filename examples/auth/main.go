package main

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
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
	err := listOrganizationsWithConfig(config.NewFromCFHome)
	if err == nil {
		err = listOrganizationsWithConfig(func() (*config.Config, error) {
			return config.NewUserPassword("https://api.sys.example.com", "admin", "password")
		})
		if err == nil {
			err = listOrganizationsWithConfig(func() (*config.Config, error) {
				return config.NewClientSecret("https://api.sys.example.com", "cf-client", "client-secret")
			})
		}
	}
	return err
}

func listOrganizationsWithConfig(cFn func() (*config.Config, error)) error {
	ctx := context.Background()
	conf, err := cFn()
	if err != nil {
		return err
	}
	conf.WithSkipTLSValidation(true)
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	err = listOrganizations(ctx, cf)
	if err != nil {
		return err
	}

	return nil
}

func listOrganizations(ctx context.Context, cf *client.Client) error {
	// grab the first space
	orgs, err := cf.Organizations.ListAll(ctx, nil)
	if err != nil {
		return err
	}

	for _, o := range orgs {
		fmt.Printf("%s\n", o.Name)
	}
	return nil
}
