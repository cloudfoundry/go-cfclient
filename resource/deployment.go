package resource

import "time"

type Deployment struct {
	Status          DeploymentStatus   `json:"status"`
	Strategy        string             `json:"strategy"`
	Options         *DeploymentOptions `json:"options,omitempty"`
	Droplet         Relationship       `json:"droplet"`
	PreviousDroplet Relationship       `json:"previous_droplet"`
	NewProcesses    []ProcessReference `json:"new_processes"`
	Revision        DeploymentRevision `json:"revision"`
	Metadata        *Metadata          `json:"metadata"`
	Relationships   AppRelationship    `json:"relationships"`
	Resource        `json:",inline"`
}

type DeploymentCreate struct {
	Relationships AppRelationship     `json:"relationships"`
	Droplet       *Relationship       `json:"droplet,omitempty"`
	Revision      *DeploymentRevision `json:"revision,omitempty"`
	Strategy      string              `json:"strategy,omitempty"`
	Options       *DeploymentOptions  `json:"options,omitempty"`
	Metadata      *Metadata           `json:"metadata,omitempty"`
}

type DeploymentUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type DeploymentList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*Deployment `json:"resources"`
}

type DeploymentRevision struct {
	GUID    string `json:"guid"`
	Version *int   `json:"version,omitempty"`
}

type ProcessReference struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
}

type DeploymentStatus struct {
	Value   string                   `json:"value"`
	Reason  string                   `json:"reason"`
	Details *DeploymentStatusDetails `json:"details,omitempty"`
	Canary  *DeploymentStatusCanary  `json:"canary,omitempty"`
}

type DeploymentStatusDetails struct {
	LastSuccessfulHealthcheck *time.Time `json:"last_successful_healthcheck,omitempty"`
	LastStatusChange          *time.Time `json:"last_status_change,omitempty"`
	Error                     string     `json:"error,omitempty"`
}

type DeploymentStatusCanary struct {
	Steps *DeploymentStatusCanarySteps `json:"steps,omitempty"`
}

type DeploymentStatusCanarySteps struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type DeploymentOptions struct {
	MaxInFlight                  *int                     `json:"max_in_flight,omitempty"`
	WebInstances                 *int                     `json:"web_instances,omitempty"`
	MemoryInMB                   *int                     `json:"memory_in_mb,omitempty"`
	DiskInMB                     *int                     `json:"disk_in_mb,omitempty"`
	LogRateLimitInBytesPerSecond *int                     `json:"log_rate_limit_in_bytes_per_second,omitempty"`
	Canary                       *DeploymentCanaryOptions `json:"canary,omitempty"`
}

type DeploymentCanaryOptions struct {
	Steps []CanaryStep `json:"steps,omitempty"`
}

type CanaryStep struct {
	InstanceWeight int `json:"instance_weight"`
}

func NewDeploymentCreate(appGUID string) *DeploymentCreate {
	return &DeploymentCreate{
		Relationships: AppRelationship{
			App: ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
		},
	}
}
