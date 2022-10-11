package resource

import "time"

// Role implements role object. Roles control access to resources in organizations and spaces. Roles are assigned to users.
type Role struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Type          string                       `json:"type,omitempty"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
}

type Included struct {
	Users         []User         `json:"users,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
	Spaces        []Space        `json:"spaces,omitempty"`
}

type ListRolesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Role     `json:"resources,omitempty"`
	Included   Included   `json:"included,omitempty"`
}

type CreateSpaceRoleRequest struct {
	RoleType      string                 `json:"type"`
	Relationships SpaceUserRelationships `json:"relationships"`
}

type CreateOrganizationRoleRequest struct {
	RoleType      string               `json:"type"`
	Relationships OrgUserRelationships `json:"relationships"`
}

type SpaceUserRelationships struct {
	Space ToOneRelationship `json:"space"`
	User  ToOneRelationship `json:"user"`
}

type OrgUserRelationships struct {
	Org  ToOneRelationship `json:"organization"`
	User ToOneRelationship `json:"user"`
}
