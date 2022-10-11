package resource

import "time"

type Space struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Relationships map[string]ToOneRelationship `json:"relationships"`
	Links         map[string]Link              `json:"links"`
	Metadata      Metadata                     `json:"metadata"`
}

type CreateSpaceRequest struct {
	Name     string
	OrgGUID  string
	Metadata *Metadata
}

type UpdateSpaceRequest struct {
	Name     string
	Metadata *Metadata
}

type SpaceUsers struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
	Metadata      Metadata                     `json:"metadata,omitempty"`
}

type ListSpacesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Space    `json:"resources,omitempty"`
}

type ListSpaceUsersResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []User     `json:"resources,omitempty"`
}
