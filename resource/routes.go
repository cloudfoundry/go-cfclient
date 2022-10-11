package resource

import "time"

type Route struct {
	Guid          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Host          string                       `json:"host"`
	Path          string                       `json:"path"`
	Url           string                       `json:"url"`
	Metadata      Metadata                     `json:"metadata"`
	Relationships map[string]ToOneRelationship `json:"relationships"`
	Links         map[string]Link              `json:"links"`
}

type ListRouteResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Route    `json:"resources,omitempty"`
}

type CreateRouteOptionalParameters struct {
	Host     string   `json:"host,omitempty"`
	Path     string   `json:"path,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

type RouteRelationships struct {
	Space  ToOneRelationship `json:"space"`
	Domain ToOneRelationship `json:"domain"`
}

type CreateRouteRequest struct {
	Relationships RouteRelationships `json:"relationships"`
	*CreateRouteOptionalParameters
}
