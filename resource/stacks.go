package resource

import "time"

// Stack implements stack object. Stacks are the base operating system and file system that your application will execute in. A stack is how you configure applications to run against different operating systems (like Windows or Linux) and different versions of those operating systems.
type Stack struct {
	GUID        string          `json:"guid"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Links       map[string]Link `json:"links,omitempty"`
	Metadata    Metadata        `json:"metadata,omitempty"`
}

type ListStacksResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []Stack    `json:"resources,omitempty"`
}
