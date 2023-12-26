package resource

type Space struct {
	Name          string              `json:"name"`
	Relationships *SpaceRelationships `json:"relationships"`
	Metadata      *Metadata           `json:"metadata"`
	Resource      `json:",inline"`
}

type SpaceCreate struct {
	Name          string              `json:"name"`
	Relationships *SpaceRelationships `json:"relationships"`
	Metadata      *Metadata           `json:"metadata,omitempty"`
}

type SpaceUpdate struct {
	Name     string    `json:"name,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type SpaceList struct {
	Pagination Pagination     `json:"pagination"`
	Resources  []*Space       `json:"resources"`
	Included   *SpaceIncluded `json:"included"`
}

type SpaceRelationships struct {
	Organization *ToOneRelationship `json:"organization"`
	Quota        *ToOneRelationship `json:"quota,omitempty"`
}

type SpaceWithIncluded struct {
	Space
	Included *SpaceIncluded `json:"included"`
}

type SpaceIncluded struct {
	Organizations []*Organization `json:"organizations"`
}

const (
	SpaceIncludeNone SpaceIncludeType = iota
	SpaceIncludeOrganization
)

type SpaceIncludeType int

func (s SpaceIncludeType) String() string {
	switch s {
	case SpaceIncludeOrganization:
		return IncludeOrganization
	default:
		return IncludeNone
	}
}

func NewSpaceCreate(name, orgGUID string) *SpaceCreate {
	return &SpaceCreate{
		Name: name,
		Relationships: &SpaceRelationships{
			Organization: &ToOneRelationship{
				Data: &Relationship{
					GUID: orgGUID,
				},
			},
		},
	}
}
