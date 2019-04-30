package cfclient

import "fmt"

// Pagination is used by the V3 apis
type Pagination struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	First        Link        `json:"first"`
	Last         Link        `json:"last"`
	Next         interface{} `json:"next"`
	Previous     interface{} `json:"previous"`
}

// Link is a HATEOAS-style link for v3 apis
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

type Metadata struct {
	Metadata struct {
		Annotations map[string]interface{} `json:"annotations"`
		Labels      map[string]interface{} `json:"labels"`
	} `json:"metadata"`
}

func (m *Metadata) AddAnnotation(key string, value string) {
	if m.Metadata.Annotations == nil {
		m.Metadata.Annotations = make(map[string]interface{})
	}
	m.Metadata.Annotations[key] = value
}

func (m *Metadata) RemoveAnnotation(key string) {
	if m.Metadata.Annotations == nil {
		m.Metadata.Annotations = make(map[string]interface{})
	}
	m.Metadata.Annotations[key] = nil
}

func (m *Metadata) AddLabel(prefix, key string, value string) {
	if m.Metadata.Labels == nil {
		m.Metadata.Labels = make(map[string]interface{})
	}
	if len(prefix) > 0 {
		m.Metadata.Labels[fmt.Sprintf("%s/%s", prefix, key)] = value
	} else {
		m.Metadata.Labels[key] = value
	}
}

func (m *Metadata) RemoveLabel(prefix, key string) {
	if m.Metadata.Labels == nil {
		m.Metadata.Labels = make(map[string]interface{})
	}
	if len(prefix) > 0 {
		m.Metadata.Labels[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Metadata.Labels[key] = nil
	}
}
