package resource

type Space struct {
	Name          string                       `json:"name,omitempty"`
	GUID          string                       `json:"guid,omitempty"`
	CreatedAt     string                       `json:"created_at,omitempty"`
	UpdatedAt     string                       `json:"updated_at,omitempty"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
	Metadata      Metadata                     `json:"metadata,omitempty"`
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
	Name          string                       `json:"name,omitempty"`
	GUID          string                       `json:"guid,omitempty"`
	CreatedAt     string                       `json:"created_at,omitempty"`
	UpdatedAt     string                       `json:"updated_at,omitempty"`
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
