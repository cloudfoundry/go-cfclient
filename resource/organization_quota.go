package resource

type OrganizationQuota struct {
	Name string `json:"name"`

	Apps     OrganizationQuotaApps     `json:"apps"`
	Services OrganizationQuotaServices `json:"services"`
	Routes   OrganizationQuotaRoutes   `json:"routes"`
	Domains  OrganizationQuotaDomains  `json:"domains"`

	Relationships OrganizationQuotaRelationships `json:"relationships"`

	Resource `json:",inline"`
}

type OrganizationQuotaCreateOrUpdate struct {
	Name          *string                         `json:"name,omitempty"`
	Apps          *OrganizationQuotaApps          `json:"apps,omitempty"`
	Services      *OrganizationQuotaServices      `json:"services,omitempty"`
	Routes        *OrganizationQuotaRoutes        `json:"routes,omitempty"`
	Domains       *OrganizationQuotaDomains       `json:"domains,omitempty"`
	Relationships *OrganizationQuotaRelationships `json:"relationships,omitempty"`
}

type OrganizationQuotaList struct {
	Pagination Pagination           `json:"pagination"`
	Resources  []*OrganizationQuota `json:"resources"`
}

type OrganizationQuotaApps struct {
	// Total memory allowed for all the started processes and running tasks in an organization
	TotalMemoryInMB *int `json:"total_memory_in_mb"`

	// Maximum memory for a single process or task
	PerProcessMemoryInMB *int `json:"per_process_memory_in_mb"`

	// Total log rate limit allowed for all the started processes and running tasks in an organization
	LogRateLimitInBytesPerSecond *int `json:"log_rate_limit_in_bytes_per_second"`

	// Total instances of all the started processes allowed in an organization
	TotalInstances *int `json:"total_instances"`

	// Maximum number of running tasks in an organization
	PerAppTasks *int `json:"per_app_tasks"`
}

type OrganizationQuotaServices struct {
	// Specifies whether instances of paid service plans can be created
	PaidServicesAllowed *bool `json:"paid_services_allowed"`

	// Total number of service instances allowed in an organization
	TotalServiceInstances *int `json:"total_service_instances"`

	// Total number of service keys allowed in an organization
	TotalServiceKeys *int `json:"total_service_keys"`
}

type OrganizationQuotaRoutes struct {
	// Total number of routes allowed in an organization
	TotalRoutes *int `json:"total_routes"`

	// Total number of ports that are reservable by routes in an organization
	TotalReservedPorts *int `json:"total_reserved_ports"`
}

type OrganizationQuotaDomains struct {
	// Total number of domains that can be scoped to an organization
	TotalDomains *int `json:"total_domains"`
}

type OrganizationQuotaRelationships struct {
	// A relationship to the organizations where the quota is applied
	Organizations ToManyRelationships `json:"organizations"`
}

func NewOrganizationQuotaCreate(name string) *OrganizationQuotaCreateOrUpdate {
	return &OrganizationQuotaCreateOrUpdate{
		Name: &name,
	}
}

func NewOrganizationQuotaUpdate() *OrganizationQuotaCreateOrUpdate {
	return &OrganizationQuotaCreateOrUpdate{}
}

func (q *OrganizationQuotaCreateOrUpdate) WithName(name string) *OrganizationQuotaCreateOrUpdate {
	q.Name = &name
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithAppsTotalMemoryInMB(mb int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &OrganizationQuotaApps{}
	}
	q.Apps.TotalMemoryInMB = &mb
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPerProcessMemoryInMB(mb int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &OrganizationQuotaApps{}
	}
	q.Apps.PerProcessMemoryInMB = &mb
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithLogRateLimitInBytesPerSecond(bytes int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &OrganizationQuotaApps{}
	}
	q.Apps.LogRateLimitInBytesPerSecond = &bytes
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalInstances(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &OrganizationQuotaApps{}
	}
	q.Apps.TotalInstances = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPerAppTasks(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &OrganizationQuotaApps{}
	}
	q.Apps.PerAppTasks = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPaidServicesAllowed(allowed bool) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &OrganizationQuotaServices{}
	}
	q.Services.PaidServicesAllowed = &allowed
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalServiceInstances(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &OrganizationQuotaServices{}
	}
	q.Services.TotalServiceInstances = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalServiceKeys(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &OrganizationQuotaServices{}
	}
	q.Services.TotalServiceKeys = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalRoutes(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Routes == nil {
		q.Routes = &OrganizationQuotaRoutes{}
	}
	q.Routes.TotalRoutes = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalReservedPorts(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Routes == nil {
		q.Routes = &OrganizationQuotaRoutes{}
	}
	q.Routes.TotalReservedPorts = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithDomains(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Domains == nil {
		q.Domains = &OrganizationQuotaDomains{}
	}
	q.Domains.TotalDomains = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithOrganizations(orgGUIDs ...string) *OrganizationQuotaCreateOrUpdate {
	if q.Relationships == nil {
		q.Relationships = &OrganizationQuotaRelationships{
			Organizations: ToManyRelationships{},
		}
	}
	for _, g := range orgGUIDs {
		r := Relationship{
			GUID: g,
		}
		q.Relationships.Organizations.Data = append(q.Relationships.Organizations.Data, r)
	}
	return q
}
