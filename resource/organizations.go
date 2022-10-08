package resource

type Organization struct {
	Name          string                       `json:"name,omitempty"`
	GUID          string                       `json:"guid,omitempty"`
	Suspended     *bool                        `json:"suspended,omitempty"`
	CreatedAt     string                       `json:"created_at,omitempty"`
	UpdatedAt     string                       `json:"updated_at,omitempty"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
	Metadata      Metadata                     `json:"metadata,omitempty"`
}

type CreateOrganizationRequest struct {
	Name      string
	Suspended *bool `json:"suspended,omitempty"`
	Metadata  *Metadata
}

type UpdateOrganizationRequest struct {
	Name      string
	Suspended *bool `json:"suspended,omitempty"`
	Metadata  *Metadata
}

type ListOrganizationsResponse struct {
	Pagination Pagination      `json:"pagination,omitempty"`
	Resources  []*Organization `json:"resources,omitempty"`
}
