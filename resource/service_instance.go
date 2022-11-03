package resource

import (
	"encoding/json"
	"time"
)

type ServiceInstance struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Tags          []string                     `json:"tags"` // Used by apps to identify service instances; they are shown in the app VCAP_SERVICES env
	Type          string                       `json:"type"` // Either managed or user-provided
	LastOperation LastOperation                `json:"last_operation"`
	Relationships ServiceInstanceRelationships `json:"relationships"`
	Metadata      *Metadata                    `json:"metadata"`
	Links         map[string]Link              `json:"links,omitempty"`

	// Information about the version of this service instance; only shown when type is managed
	MaintenanceInfo *ServiceInstanceMaintenanceInfo `json:"maintenance_info,omitempty"`

	// Whether an upgrade of this service instance is available on the current Service Plan
	// Details are available in the maintenance_info object; Only shown when type is managed
	UpgradeAvailable *bool `json:"upgrade_available,omitempty"`

	// The URL to the service instance dashboard (or null if there is none); only shown when type is managed
	DashboardURL *string `json:"dashboard_url,omitempty"`
}

type ServiceInstanceCreate struct {
	Type          string                       `json:"type"` // Either managed or user-provided
	Name          string                       `json:"name"`
	Relationships ServiceInstanceRelationships `json:"relationships"`
	Metadata      *Metadata                    `json:"metadata,omitempty"`
	Parameters    *json.RawMessage             `json:"parameters,omitempty"` // A JSON object that is passed to the service broker
	Tags          []string                     `json:"tags,omitempty"`
}

type ServiceInstanceList struct {
	Pagination Pagination         `json:"pagination"`
	Resources  []*ServiceInstance `json:"resources"`
}

type ServiceInstanceMaintenanceInfo struct {
	// The current semantic version of this service instance
	// Comparing this version with the version of the Service Plan can be used to determine
	// whether this service instance is up-to-date with the Service Plan
	Version string `json:"version"`

	// A textual explanation associated with this version
	Description string `json:"description,omitempty"`
}

type ServiceInstanceRelationships struct {
	// The service plan the service instance relates to; only shown when type is managed
	ServicePlan *ToOneRelationship `json:"service_plan,omitempty"`

	// The space the service instance is contained in
	Space ToOneRelationship `json:"space"`
}

func NewServiceInstanceCreateManaged(name, spaceGUID, servicePlanGUID string) *ServiceInstanceCreate {
	return &ServiceInstanceCreate{
		Type: "managed",
		Name: name,
		Relationships: ServiceInstanceRelationships{
			ServicePlan: &ToOneRelationship{
				Data: &Relationship{
					GUID: servicePlanGUID,
				},
			},
			Space: ToOneRelationship{
				Data: &Relationship{
					GUID: spaceGUID,
				},
			},
		},
	}
}

func NewServiceInstanceCreateUserProvided(name, spaceGUID string) *ServiceInstanceCreate {
	return &ServiceInstanceCreate{
		Type: "user-provided",
		Name: name,
		Relationships: ServiceInstanceRelationships{
			Space: ToOneRelationship{
				Data: &Relationship{
					GUID: spaceGUID,
				},
			},
		},
	}
}
