package resource

type Meta struct {
	GUID      string `json:"guid"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Pagination is used by the  apis
type Pagination struct {
	TotalResults int  `json:"total_results"`
	TotalPages   int  `json:"total_pages"`
	First        Link `json:"first"`
	Last         Link `json:"last"`
	Next         Link `json:"next"`
	Previous     Link `json:"previous"`
}

// Link is a HATEOAS-style link for  apis
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

// ToOneRelationship is a relationship to a single object
type ToOneRelationship struct {
	Data Relationship `json:"data,omitempty"`
}

// ToManyRelationships is a relationship to multiple objects
type ToManyRelationships struct {
	Data []Relationship `json:"data,omitempty"`
}

type Relationship struct {
	GUID string `json:"guid,omitempty"`
}

type Metadata struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
