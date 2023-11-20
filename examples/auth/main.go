package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
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
		err = listOrganizationsWithConfig(func(options ...config.Option) (*config.Config, error) {
			options = append(options, config.UserPassword("admin", "password"))
			return config.New("https://api.sys.example.com", options...)
		})
		if err == nil {
			err = listOrganizationsWithConfig(func(options ...config.Option) (*config.Config, error) {
				options = append(options, config.ClientCredentials("cf-client", "client-secret"))
				return config.New("https://api.sys.example.com", options...)
			})
		}
	}
	return err
}

func listOrganizationsWithConfig(cFn func(options ...config.Option) (*config.Config, error)) error {
	ctx := context.Background()
	conf, err := cFn(config.SkipTLSValidation())
	if err != nil {
		return err
	}
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
