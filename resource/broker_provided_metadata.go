package resource

import (
	"fmt"
	"strings"
)

// BrokerProvidedMetadata allows the Service Broker to tag API resources with information that does not directly affect their functionality.
type BrokerProvidedMetadata struct {
	Labels     map[string]*string `json:"labels"`
	Attributes map[string]*string `json:"attributes"`
}

// NewBrokerProvidedMetadata creates a new broker provided metadata instance
func NewBrokerProvidedMetadata() *BrokerProvidedMetadata {
	return &BrokerProvidedMetadata{}
}

// WithAttribute is a fluent method alias for SetAttribute
func (m *BrokerProvidedMetadata) WithAttribute(prefix, key string, v string) *BrokerProvidedMetadata {
	m.SetAttribute(prefix, key, v)
	return m
}

// WithLabel is a fluent method alias for SetLabel
func (m *BrokerProvidedMetadata) WithLabel(prefix, key string, v string) *BrokerProvidedMetadata {
	m.SetLabel(prefix, key, v)
	return m
}

// SetAttribute to the broker provided metadata instance
//
// The prefix and value are optional and may be an empty string. The key must be at least 1 character in length.
func (m *BrokerProvidedMetadata) SetAttribute(prefix, key string, v string) {
	if m.Attributes == nil {
		m.Attributes = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Attributes[fmt.Sprintf("%s/%s", prefix, key)] = &v
	} else {
		m.Attributes[key] = &v
	}
}

// RemoveAttribute removes an attribute by setting the specified key's value to nil which can then be passed to the API
func (m *BrokerProvidedMetadata) RemoveAttribute(prefix, key string) {
	if m.Attributes == nil {
		m.Attributes = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Attributes[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Attributes[key] = nil
	}
}

// SetLabel to the broker provided metadata instance
//
// The prefix and value are optional and may be an empty string. The key must be at least 1 character in length.
func (m *BrokerProvidedMetadata) SetLabel(prefix, key string, v string) {
	if m.Labels == nil {
		m.Labels = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = &v
	} else {
		m.Labels[key] = &v
	}
}

// RemoveLabel removes a label by setting the specified key's value to nil which can then be passed to the API
func (m *BrokerProvidedMetadata) RemoveLabel(prefix, key string) {
	if m.Labels == nil {
		m.Labels = make(map[string]*string)
	}
	if len(prefix) > 0 {
		m.Labels[fmt.Sprintf("%s/%s", prefix, key)] = nil
	} else {
		m.Labels[key] = nil
	}
}

// Clear automatically calls Remove on all attributes and labels present in the broker provided metadata instance
func (m *BrokerProvidedMetadata) Clear() {
	splitKey := func(k string) (string, string) {
		p := strings.Split(k, "/")
		if len(p) == 1 {
			return "", p[0]
		}
		return p[0], p[1]
	}
	for k := range m.Attributes {
		prefix, key := splitKey(k)
		m.RemoveAttribute(prefix, key)
	}
	for k := range m.Labels {
		prefix, key := splitKey(k)
		m.RemoveLabel(prefix, key)
	}
}
