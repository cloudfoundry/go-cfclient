package resource

type OrganizationQuota struct {
	Name string `json:"name"`

	Apps     AppsQuota     `json:"apps"`
	Services ServicesQuota `json:"services"`
	Routes   RoutesQuota   `json:"routes"`
	Domains  DomainsQuota  `json:"domains"`

	Relationships OrganizationQuotaRelationships `json:"relationships"`

	Resource `json:",inline"`
}

type OrganizationQuotaCreateOrUpdate struct {
	Name          *string                         `json:"name,omitempty"`
	Apps          *AppsQuota                      `json:"apps,omitempty"`
	Services      *ServicesQuota                  `json:"services,omitempty"`
	Routes        *RoutesQuota                    `json:"routes,omitempty"`
	Domains       *DomainsQuota                   `json:"domains,omitempty"`
	Relationships *OrganizationQuotaRelationships `json:"relationships,omitempty"`
}

type OrganizationQuotaList struct {
	Pagination Pagination           `json:"pagination"`
	Resources  []*OrganizationQuota `json:"resources"`
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
		q.Apps = &AppsQuota{}
	}
	q.Apps.TotalMemoryInMB = &mb
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPerProcessMemoryInMB(mb int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &AppsQuota{}
	}
	q.Apps.PerProcessMemoryInMB = &mb
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithLogRateLimitInBytesPerSecond(bytes int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &AppsQuota{}
	}
	q.Apps.LogRateLimitInBytesPerSecond = &bytes
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalInstances(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &AppsQuota{}
	}
	q.Apps.TotalInstances = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPerAppTasks(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Apps == nil {
		q.Apps = &AppsQuota{}
	}
	q.Apps.PerAppTasks = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithPaidServicesAllowed(allowed bool) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &ServicesQuota{}
	}
	q.Services.PaidServicesAllowed = &allowed
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalServiceInstances(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &ServicesQuota{}
	}
	q.Services.TotalServiceInstances = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalServiceKeys(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Services == nil {
		q.Services = &ServicesQuota{}
	}
	q.Services.TotalServiceKeys = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalRoutes(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Routes == nil {
		q.Routes = &RoutesQuota{}
	}
	q.Routes.TotalRoutes = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithTotalReservedPorts(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Routes == nil {
		q.Routes = &RoutesQuota{}
	}
	q.Routes.TotalReservedPorts = &count
	return q
}

func (q *OrganizationQuotaCreateOrUpdate) WithDomains(count int) *OrganizationQuotaCreateOrUpdate {
	if q.Domains == nil {
		q.Domains = &DomainsQuota{}
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
