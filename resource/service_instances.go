package resource

import (
	"time"
)

type ServiceInstance struct {
	Guid          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Metadata      Metadata                     `json:"metadata"`
	Links         map[string]Link              `json:"links"`
}

type ListServiceInstancesResponse struct {
	Pagination Pagination        `json:"pagination,omitempty"`
	Resources  []ServiceInstance `json:"resources,omitempty"`
}
