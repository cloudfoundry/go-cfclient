package operation

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type StrategyMode int

const (
	StrategyNone StrategyMode = iota
	StrategyBlueGreen
	StrategyRolling
)

// AppPushOperation can be used to push buildpack apps
type AppPushOperation struct {
	orgName   string
	spaceName string
	client    *client.Client
	strategy  StrategyMode
}

// NewAppPushOperation creates a new AppPushOperation
func NewAppPushOperation(client *client.Client, orgName, spaceName string) *AppPushOperation {
	apo := AppPushOperation{
		orgName:   orgName,
		spaceName: spaceName,
		client:    client,
	}
	apo.strategy = StrategyNone
	return &apo
}
func (p *AppPushOperation) WithStrategy(s StrategyMode) {
	switch s {
	case StrategyBlueGreen, StrategyRolling:
		p.strategy = s
	default:
		p.strategy = StrategyNone
	}
}

// Push creates or updates an application using the specified manifest and zipped source files
func (p *AppPushOperation) Push(ctx context.Context, appManifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	org, err := p.findOrg(ctx)
	if err != nil {
		return nil, err
	}
	space, err := p.findSpace(ctx, org.GUID)
	if err != nil {
		return nil, err
	}
	return p.pushWithStrategyApp(ctx, space, appManifest, zipFile)
}
func (p *AppPushOperation) pushWithStrategyApp(ctx context.Context, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	switch p.strategy {
	case StrategyBlueGreen:
		return p.pushBlueGreenApp(ctx, space, manifest, zipFile)
	case StrategyRolling:
		return p.pushRollingApp(ctx, space, manifest, zipFile)
	default:
		return p.pushApp(ctx, space, manifest, zipFile)
	}
}

func (p *AppPushOperation) pushBlueGreenApp(ctx context.Context, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	originalApp, err := p.findApp(ctx, manifest.Name, space)
	if err != nil && err != client.ErrExactlyOneResultNotReturned {
		return nil, err
	}
	if err == client.ErrExactlyOneResultNotReturned || originalApp.State != "STARTED" {
		return p.pushApp(ctx, space, manifest, zipFile)
	}

	tempAppName := originalApp.Name + "-venerable"

	// Check if temporary app name already exists if yes gracefully delete the app and continue
	tempApp, err := p.findApp(ctx, tempAppName, space)
	if err == nil {
		err = p.gracefulDeletion(ctx, tempApp)
		if err != nil {
			return nil, err
		}
	}

	// Update the existing app's name
	_, err = p.client.Applications.Update(ctx, originalApp.GUID, &resource.AppUpdate{
		Name: tempAppName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update app name failed with: %s", err.Error())
	}

	// Apply the manifest
	newApp, err := p.pushApp(ctx, space, manifest, zipFile)
	if err != nil {
		// If push fails change back original app name
		_, err = p.client.Applications.Update(ctx, originalApp.GUID, &resource.AppUpdate{
			Name: originalApp.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update app name back to original name: failed with %s", err.Error())
		}
		return nil, fmt.Errorf("blue green deployment failed with: %s", err.Error())
	}
	if newApp.State == "STARTED" {
		err = p.gracefulDeletion(ctx, originalApp)
		return newApp, err
	}
	return newApp, fmt.Errorf("failed to verify application start: %s", err.Error())
}

func (p *AppPushOperation) pushRollingApp(ctx context.Context, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	originalApp, err := p.findApp(ctx, manifest.Name, space)
	if err != nil && err != client.ErrExactlyOneResultNotReturned {
		return nil, err
	}
	if err == client.ErrExactlyOneResultNotReturned || originalApp.State != "STARTED" {
		return p.pushApp(ctx, space, manifest, zipFile)
	}
	// Get the fallback revision in case of rollback
	fallbackRevision, _ := p.client.Revisions.SingleForAppDeployed(ctx, originalApp.GUID, nil)

	err = p.applySpaceManifest(ctx, space, manifest)
	if err != nil {
		return nil, err
	}

	var pkg *resource.Package
	if manifest.Docker != nil {
		pkg, err = p.uploadDockerPackage(ctx, originalApp, manifest.Docker)
	} else {
		pkg, err = p.uploadBitsPackage(ctx, originalApp, zipFile)
	}
	if err != nil {
		return nil, err
	}

	droplet, err := p.buildDroplet(ctx, pkg, manifest)
	if err != nil {
		return nil, err
	}

	deployment, err := p.createNewDeployment(ctx, originalApp, droplet)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy with: %s", err.Error())
	}
	// In case application crashed due to new deployment, deployment will be stuck with value "ACTIVE" and reason "DEPLOYING"
	// This will be considered as deployment failed after timeout
	depPollErr := p.waitForDeployment(ctx, deployment.GUID, *manifest.Instances)

	// Check the app state if app not started or deployment failed rollback the deployment
	originalApp, err = p.findApp(ctx, manifest.Name, space)
	if err != nil {
		return nil, fmt.Errorf("failed to verify application status with: %s", err.Error())
	}
	if originalApp.State != "STARTED" || depPollErr != nil {
		rollBackDeployment, rollBackErr := p.rollBackDeployment(ctx, originalApp, fallbackRevision)
		if rollBackErr != nil {
			return nil, fmt.Errorf("failed to confirm rollback deployment with: %s", rollBackErr.Error())
		}
		depRollPollErr := p.waitForDeployment(ctx, rollBackDeployment.GUID, *manifest.Instances)
		if depRollPollErr != nil {
			return nil, fmt.Errorf("failed to deploy with: %s \nfailed to confirm roll back to last deployment with: %s", depPollErr.Error(), depRollPollErr.Error())
		}
		return nil, fmt.Errorf("failed to deploy with: %s \nrolled back to last deployment", depPollErr.Error())
	}

	return originalApp, nil
}

// Poll for deployment status and wait for the deployment to be in the final state
// Timeout is calculated based on the number of instances
func (p *AppPushOperation) waitForDeployment(ctx context.Context, deploymentGUID string, instances uint) error {
	// If instances is not set default to 1
	if instances == 0 {
		instances = 1
	}
	pollOptions := client.NewPollingOptions()
	pollOptions.Timeout = time.Duration(instances) * time.Minute

	depPollErr := client.PollForStateOrTimeout(func() (string, error) {
		deployment, err := p.client.Deployments.Get(ctx, deploymentGUID)
		if err != nil {
			return "", err
		}
		return deployment.Status.Value, nil
	}, "FINALIZED", pollOptions)
	return depPollErr
}

func (p *AppPushOperation) createNewDeployment(ctx context.Context, originalApp *resource.App, droplet *resource.Droplet) (*resource.Deployment, error) {
	return p.client.Deployments.Create(ctx, &resource.DeploymentCreate{
		Relationships: resource.AppRelationship{
			App: resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: originalApp.GUID,
				},
			},
		},
		Droplet: &resource.Relationship{
			GUID: droplet.GUID,
		},
	})
}

