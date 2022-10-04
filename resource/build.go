package resource

type Build struct {
	State     string          `json:"state,omitempty"`
	Error     string          `json:"error,omitempty"`
	Lifecycle Lifecycle       `json:"lifecycle,omitempty"`
	Package   Relationship    `json:"package,omitempty"`
	Droplet   Relationship    `json:"droplet,omitempty"`
	GUID      string          `json:"guid,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
	CreatedBy CreatedBy       `json:"created_by,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  Metadata        `json:"metadata,omitempty"`
}

type CreatedBy struct {
	GUID  string `json:"guid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
