package resource

import "time"

type Build struct {
	GUID      string          `json:"guid"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	State     string          `json:"state,omitempty"`
	Error     string          `json:"error,omitempty"`
	Lifecycle Lifecycle       `json:"lifecycle,omitempty"`
	Package   Relationship    `json:"package,omitempty"`
	Droplet   Relationship    `json:"droplet,omitempty"`
	CreatedBy CreatedBy       `json:"created_by,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  Metadata        `json:"metadata,omitempty"`
}

type CreatedBy struct {
	GUID  string `json:"guid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