func (p *AppPushOperation) rollBackDeployment(ctx context.Context, originalApp *resource.App, fallbackRevision *resource.Revision) (*resource.Deployment, error) {
	return p.client.Deployments.Create(ctx, &resource.DeploymentCreate{
		Relationships: resource.AppRelationship{
			App: resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: originalApp.GUID,
				},
			},
		},
		Revision: &resource.DeploymentRevision{
			GUID: fallbackRevision.GUID,
		},
	})
}

// Stop the application and delete it
// https://github.com/cloudfoundry/cloud_controller_ng/issues/1017
func (p *AppPushOperation) gracefulDeletion(ctx context.Context, app *resource.App) error {
	app, err := p.client.Applications.Stop(ctx, app.GUID)
	if err != nil {
		return fmt.Errorf("failed to stop the application with: %s", err.Error())
	}
	jobId, err := p.client.Applications.Delete(ctx, app.GUID)
	if err != nil {
		return err
	}
	err = p.client.Jobs.PollComplete(ctx, jobId, &client.PollingOptions{
		Timeout:       20 * time.Minute,
		CheckInterval: time.Second * 5,
		FailedState:   string(resource.JobStateFailed),
	})
	if err != nil {
		return err
	}
	return nil
}

// pushApp pushes an application
//
// After an application is created and packages are uploaded, a droplet must be created via a build in order for
// an application to be deployed or tasks to be run. The current droplet must be assigned to an application before
// it may be started. When tasks are created, they either use a specific droplet guid, or use the current droplet
// assigned to an application.
func (p *AppPushOperation) pushApp(ctx context.Context, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	err := p.applySpaceManifest(ctx, space, manifest)
	if err != nil {
		return nil, err
	}

	app, err := p.findApp(ctx, manifest.Name, space)
	if err != nil {
		return nil, err
	}

	var pkg *resource.Package
	if app.Lifecycle.Type == resource.LifecycleDocker.String() {
		pkg, err = p.uploadDockerPackage(ctx, app, manifest.Docker)
	} else {
		pkg, err = p.uploadBitsPackage(ctx, app, zipFile)
	}
	if err != nil {
		return nil, err
	}

	droplet, err := p.buildDroplet(ctx, pkg, manifest)
	if err != nil {
		return nil, err
	}

	_, err = p.client.Droplets.SetCurrentAssociationForApp(ctx, app.GUID, droplet.GUID)
	if err != nil {
		return nil, err
	}

	return p.client.Applications.Start(ctx, app.GUID)
}

