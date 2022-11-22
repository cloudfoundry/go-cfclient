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

	code, err := cf.SSHCode(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("SSH Code: %s\n", code)
	return nil
}
