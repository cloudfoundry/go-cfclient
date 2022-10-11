package resource

import "time"

// User implements the user object
type User struct {
	GUID             string          `json:"guid"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Username         string          `json:"username,omitempty"`
	PresentationName string          `json:"presentation_name,omitempty"`
	Origin           string          `json:"origin,omitempty"`
	Links            map[string]Link `json:"links,omitempty"`
	Metadata         Metadata        `json:"metadata,omitempty"`
}

type ListUsersResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []User     `json:"resources,omitempty"`
}
