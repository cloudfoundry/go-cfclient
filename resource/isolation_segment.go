package resource

type IsolationSegment struct {
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata"`
	Resource `json:",inline"`
}

type IsolationSegmentCreate struct {
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type IsolationSegmentUpdate struct {
	Name     *string   `json:"name,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type IsolationSegmentRelationship struct {
	Data  []Relationship  `json:"data"`
	Links map[string]Link `json:"links"`
}

type IsolationSegmentList struct {
	Pagination Pagination          `json:"pagination"`
	Resources  []*IsolationSegment `json:"resources"`
}

func NewIsolationSegmentCreate(name string) *IsolationSegmentCreate {
	return &IsolationSegmentCreate{
		Name: name,
	}
}
