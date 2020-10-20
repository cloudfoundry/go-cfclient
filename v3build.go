package cfclient

type V3Build struct {
	State     string          `json:"state,omitempty"`
	Error     string          `json:"error,omitempty"`
	Lifecycle V3Lifecycle     `json:"lifecycle,omitempty"`
	Package   V3Relationship  `json:"package,omitempty"`
	Droplet   V3Relationship  `json:"droplet,omitempty"`
	GUID      string          `json:"guid,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	UpdatedAt string          `json:"updated_at,omitempty"`
	CreatedBy V3CreatedBy     `json:"created_by,omitempty"`
	Links     map[string]Link `json:"links,omitempty"`
	Metadata  V3Metadata      `json:"metadata,omitempty"`
}

type V3CreatedBy struct {
	GUID  string `json:"guid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
