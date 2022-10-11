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

type OrganizationCreate struct {
	Name      string    `json:"name"`
	Suspended *bool     `json:"suspended,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type OrganizationUpdate struct {
	Name      string    `json:"name,omitempty"`
	Suspended *bool     `json:"suspended,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type OrganizationList struct {
	Pagination Pagination      `json:"pagination,omitempty"`
	Resources  []*Organization `json:"resources,omitempty"`
}
