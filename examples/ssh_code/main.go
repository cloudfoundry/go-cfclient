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
	ctx := context.Background()
	conf, err := config.NewFromCFHome(config.SkipTLSValidation())
	if err != nil {
		return err
	}
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	code, err := cf.SSHCode(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("SSH Code: %s\n", code)
	return nil
}
