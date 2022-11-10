package main

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"os"
)

func main() {
	err := runPush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func runPush() error {
	conf, err := client.NewConfigFromCFHome()
	if err != nil {
		return err
	}
	conf.SkipSSLValidation(true)
	cf, err := client.New(conf)
	if err != nil {
		return err
	}

	spaceOpts := client.NewSpaceListOptions()
	spaceOpts.Names.Values = []string{"dev"}
	spaces, err := cf.Spaces.ListAll(spaceOpts)
	if err != nil {
		return err
	}
	if len(spaces) != 1 {
		return fmt.Errorf("expected one space named dev, but found %d", len(spaces))
	}
	space := spaces[0]

	jobGUID, err := cf.Manifests.ApplyManifest(space.GUID, manifest)
	if err != nil {
		return err
	}
	err = cf.Jobs.PollComplete(jobGUID, nil)
	if err != nil {
		return err
	}

	appOpts := client.NewAppListOptions()
	appOpts.Names.Values = []string{"spring-music"}
	apps, err := cf.Applications.ListAll(appOpts)
	if err != nil {
		return err
	}
	if len(apps) != 1 {
		return fmt.Errorf("expected one space named dev, but found %d", len(apps))
	}
	app := apps[0]

	newPkg := resource.NewPackageCreate(app.GUID)
	pkg, err := cf.Packages.Create(newPkg)
	if err != nil {
		return err
	}

	// might need full path or copy it locally
	f, err := os.Open("spring-music-1.0.jar")
	if err != nil {
		return err
	}
	err = cf.Packages.UploadBits(pkg.GUID, f)
	if err != nil {
		return err
	}
	err = cf.Packages.PollReady(pkg.GUID, nil)
	if err != nil {
		return err
	}

	newBuild := resource.NewBuildCreate(pkg.GUID)
	newBuild.Lifecycle = &resource.Lifecycle{
		Type: "buildpack",
		BuildpackData: resource.BuildpackLifecycle{
			Buildpacks: []string{"java_buildpack_offline"},
			Stack:      "cflinuxfs3",
		},
	}
	build, err := cf.Builds.Create(newBuild)
	if err != nil {
		return err
	}
	err = cf.Builds.PollStaged(build.GUID, nil)
	if err != nil {
		return err
	}

	droplets, err := cf.Droplets.ListForPackageAll(pkg.GUID, nil)
	if err != nil {
		return err
	}
	if len(droplets) != 1 {
		return fmt.Errorf("expected one droplet, but found %d", len(droplets))
	}
	droplet := droplets[0]

	_, err = cf.Droplets.SetCurrentAssociationForApp(app.GUID, droplet.GUID)
	if err != nil {
		return err
	}

	app, err = cf.Applications.Start(app.GUID)
	if err != nil {
		return err
	}

	fmt.Printf("Finished pushing %s\n", app.Name)

	return nil
}

const manifest = `
---
applications:
- name: spring-music
  memory: 1G
  random-route: true
  stack: cflinuxfs3
  buildpacks:
  - java_buildpack_offline
`
