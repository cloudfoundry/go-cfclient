package main

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/operation"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("expected arguments: org, space, /path/to/spring-music.jar")
		os.Exit(1)
	}
	org := os.Args[1]
	space := os.Args[2]
	pathToZip := os.Args[3]

	err := runPush(org, space, pathToZip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func runPush(org, space, pathToZip string) error {
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

	var manifest *operation.Manifest
	err = yaml.Unmarshal([]byte(yamlManifest), &manifest)
	if err != nil {
		return err
	}

	zipFile, err := os.Open(pathToZip)
	if err != nil {
		return err
	}
	pushOp := operation.NewAppPushOperation(cf, org, space)
	app, err := pushOp.Push(ctx, manifest.Applications[0], zipFile)
	if err != nil {
		return err
	}
	fmt.Printf("successfully pushed %s, state: %s\n", app.Name, app.State)
	return nil
}

const yamlManifest = `
---
applications:
- name: spring-music-example
  memory: 1G
  random-route: true
  stack: cflinuxfs3
  buildpacks:
  - java_buildpack_offline
`
