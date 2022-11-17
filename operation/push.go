package operation

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"gopkg.in/yaml.v3"
	"io"
)

// AppPushOperation can be used to push buildpack apps
type AppPushOperation struct {
	orgName   string
	spaceName string
	client    *client.Client
}

// NewAppPushOperation creates a new AppPushOperation
func NewAppPushOperation(client *client.Client, orgName, spaceName string) *AppPushOperation {
	return &AppPushOperation{
		orgName:   orgName,
		spaceName: spaceName,
		client:    client,
	}
}

// Push creates or updates an application using the specified manifest and zipped source files
func (p *AppPushOperation) Push(appManifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	org, err := p.findOrg()
	if err != nil {
		return nil, err
	}
	space, err := p.findSpace(org.GUID)
	if err != nil {
		return nil, err
	}
	return p.pushApp(space, appManifest, zipFile)
}

// pushApp pushes an application
//
// After an application is created and packages are uploaded, a droplet must be created via a build in order for
// an application to be deployed or tasks to be run. The current droplet must be assigned to an application before
// it may be started. When tasks are created, they either use a specific droplet guid, or use the current droplet
// assigned to an application.
func (p *AppPushOperation) pushApp(space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	err := p.applySpaceManifest(space, manifest)
	if err != nil {
		return nil, err
	}

	app, err := p.findApp(manifest.Name, space)
	if err != nil {
		return nil, err
	}

	pkg, err := p.uploadPackage(app, zipFile)
	if err != nil {
		return nil, err
	}

	droplet, err := p.buildDroplet(pkg, manifest)
	if err != nil {
		return nil, err
	}

	_, err = p.client.Droplets.SetCurrentAssociationForApp(app.GUID, droplet.GUID)
	if err != nil {
		return nil, err
	}

	return p.client.Applications.Start(app.GUID)
}

func (p *AppPushOperation) applySpaceManifest(space *resource.Space, manifest *AppManifest) error {
	// wrap it in a manifest that has an applications array as required by the API
	multiAppsManifest := &Manifest{
		Applications: []*AppManifest{manifest},
	}
	manifestBytes, err := yaml.Marshal(&multiAppsManifest)
	if err != nil {
		return fmt.Errorf("error marshalling application manifest: %w", err)
	}

	jobGUID, err := p.client.Manifests.ApplyManifest(space.GUID, string(manifestBytes))
	if err != nil {
		return fmt.Errorf("error applying application manifest to space %s: %w", space.Name, err)
	}
	err = p.client.Jobs.PollComplete(jobGUID, nil)
	if err != nil {
		return fmt.Errorf("error waiting for application manifest to finish applying to space %s: %w", space.Name, err)
	}
	return nil
}

func (p *AppPushOperation) findApp(appName string, space *resource.Space) (*resource.App, error) {
	appOpts := client.NewAppListOptions()
	appOpts.Names.Values = []string{appName}
	appOpts.SpaceGUIDs.Values = []string{space.GUID}
	apps, err := p.client.Applications.ListAll(appOpts)
	if err != nil {
		return nil, err
	}
	if len(apps) != 1 {
		return nil, fmt.Errorf("expected to find one application named %s in space %s, but found %d",
			appName, space.Name, len(apps))
	}
	return apps[0], nil
}

func (p *AppPushOperation) uploadPackage(app *resource.App, zipFile io.Reader) (*resource.Package, error) {
	newPkg := resource.NewPackageCreate(app.GUID)
	pkg, err := p.client.Packages.Create(newPkg)
	if err != nil {
		return nil, fmt.Errorf("error creating package for app %s: %w", app.Name, err)
	}

	err = p.client.Packages.UploadBits(pkg.GUID, zipFile)
	if err != nil {
		return nil, fmt.Errorf("error uploading package bits for app %s: %w", app.Name, err)
	}
	err = p.client.Packages.PollReady(pkg.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for package to process for app %s: %w", app.Name, err)
	}
	return pkg, nil
}

func (p *AppPushOperation) buildDroplet(pkg *resource.Package, manifest *AppManifest) (*resource.Droplet, error) {
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
		return nil, fmt.Errorf("error creating build from package for app %s: %w", manifest.Name, err)
	}
	err = p.client.Builds.PollStaged(build.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for app %s package to build: %w", manifest.Name, err)
	}

	opts := client.NewDropletPackageListOptions()
	opts.States.Values = []string{string(resource.DropletStateStaged)}
	droplets, err := p.client.Droplets.ListForPackageAll(pkg.GUID, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding droplet for app %s: %w", manifest.Name, err)
	}
	if len(droplets) != 1 {
		return nil, fmt.Errorf("expected one droplet, but found %d", len(droplets))
	}
	return droplets[0], nil
}

func (p *AppPushOperation) findOrg() (*resource.Organization, error) {
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

func (p *AppPushOperation) findSpace(orgGUID string) (*resource.Space, error) {
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
