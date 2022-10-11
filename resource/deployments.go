package resource

import "time"

type CreateDeploymentOptionalParameters struct {
	Droplet  *Relationship       `json:"droplet,omitempty"`
	Revision *DeploymentRevision `json:"revision,omitempty"`
	Strategy *string             `json:"strategy,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type CreateDeploymentRequest struct {
	*CreateDeploymentOptionalParameters `json:",inline"`
	Relationships                       struct {
		App ToOneRelationship `json:"app"`
	} `json:"relationships"`
}

type DeploymentRevision struct {
	GUID    string `json:"guid"`
	Version int    `json:"version"`
}

type ProcessReference struct {
	GUID string `json:"guid"`
	Type string `type:"type"`
}

type DeploymentStatus struct {
	Value   string            `json:"value"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details"`
}

type Deployment struct {
	GUID            string                       `json:"guid"`
	State           string                       `json:"state"`
	Status          DeploymentStatus             `json:"status"`
	Strategy        string                       `json:"strategy"`
	Droplet         Relationship                 `json:"droplet"`
	PreviousDroplet Relationship                 `json:"previous_droplet"`
	NewProcesses    []ProcessReference           `json:"new_processes"`
	Revision        DeploymentRevision           `json:"revision"`
	CreatedAt       time.Time                    `json:"created_at"`
	UpdatedAt       time.Time                    `json:"updated_at"`
	Links           map[string]Link              `json:"links,omitempty"`
	Metadata        Metadata                     `json:"metadata,omitempty"`
	Relationships   map[string]ToOneRelationship `json:"relationships,omitempty"`
}