func (p *AppPushOperation) applySpaceManifest(ctx context.Context, space *resource.Space, manifest *AppManifest) error {
	// wrap it in a manifest that has an applications array as required by the API
	multiAppsManifest := &Manifest{
		Applications: []*AppManifest{manifest},
	}
	manifestBytes, err := yaml.Marshal(&multiAppsManifest)
	if err != nil {
		return fmt.Errorf("error marshalling application manifest: %w", err)
	}

	jobGUID, err := p.client.Manifests.ApplyManifest(ctx, space.GUID, string(manifestBytes))
	if err != nil {
		return fmt.Errorf("error applying application manifest to space %s: %w", space.Name, err)
	}
	err = p.client.Jobs.PollComplete(ctx, jobGUID, nil)
	if err != nil {
		return fmt.Errorf("error waiting for application manifest to finish applying to space %s: %w", space.Name, err)
	}
	return nil
}

func (p *AppPushOperation) findApp(ctx context.Context, appName string, space *resource.Space) (*resource.App, error) {
	appOpts := client.NewAppListOptions()
	appOpts.Names.EqualTo(appName)
	appOpts.SpaceGUIDs.EqualTo(space.GUID)
	app, err := p.client.Applications.Single(ctx, appOpts)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *AppPushOperation) uploadDockerPackage(ctx context.Context, app *resource.App, docker *AppManifestDocker) (*resource.Package, error) {
	newPkg := resource.NewDockerPackageCreate(app.GUID, docker.Image, docker.Username, os.Getenv("CF_DOCKER_PASSWORD"))
	pkg, err := p.client.Packages.Create(ctx, newPkg)
	if err != nil {
		return nil, fmt.Errorf("error creating docker package for app %s: %w", app.Name, err)
	}
	return pkg, nil
}

func (p *AppPushOperation) uploadBitsPackage(ctx context.Context, app *resource.App, zipFile io.Reader) (*resource.Package, error) {
	newPkg := resource.NewPackageCreate(app.GUID)
	pkg, err := p.client.Packages.Create(ctx, newPkg)
	if err != nil {
		return nil, fmt.Errorf("error creating package bits for app %s: %w", app.Name, err)
	}
	_, err = p.client.Packages.Upload(ctx, pkg.GUID, zipFile)
	if err != nil {
		return nil, fmt.Errorf("error uploading package bits for app %s: %w", app.Name, err)
	}
	err = p.client.Packages.PollReady(ctx, pkg.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for package to process for app %s: %w", app.Name, err)
	}
	return pkg, nil
}

func (p *AppPushOperation) buildDroplet(ctx context.Context, pkg *resource.Package, manifest *AppManifest) (*resource.Droplet, error) {
	newBuild := resource.NewBuildCreate(pkg.GUID)
	if pkg.Type == resource.LifecycleDocker.String() {
		newBuild.Lifecycle = &resource.Lifecycle{Type: pkg.Type}
	} else {
		newBuild.Lifecycle = &resource.Lifecycle{
			Type: resource.LifecycleBuildpack.String(),
			BuildpackData: resource.BuildpackLifecycle{
				Buildpacks: manifest.Buildpacks,
				Stack:      manifest.Stack,
			},
		}
	}
	build, err := p.client.Builds.Create(ctx, newBuild)
	if err != nil {
		return nil, fmt.Errorf("error creating build from package for app %s: %w", manifest.Name, err)
	}
	err = p.client.Builds.PollStaged(ctx, build.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for app %s package to build: %w", manifest.Name, err)
	}

	opts := client.NewDropletPackageListOptions()
	opts.States.EqualTo(resource.DropletStateStaged.String())
	droplet, err := p.client.Droplets.SingleForPackage(ctx, pkg.GUID, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding droplet for app %s: %w", manifest.Name, err)
	}
	return droplet, nil
}

func (p *AppPushOperation) findOrg(ctx context.Context) (*resource.Organization, error) {
	opts := client.NewOrganizationListOptions()
	opts.Names.EqualTo(p.orgName)
	org, err := p.client.Organizations.Single(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find org %s: %w", p.orgName, err)
	}
	return org, nil
}

func (p *AppPushOperation) findSpace(ctx context.Context, orgGUID string) (*resource.Space, error) {
	opts := client.NewSpaceListOptions()
	opts.Names.EqualTo(p.spaceName)
	opts.OrganizationGUIDs.EqualTo(orgGUID)
	space, err := p.client.Spaces.Single(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find space %s: %w", p.spaceName, err)
	}
	return space, nil
}
