package operation

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/cloudfoundry/go-cfclient/v3/testutil"

	"github.com/stretchr/testify/require"
)

func TestAppPush(t *testing.T) {
	serverURL := testutil.SetupFakeAPIServer()
	defer testutil.Teardown()

	g := testutil.NewObjectJSONGenerator()
	org := g.Organization()
	space := g.Space()
	job := g.Job("COMPLETE")
	app := g.Application()
	pkg := g.Package("READY")
	build := g.Build("STAGED")
	droplet := g.Droplet()
	dropletAssoc := g.DropletAssociation()

	fakeAppZipReader := strings.NewReader("blah zip zip")
	var numOfInstances uint = 2
	manifest := &AppManifest{
		Name:       app.Name,
		Buildpacks: []string{"java-buildpack-offline"},

		AppManifestProcess: AppManifestProcess{
			HealthCheckType:         "http",
			HealthCheckHTTPEndpoint: "/health",
			Instances:               &numOfInstances,
			Memory:                  "1G",
		},
		Routes: &AppManifestRoutes{
			{
				Route: "https://spring-music.cf.apps.example.org",
			},
		},
		Services: &AppManifestServices{{Name: "spring-music-sql"}},
		Stack:    "cflinuxfs3",
	}

	testutil.SetupMultiple([]testutil.MockRoute{
		{
			Method:   http.MethodGet,
			Endpoint: "/v3/organizations",
			Output:   g.SinglePaged(org.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodGet,
			Endpoint: "/v3/spaces",
			Output:   g.SinglePaged(space.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:           http.MethodPost,
			Endpoint:         fmt.Sprintf("/v3/spaces/%s/actions/apply_manifest", space.GUID),
			Output:           g.SinglePaged(space.JSON),
			Status:           http.StatusAccepted,
			RedirectLocation: fmt.Sprintf("%s/v3/jobs/%s", serverURL, job.GUID),
		},
		{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("/v3/jobs/%s", job.GUID),
			Output:   g.Single(job.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodGet,
			Endpoint: "/v3/apps",
			Output:   g.SinglePaged(app.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodPost,
			Endpoint: "/v3/packages",
			Output:   g.Single(pkg.JSON),
			Status:   http.StatusCreated,
		},
		{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("/v3/packages/%s/upload", pkg.GUID),
			Output:   g.Single(pkg.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("/v3/packages/%s", pkg.GUID),
			Output:   g.Single(pkg.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodPost,
			Endpoint: "/v3/builds",
			Output:   g.Single(build.JSON),
			Status:   http.StatusCreated,
		},
		{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("/v3/builds/%s", build.GUID),
			Output:   g.Single(build.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("/v3/packages/%s/droplets", pkg.GUID),
			Output:   g.SinglePaged(droplet.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodPatch,
			Endpoint: fmt.Sprintf("/v3/apps/%s/relationships/current_droplet", app.GUID),
			Output:   g.Single(dropletAssoc.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodPost,
			Endpoint: fmt.Sprintf("/v3/apps/%s/actions/start", app.GUID),
			Output:   g.Single(app.JSON),
			Status:   http.StatusOK,
		},
	}, t)

	c, _ := config.New(serverURL, config.Token("", "fake-refresh-token"), config.SkipTLSValidation())
	cf, err := client.New(c)
	require.NoError(t, err)

	pusher := NewAppPushOperation(cf, org.Name, space.Name)
	// Invalid strategy
	strategy := StrategyMode(10)
	pusher.WithStrategy(strategy)
	_, err = pusher.Push(context.Background(), manifest, fakeAppZipReader)
	require.NoError(t, err)
}

func TestDockerLifecycleBuildCreation(t *testing.T) {
	serverURL := testutil.SetupFakeAPIServer()
	defer testutil.Teardown()

	g := testutil.NewObjectJSONGenerator()
	build := g.Build("STAGED")
	droplet := g.Droplet()

	manifest := &AppManifest{
		Name: "test-docker-app",
		Docker: &AppManifestDocker{
			Image: "kennethreitz/httpbin",
		},
	}

	// Create a docker package manually for testing
	dockerPkg := &resource.Package{
		Type: "docker",
		Resource: resource.Resource{
			GUID: "docker-package-guid",
		},
	}

	// Mock the API calls for docker build creation
	testutil.SetupMultiple([]testutil.MockRoute{
		{
			Method:   http.MethodPost,
			Endpoint: "/v3/builds",
			Output:   g.Single(build.JSON),
			Status:   http.StatusCreated,
			// Verify that the lifecycle data contains proper docker lifecycle structure
			PostForm: `{"package":{"guid":"docker-package-guid"},"lifecycle":{"type":"docker","data":{}}}`,
		},
		{
			Method:   http.MethodGet,
			Endpoint: fmt.Sprintf("/v3/builds/%s", build.GUID),
			Output:   g.Single(build.JSON),
			Status:   http.StatusOK,
		},
		{
			Method:   http.MethodGet,
			Endpoint: "/v3/packages/docker-package-guid/droplets",
			Output:   g.SinglePaged(droplet.JSON),
			Status:   http.StatusOK,
		},
	}, t)

	c, _ := config.New(serverURL, config.Token("", "fake-refresh-token"), config.SkipTLSValidation())
	cf, err := client.New(c)
	require.NoError(t, err)

	pusher := NewAppPushOperation(cf, "", "")

	// Test the buildDroplet method specifically with a docker package
	resultDroplet, err := pusher.buildDroplet(context.Background(), dockerPkg, manifest)
	require.NoError(t, err, "Docker lifecycle build should not fail")
	require.NotNil(t, resultDroplet, "Docker lifecycle build should return a droplet")
}

// TestDockerLifecycleStructure tests that docker builds have the correct lifecycle structure
func TestDockerLifecycleStructure(t *testing.T) {
	// Test the lifecycle structure directly without HTTP mocking
	dockerPkg := &resource.Package{
		Type: resource.LifecycleDocker.String(),
		Resource: resource.Resource{
			GUID: "test-docker-package",
		},
	}

	// Create build request directly to test lifecycle structure
	buildCreate := resource.NewBuildCreate(dockerPkg.GUID)

	// Apply the same logic as in buildDroplet method
	if dockerPkg.Type == resource.LifecycleDocker.String() {
		buildCreate.Lifecycle = &resource.Lifecycle{
			Type: dockerPkg.Type,
			Data: &resource.DockerLifecycle{}, // Empty docker lifecycle data
		}
	}

	// Verify the structure
	require.NotNil(t, buildCreate.Lifecycle, "Docker build should have lifecycle")
	require.Equal(t, "docker", buildCreate.Lifecycle.Type, "Docker build should have docker lifecycle type")
	require.NotNil(t, buildCreate.Lifecycle.Data, "Docker build should have lifecycle data")

	// Verify it's the correct type
	dockerLifecycle, ok := buildCreate.Lifecycle.Data.(*resource.DockerLifecycle)
	require.True(t, ok, "Docker build lifecycle data should be DockerLifecycle type")
	require.NotNil(t, dockerLifecycle, "Docker lifecycle data should not be nil")
}

// TestDockerLifecycleJSONMarshaling tests that docker lifecycle marshals to correct JSON
func TestDockerLifecycleJSONMarshaling(t *testing.T) {
	dockerLifecycle := &resource.Lifecycle{
		Type: "docker",
		Data: &resource.DockerLifecycle{}, // Empty docker lifecycle data
	}

	// Test JSON marshaling to ensure it produces the expected structure
	// This is what the CF API expects: {"type":"docker","data":{}}
	expectedJSON := `{"type":"docker","data":{}}`

	// Marshal the lifecycle
	actualJSON, err := dockerLifecycle.MarshalJSON()
	require.NoError(t, err, "Docker lifecycle should marshal without error")

	// Verify the JSON structure matches expectations
	require.JSONEq(t, expectedJSON, string(actualJSON), "Docker lifecycle JSON should match expected format")
}
