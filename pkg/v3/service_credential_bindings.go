package v3

import (
	"time"
)

// ServiceCredentialBindings implements the service credential binding object. a credential binding can be a binding between apps and a service instance or a service key
type ServiceCredentialBindings struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Type          string                       `json:"type"`
	LastOperation LastOperation                `json:"last_operation"`
	Metadata      Metadata                     `json:"metadata"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links"`
}

type LastOperation struct {
	Type        string `json:"type"`
	State       string `json:"state"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type ListServiceCredentialBindingsResponse struct {
	Pagination Pagination                  `json:"pagination,omitempty"`
	Resources  []ServiceCredentialBindings `json:"resources,omitempty"`
}
