package resource

type Organization struct {
	Name          string            `json:"name"`
	Suspended     bool              `json:"suspended"`
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
	Routes           int `json:"routes"`
	ServiceInstances int `json:"service_instances"`
	ReservedPorts    int `json:"reserved_ports"`
	Domains          int `json:"domains"`
	PerAppTasks      int `json:"per_app_tasks"`
	ServiceKeys      int `json:"service_keys"`
}

type QuotaRelationship struct {
	Quota ToOneRelationship `json:"quota"`
}

func NewOrganizationCreate(name string) *OrganizationCreate {
	return &OrganizationCreate{
		Name: name,
	}
}
