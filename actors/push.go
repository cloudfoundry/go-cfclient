package actors

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/client"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"io"
	"strconv"
	"strings"
	"time"
)

type AppPusher struct {
	orgName   string
	spaceName string
	client    *client.Client
}

func NewAppPusher(client *client.Client, orgName, spaceName string) *AppPusher {
	return &AppPusher{
		orgName:   orgName,
		spaceName: spaceName,
		client:    client,
	}
}

func (p *AppPusher) Push(appManifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	org, err := p.findOrg()
	if err != nil {
		return nil, err
	}
	space, err := p.findSpace(org.GUID)
	if err != nil {
		return nil, err
	}
	app, err := p.findAppOrNil(org.GUID, space.GUID, appManifest.Name)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return p.pushNewApp(org, space, appManifest, zipFile)
	} else {
		return p.pushUpdatedApp(org, space, app, appManifest, zipFile)
	}
}

func (p *AppPusher) pushUpdatedApp(org *resource.Organization, space *resource.Space, app *resource.App, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	panic("pushUpdatedApp not implemented")
}

// pushNewApp pushes a new application that doesn't already exist
//
// After an application is created and packages are uploaded, a droplet must be created via a build in order for
// an application to be deployed or tasks to be run. The current droplet must be assigned to an application before
// it may be started. When tasks are created, they either use a specific droplet guid, or use the current droplet
// assigned to an application.
func (p *AppPusher) pushNewApp(org *resource.Organization, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	// Create the application object
	newApp := resource.NewAppCreate(manifest.Name, space.GUID)
	newApp.EnvironmentVariables = manifest.Env
	newApp.Lifecycle.Type = "buildpack"
	newApp.Lifecycle.BuildpackData = resource.BuildpackLifecycle{
		Buildpacks: manifest.Buildpacks,
		Stack:      manifest.Stack,
	}
	app, err := p.client.Applications.Create(newApp)
	if err != nil {
		return nil, fmt.Errorf("could not create new app %s: %w", manifest.Name, err)
	}

	// create a package and then upload package bits
	newPkg := resource.NewPackageCreate(app.GUID)
	pkg, err := p.client.Packages.Create(newPkg)
	if err != nil {
		return nil, fmt.Errorf("could not create new app %s package: %w", manifest.Name, err)
	}
	err = p.client.Packages.UploadBits(pkg.GUID, zipFile)
	if err != nil {
		return nil, fmt.Errorf("could not upload app %s package bits: %w", manifest.Name, err)
	}

	// build droplet
	newBuild := resource.NewBuildCreate(pkg.GUID)
	newBuild.Lifecycle = &resource.Lifecycle{
		Type: "buildpack",
		BuildpackData: resource.BuildpackLifecycle{
			Buildpacks: manifest.Buildpacks,
			Stack:      manifest.Stack,
		},
	}
	build, err := p.client.Builds.Create(newBuild)
	if err != nil {
		return nil, fmt.Errorf("could not create build for app %s: %w", manifest.Name, err)
	}

	// wait for build to finish
	done := false
	for !done {
		time.Sleep(time.Second * 2)
		b, err := p.client.Builds.Get(build.GUID)
		if err != nil {
			return nil, fmt.Errorf("could not get build for app %s: %w", manifest.Name, err)
		}
		switch b.State {
		case "STAGING":
			continue
		case "STAGED":
			done = true
		case "FAILED":
			return nil, fmt.Errorf("failed to stage app %s: %w", manifest.Name, err)
		}
	}

	// set the app's process attributes
	processes, err := p.client.Processes.ListForAppAll(app.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get processes for app %s: %w", manifest.Name, err)
	}
	if len(processes) != 1 {
		return nil, fmt.Errorf("expected one process for app %s but got %d", manifest.Name, len(processes))
	}
	process := processes[0]

	processUpdate := resource.NewProcessUpdate()
	if manifest.Command != "" {
		processUpdate.WithCommand(manifest.Command)
	}
	if manifest.HealthCheckHTTPEndpoint != "" {
		processUpdate.WithHealthCheckEndpoint(manifest.HealthCheckHTTPEndpoint)
	}
	if manifest.Timeout != 0 {
		processUpdate.WithHealthCheckInvocationTimeout(manifest.Timeout)
	}
	if manifest.HealthCheckType != "" {
		processUpdate.WithHealthCheckType(manifest.HealthCheckType)
	}
	process, err = p.client.Processes.Update(process.GUID, processUpdate)
	if err != nil {
		return nil, fmt.Errorf("could not update process for app %s: %w", manifest.Name, err)
	}

	processScale := resource.NewProcessScale()
	if manifest.Memory != "" {
		mb, err := sizeStringToIntMB(manifest.Memory)
		if err != nil {
			return nil, err
		}
		processScale.WithMemoryInMB(mb)
	}
	if manifest.DiskQuota != "" {
		mb, err := sizeStringToIntMB(manifest.DiskQuota)
		if err != nil {
			return nil, err
		}
		processScale.WithDiskInMB(mb)
	}
	if manifest.Instances > 0 {
		processScale.WithInstances(manifest.Instances)
	}
	if manifest.LogRateLimit != "" {
		// TODO add a type/func that can natively parse and convert 1MB 1GB etc
		//processScale.WithLogRateLimitInBytesPerSecond(manifest.LogRateLimit)
	}
	process, err = p.client.Processes.Scale(process.GUID, processScale)
	if err != nil {
		return nil, fmt.Errorf("could not scale process for app %s: %w", manifest.Name, err)
	}

	// finally start the app
	app, err = p.client.Applications.Start(app.GUID)
	if err != nil {
		return nil, fmt.Errorf("could not start the app %s: %w", manifest.Name, err)
	}

	return app, nil
}

