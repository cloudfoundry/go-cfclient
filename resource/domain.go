package resource

type Domain struct {
	Name               string              `json:"name"`
	Internal           bool                `json:"internal"`
	RouterGroup        *Relationship       `json:"router_group"`
	SupportedProtocols []string            `json:"supported_protocols"`
	Relationships      DomainRelationships `json:"relationships"`
	Metadata           *Metadata           `json:"metadata"`
	Resource           `json:",inline"`
}

type DomainCreate struct {
	Name string `json:"name"`

	Internal      *bool                `json:"internal,omitempty"`
	RouterGroup   *Relationship        `json:"router_group,omitempty"`
	Relationships *DomainRelationships `json:"relationships,omitempty"`
	Metadata      *Metadata            `json:"metadata,omitempty"`
}

type DomainUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type DomainList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Domain  `json:"resources"`
}

type DomainRelationships struct {
	Organization        *ToOneRelationship   `json:"organization,omitempty"`
	SharedOrganizations *ToManyRelationships `json:"shared_organizations,omitempty"`
}

func NewDomainCreate(name string) *DomainCreate {
	return &DomainCreate{
		Name: name,
	}
}

func NewDomainShare(orgGUID string) *ToManyRelationships {
	return &ToManyRelationships{
		Data: []Relationship{
			{
				GUID: orgGUID,
			},
		},
	}
}
