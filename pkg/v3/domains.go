package v3

import (
	"time"
)

type DomainRelationships struct {
	Organization        ToOneRelationship   `json:"organization"`
	SharedOrganizations ToManyRelationships `json:"shared_organizations"`
}

type Domain struct {
	GUID          string              `json:"guid"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Name          string              `json:"name"`
	Internal      bool                `json:"internal"`
	Metadata      Metadata            `json:"metadata"`
	Relationships DomainRelationships `json:"relationships"`
	Links         map[string]Link     `json:"links"`
}

type ListDomainsResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Domain   `json:"resources,omitempty"`
}
