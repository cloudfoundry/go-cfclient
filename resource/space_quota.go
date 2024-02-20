package resource

type SpaceQuota struct {
	// 	Name of the quota
	Name string `json:"name"`

	// quotas that affect apps
	Apps AppsQuota `json:"apps"`

	// quotas that affect services
	Services ServicesQuota `json:"services"`

	// quotas that affect routes
	Routes RoutesQuota `json:"routes"`

	// relationships to the organization and spaces where the quota belongs
	Relationships SpaceQuotaRelationships `json:"relationships"`

	Resource `json:",inline"`
}

type SpaceQuotaList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*SpaceQuota `json:"resources"`
}

type SpaceQuotaCreateOrUpdate struct {
	// 	Name of the quota
	Name *string `json:"name,omitempty"`

	// relationships to the organization and spaces where the quota belongs
	Relationships *SpaceQuotaRelationships `json:"relationships,omitempty"`

	// quotas that affect apps
	Apps *AppsQuota `json:"apps,omitempty"`

	// quotas that affect services
	Services *ServicesQuota `json:"services,omitempty"`

	// quotas that affect routes
	Routes *RoutesQuota `json:"routes,omitempty"`
}

type SpaceQuotaRelationships struct {
	Organization *ToOneRelationship   `json:"organization,omitempty"`
	Spaces       *ToManyRelationships `json:"spaces,omitempty"`
}

func NewSpaceQuotaCreate(name, orgGUID string) *SpaceQuotaCreateOrUpdate {
	return &SpaceQuotaCreateOrUpdate{
		Name: &name,
		Relationships: &SpaceQuotaRelationships{
			Organization: &ToOneRelationship{
				Data: &Relationship{
					GUID: orgGUID,
				},
			},
		},
	}
}

func NewSpaceQuotaUpdate() *SpaceQuotaCreateOrUpdate {
	return &SpaceQuotaCreateOrUpdate{}
}

func (s *SpaceQuotaCreateOrUpdate) WithName(name string) *SpaceQuotaCreateOrUpdate {
	s.Name = &name
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalMemoryInMB(mb int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &AppsQuota{}
	}
	s.Apps.TotalMemoryInMB = &mb
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPerProcessMemoryInMB(mb int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &AppsQuota{}
	}
	s.Apps.PerProcessMemoryInMB = &mb
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithLogRateLimitInBytesPerSecond(mbps int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &AppsQuota{}
	}
	s.Apps.LogRateLimitInBytesPerSecond = &mbps
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalInstances(count int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &AppsQuota{}
	}
	s.Apps.TotalInstances = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPerAppTasks(count int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &AppsQuota{}
	}
	s.Apps.PerAppTasks = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPaidServicesAllowed(allowed bool) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &ServicesQuota{}
	}
	s.Services.PaidServicesAllowed = &allowed
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalServiceInstances(count int) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &ServicesQuota{}
	}
	s.Services.TotalServiceInstances = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalServiceKeys(count int) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &ServicesQuota{}
	}
	s.Services.TotalServiceKeys = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalRoutes(count int) *SpaceQuotaCreateOrUpdate {
	if s.Routes == nil {
		s.Routes = &RoutesQuota{}
	}
	s.Routes.TotalRoutes = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalReservedPorts(count int) *SpaceQuotaCreateOrUpdate {
	if s.Routes == nil {
		s.Routes = &RoutesQuota{}
	}
	s.Routes.TotalReservedPorts = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithSpaces(spaceGUIDs ...string) *SpaceQuotaCreateOrUpdate {
	if s.Relationships == nil {
		s.Relationships = &SpaceQuotaRelationships{
			Spaces: &ToManyRelationships{},
		}
	}
	for _, g := range spaceGUIDs {
		r := Relationship{
			GUID: g,
		}
		s.Relationships.Spaces.Data = append(s.Relationships.Spaces.Data, r)
	}
	return s
}
