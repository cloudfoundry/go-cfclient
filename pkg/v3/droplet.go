package v3

// Droplet is the result of staging an application package.
// There are two types (lifecycles) of droplets: buildpack and
// docker. In the case of buildpacks, the droplet contains the
// bits produced by the buildpack.
type Droplet struct {
	State             string            `json:"state,omitempty"`
	Error             string            `json:"error,omitempty"`
	Lifecycle         Lifecycle         `json:"lifecycle,omitempty"`
	GUID              string            `json:"guid,omitempty"`
	CreatedAt         string            `json:"created_at,omitempty"`
	UpdatedAt         string            `json:"updated_at,omitempty"`
	Links             map[string]Link   `json:"links,omitempty"`
	ExecutionMetadata string            `json:"execution_metadata,omitempty"`
	ProcessTypes      map[string]string `json:"process_types,omitempty"`
	Metadata          Metadata          `json:"metadata,omitempty"`

	// Only specified when the droplet is using the Docker lifecycle.
	Image string `json:"image,omitempty"`

	// The following fields are specified when the droplet is using
	// the buildpack lifecycle.
	Checksum struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"checksum,omitempty"`
	Stack      string              `json:"stack,omitempty"`
	Buildpacks []DetectedBuildpack `json:"buildpacks,omitempty"`
}

type DetectedBuildpack struct {
	Name          string `json:"name,omitempty"`           // system buildpack name
	BuildpackName string `json:"buildpack_name,omitempty"` // name reported by the buildpack
	DetectOutput  string `json:"detect_output,omitempty"`  // output during detect process
	Version       string `json:"version,omitempty"`
}

type CurrentDropletResponse struct {
	Data  Relationship    `json:"data,omitempty"`
	Links map[string]Link `json:"links,omitempty"`
}
