package operation

import "github.com/cloudfoundry-community/go-cfclient/v3/resource"

type AppHealthCheckType string

const (
	Http    AppHealthCheckType = "http"
	Port    AppHealthCheckType = "port"
	Process AppHealthCheckType = "process"
)

type AppProcessType string

const (
	Web    AppProcessType = "web"
	Worker AppProcessType = "worker"
)

type AppRouteProtocol string

const (
	HTTP1 AppRouteProtocol = "http1"
	HTTP2 AppRouteProtocol = "http2"
	TCP   AppRouteProtocol = "tcp"
)

type Manifest struct {
	Version      string         `yaml:"version,omitempty"`
	Applications []*AppManifest `yaml:"applications"`
}

type AppManifest struct {
	Name               string                `yaml:"name"`
	Path               string                `yaml:"path,omitempty"`
	Buildpacks         []string              `yaml:"buildpacks,omitempty"`
	Docker             *AppManifestDocker    `yaml:"docker,omitempty"`
	Env                map[string]string     `yaml:"env,omitempty"`
	RandomRoute        bool                  `yaml:"random-route,omitempty"`
	NoRoute            bool                  `yaml:"no-route,omitempty"`
	DefaultRoute       bool                  `yaml:"default-route,omitempty"`
	Routes             *AppManifestRoutes    `yaml:"routes,omitempty"`
	Services           *AppManifestServices  `yaml:"services,omitempty"`
	Sidecars           *AppManifestSideCars  `yaml:"sidecars,omitempty"`
	Processes          *AppManifestProcesses `yaml:"processes,omitempty"`
	Stack              string                `yaml:"stack,omitempty"`
	Metadata           *resource.Metadata    `yaml:"metadata,omitempty"`
	AppManifestProcess `yaml:",inline"`
}

type AppManifestProcesses []AppManifestProcess

type AppManifestProcess struct {
	Type                         AppProcessType     `yaml:"type,omitempty"`
	Command                      string             `yaml:"command,omitempty"`
	DiskQuota                    string             `yaml:"disk_quota,omitempty"`
	HealthCheckType              AppHealthCheckType `yaml:"health-check-type,omitempty"`
	HealthCheckHTTPEndpoint      string             `yaml:"health-check-http-endpoint,omitempty"`
	HealthCheckInvocationTimeout uint               `yaml:"health-check-invocation-timeout,omitempty"`
	Instances                    uint               `yaml:"instances,omitempty"`
	LogRateLimitPerSecond        string             `yaml:"log-rate-limit-per-second,omitempty"`
	Memory                       string             `yaml:"memory,omitempty"`
	Timeout                      uint               `yaml:"timeout,omitempty"`
}

type AppManifestDocker struct {
	Image    string `yaml:"image,omitempty"`
	Username string `yaml:"username,omitempty"`
}

type AppManifestServices []AppManifestService

type AppManifestService struct {
	Name        string                 `yaml:"name"`
	BindingName string                 `yaml:"binding_name,omitempty"`
	Parameters  map[string]interface{} `yaml:"parameters,omitempty"`
}

type AppManifestRoutes []AppManifestRoute

type AppManifestRoute struct {
	Route    string           `yaml:"route"`
	Protocol AppRouteProtocol `yaml:"protocol,omitempty"`
}

type AppManifestSideCars []AppManifestSideCar

type AppManifestSideCar struct {
	Name         string   `yaml:"name"`
	ProcessTypes []string `yaml:"process_types,omitempty"`
	Command      string   `yaml:"command,omitempty"`
	Memory       string   `yaml:"memory,omitempty"`
}

func NewManifest(applications ...*AppManifest) *Manifest {
	return &Manifest{
		Version:      "1",
		Applications: applications,
	}
}

func NewAppManifest(appName string) *AppManifest {
	return &AppManifest{
		Name: appName,
		AppManifestProcess: AppManifestProcess{
			HealthCheckType:         "port",
			HealthCheckHTTPEndpoint: "/",
			Instances:               1,
			Memory:                  "256M",
		},
	}
}
