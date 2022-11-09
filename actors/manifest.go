package actors

type Manifest struct {
	Applications []AppManifest `yaml:"applications"`
}

type AppManifest struct {
	Name       string   `yaml:"name"`
	Buildpacks []string `yaml:"buildpacks"`
	Command    string   `yaml:"command"`
	DiskQuota  string   `yaml:"disk_quota"`
	Docker     struct {
		Image    string `yaml:"image,omitempty"`
		Username string `yaml:"username,omitempty"`
	} `yaml:"docker,omitempty"`
	Env                     map[string]string `yaml:"env"`
	HealthCheckType         string            `yaml:"health-check-type"`
	HealthCheckHTTPEndpoint string            `yaml:"health-check-http-endpoint,omitempty"`
	Instances               int               `yaml:"instances"`
	LogRateLimit            string            `yaml:"log-rate-limit"`
	Memory                  string            `yaml:"memory"`
	NoRoute                 bool              `yaml:"no-route,omitempty"`
	Routes                  []struct {
		Route string `yaml:"route,omitempty"`
	} `yaml:"routes,omitempty"`
	Services []string `yaml:"services"`
	Stack    string   `yaml:"stack"`
	Timeout  int      `yaml:"timeout,omitempty"`
}