func (p *AppPusher) findOrg() (*resource.Organization, error) {
	opts := client.NewOrgListOptions()
	opts.Names.Values = []string{p.orgName}
	orgs, err := p.client.Organizations.ListAll(opts)
	if err != nil {
		return nil, fmt.Errorf("could not find org %s: %w", p.orgName, err)
	}
	if len(orgs) != 1 {
		return nil, fmt.Errorf("expected to find one org named %s, but found %d", p.orgName, len(orgs))
	}
	return orgs[0], nil
}

func (p *AppPusher) findSpace(orgGUID string) (*resource.Space, error) {
	opts := client.NewSpaceListOptions()
	opts.Names.Values = []string{p.spaceName}
	opts.OrganizationGUIDs.Values = []string{orgGUID}
	spaces, err := p.client.Spaces.ListAll(opts)
	if err != nil {
		return nil, fmt.Errorf("could not find space %s: %w", p.spaceName, err)
	}
	if len(spaces) != 1 {
		return nil, fmt.Errorf("expected to find one space named %s, but found %d", p.spaceName, len(spaces))
	}
	return spaces[0], nil
}

func (p *AppPusher) findAppOrNil(orgGUID, spaceGUID, appName string) (*resource.App, error) {
	opts := client.NewAppListOptions()
	opts.OrganizationGUIDs.Values = []string{orgGUID}
	opts.SpaceGUIDs.Values = []string{spaceGUID}
	opts.Names.Values = []string{appName}
	apps, err := p.client.Applications.ListAll(opts)
	if err != nil {
		return nil, fmt.Errorf("could not find app %s: %w", appName, err)
	}
	if len(apps) == 1 {
		return apps[0], nil
	}
	return nil, nil
}

func sizeStringToIntMB(size string) (int, error) {
	size = strings.ToUpper(size)
	if strings.HasSuffix(size, "M") || strings.HasSuffix(size, "MB") {
		s := strings.TrimSuffix(strings.TrimSuffix(size, "M"), "MB")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("could not convert MB size string to int: %w", err)
		}
		return i, nil
	} else if strings.HasSuffix(size, "G") || strings.HasSuffix(size, "GB") {
		s := strings.TrimSuffix(strings.TrimSuffix(size, "G"), "GB")
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("could not convert GB size string to int: %w", err)
		}
		i = i * 1024
		return i, nil
	}
	return 0, fmt.Errorf("unsupported size string %s", size)
}
