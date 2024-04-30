package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
)

const apiURL = "https://api.sys.example.com"
const loginURL = "https://login.sys.example.com"
const tokenURL = "https://uaa.sys.example.com"
const username = "admin"
const password = "password"
const clientID = "cf"
const clientSecret = "secret"
const accessToken = "<access-token>"
const refreshToken = "<refresh-token>"

func main() {
	err := execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func execute() error {
	// use the CF CLI config and the stored access/refresh token
	cfg, err := config.NewFromCFHome()
	if err != nil {
		return err
	}
	err = listOrganizationsWithConfig(cfg)
	if err != nil {
		return err
	}

	// use the hardcoded CF API endpoint and user/pass and skip TLS validation
	cfg, err = config.New(apiURL,
		config.UserPassword(username, password),
		config.SkipTLSValidation())
	if err != nil {
		return err
	}
	err = listOrganizationsWithConfig(cfg)
	if err != nil {
		return err
	}

	// use the hardcoded CF API endpoint and client/secret and skip TLS validation
	cfg, err = config.New(apiURL,
		config.ClientCredentials(clientID, clientSecret),
		config.SkipTLSValidation())
	if err != nil {
		return err
	}
	err = listOrganizationsWithConfig(cfg)
	if err != nil {
		return err
	}

	// use the hardcoded CF API endpoint and OAuth token
	cfg, err = config.New(apiURL,
		config.Token(accessToken, refreshToken),
		config.SkipTLSValidation())
	if err != nil {
		return err
	}
	err = listOrganizationsWithConfig(cfg)
	if err != nil {
		return err
	}

	// Unnecessarily use all config options
	cfg, err = config.New(apiURL,
		config.UserPassword(username, password),
		config.SkipTLSValidation(),
		config.AuthTokenURL(loginURL, tokenURL),
		config.UserAgent("MyApp-Client/1.0"),
		config.HttpClient(&http.Client{}),
		config.RequestTimeout(10*time.Second),
		config.Origin("uaa"),
		config.Scopes("cloud_controller.read", "cloud_controller_service_permissions.read"),
		config.SSHOAuthClient("ssh-proxy"))
	if err != nil {
		return err
	}
	return listOrganizationsWithConfig(cfg)
}

func listOrganizationsWithConfig(cfg *config.Config) error {
	ctx := context.Background()
	cf, err := client.New(cfg)
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
