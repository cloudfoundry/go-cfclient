package resource

import "time"

// Droplet is the result of staging an application package.
// There are two types (lifecycles) of droplets: buildpack and
// docker. In the case of buildpacks, the droplet contains the
// bits produced by the buildpack.
type Droplet struct {
	GUID              string            `json:"guid"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	State             string            `json:"state"`
	Error             *string           `json:"error"`
	Lifecycle         Lifecycle         `json:"lifecycle"`
	Links             map[string]Link   `json:"links"`
	ExecutionMetadata string            `json:"execution_metadata"`
	ProcessTypes      map[string]string `json:"process_types"`
	Metadata          Metadata          `json:"metadata"`

	// Link to the associated app, if any
	Relationships DropletRelationships `json:"relationships"`

	// Only specified when the droplet is using the Docker lifecycle.
	Image *string `json:"image"`

	// The following fields are specified when the droplet is using
	// the buildpack lifecycle.
	Checksum struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"checksum"`
	Stack      string              `json:"stack"`
	Buildpacks []DetectedBuildpack `json:"buildpacks"`
}

type DetectedBuildpack struct {
	Name          string `json:"name"`           // system buildpack name
	BuildpackName string `json:"buildpack_name"` // name reported by the buildpack
	DetectOutput  string `json:"detect_output"`  // output during detect process
	Version       string `json:"version"`
}

type CurrentDropletResponse struct {
	Data  Relationship    `json:"data"`
	Links map[string]Link `json:"links"`
}

type DropletRelationships struct {
	App ToOneRelationship `json:"app"`
}

type DropletCreate struct {
	Relationships DropletRelationships `json:"relationships"`
	ProcessTypes  map[string]string    `json:"process_types"`
}

func NewDropletCreate(appGUID string) *DropletCreate {
	return &DropletCreate{
		Relationships: DropletRelationships{
			App: ToOneRelationship{
				Data: Relationship{
					GUID: appGUID,
				},
			},
		},
		ProcessTypes: make(map[string]string),
	}
}

type DropletList struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []*Droplet `json:"resources,omitempty"`
}

type DropletUpdate struct {
	Metadata Metadata `json:"metadata,omitempty"`
	Image    string   `json:"image"`
}

type DropletCopy struct {
	// Link to the associated app
	Relationships DropletRelationships `json:"relationships"`
}

func NewDropletCopy(appGUID string) *DropletCopy {
	return &DropletCopy{
		Relationships: DropletRelationships{
			App: ToOneRelationship{
				Data: Relationship{
					GUID: appGUID,
				},
			},
		},
	}
}
