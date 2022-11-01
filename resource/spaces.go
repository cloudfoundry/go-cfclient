package resource

import (
	"time"
)

type Space struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Relationships map[string]ToOneRelationship `json:"relationships"`
	Links         map[string]Link              `json:"links"`
	Metadata      Metadata                     `json:"metadata"`
}

type SpaceCreate struct {
	Name          string                   `json:"name"`
	Relationships SpaceCreateRelationships `json:"relationships"`
	Metadata      *Metadata                `json:"metadata,omitempty"`
}
type SpaceCreateData struct {
	GUID string `json:"guid"`
}
type SpaceCreateOrganization struct {
	Data SpaceCreateData `json:"data"`
}
type SpaceCreateRelationships struct {
	Organization SpaceCreateOrganization `json:"organization"`
}

func NewSpaceCreate(name, orgGUID string) *SpaceCreate {
	return &SpaceCreate{
		Name: name,
		Relationships: SpaceCreateRelationships{
			Organization: SpaceCreateOrganization{
				Data: SpaceCreateData{
					GUID: orgGUID,
				},
			},
		},
	}
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

type SpaceUser struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	Relationships map[string]ToOneRelationship `json:"relationships"`
	Links         map[string]Link              `json:"links"`
	Metadata      Metadata                     `json:"metadata"`
}

type SpaceUserList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*User    `json:"resources"`
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
		return "organization"
	}
	return ""
}
