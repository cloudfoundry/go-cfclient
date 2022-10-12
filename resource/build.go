package resource

import "time"

type Build struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     string    `json:"state"`
	Error     *string   `json:"error,omitempty"`

	StagingMemoryInMB                 int `json:"staging_memory_in_mb"`
	StagingDiskInMB                   int `json:"staging_disk_in_mb"`
	StagingLogRateLimitBytesPerSecond int `json:"staging_log_rate_limit_bytes_per_second"`

	Lifecycle     Lifecycle                    `json:"lifecycle"`
	Package       Relationship                 `json:"package"`
	Droplet       *Relationship                `json:"droplet,omitempty"`
	CreatedBy     CreatedBy                    `json:"created_by"`
	Links         map[string]Link              `json:"links"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Metadata      Metadata                     `json:"metadata"`
}

// The 3 lifecycle states
const (
	BuildStateStaging = "STAGING"
	BuildStateStaged  = "STAGED"
	BuildStateFailed  = "FAILED"
)

type CreatedBy struct {
	GUID  string `json:"guid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BuildCreate struct {
	Package                           Relationship `json:"package"`
	Lifecycle                         *Lifecycle   `json:"lifecycle,omitempty"`
	StagingMemoryInMB                 int          `json:"staging_memory_in_mb,omitempty"`
	StagingDiskInMB                   int          `json:"staging_disk_in_mb,omitempty"`
	StagingLogRateLimitBytesPerSecond int          `json:"staging_log_rate_limit_bytes_per_second,omitempty"`
	Metadata                          *Metadata    `json:"metadata,omitempty"`
}

func NewBuildCreate(packageGUID string) *BuildCreate {
	return &BuildCreate{
		Package: Relationship{
			GUID: packageGUID,
		},
	}
}

type BuildUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

func NewBuildUpdate() *BuildUpdate {
	return &BuildUpdate{
		Metadata: &Metadata{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
	}
}

type BuildList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Build   `json:"resources"`
}
