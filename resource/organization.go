package resource

type Organization struct {
	Name          string            `json:"name"`
	Suspended     *bool             `json:"suspended,omitempty"`
	Relationships QuotaRelationship `json:"relationships,omitempty"`
	Metadata      *Metadata         `json:"metadata,omitempty"`
	Resource      `json:",inline"`
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

type OrganizationUsageSummary struct {
	UsageSummary UsageSummary    `json:"usage_summary"`
	Links        map[string]Link `json:"links,omitempty"`
}

type OrganizationList struct {
	Pagination Pagination      `json:"pagination,omitempty"`
	Resources  []*Organization `json:"resources,omitempty"`
}

type UsageSummary struct {
	StartedInstances int `json:"started_instances"`
	MemoryInMb       int `json:"memory_in_mb"`
}

type QuotaRelationship struct {
	Quota ToOneRelationship `json:"quota"`
}

func NewOrganizationCreate(name string) *OrganizationCreate {
	return &OrganizationCreate{
		Name: name,
	}
}
