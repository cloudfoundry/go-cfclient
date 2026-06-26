package resource

type RoutePolicy struct {
	Source        string                   `json:"source"`
	Metadata      *Metadata                `json:"metadata"`
	Relationships RoutePolicyRelationships `json:"relationships"`
	Resource      `json:",inline"`
}

type RoutePolicyCreate struct {
	Relationships RoutePolicyRelationships `json:"relationships"`
	Source        string                   `json:"source"`
}

type RoutePolicyUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type RoutePolicyList struct {
	Pagination Pagination           `json:"pagination"`
	Resources  []*RoutePolicy       `json:"resources"`
	Included   *RoutePolicyIncluded `json:"included"`
}

type RoutePolicyRelationships struct {
	Route        *ToOneRelationship `json:"route,omitempty"`
	App          *ToOneRelationship `json:"app,omitempty"`
	Space        *ToOneRelationship `json:"space,omitempty"`
	Organization *ToOneRelationship `json:"organization,omitempty"`
}

type RoutePolicyWithIncluded struct {
	RoutePolicy RoutePolicy          `json:"route_policy"`
	Included    *RoutePolicyIncluded `json:"included"`
}

type RoutePolicyIncluded struct {
	Routes        []*Route        `json:"routes"`
	Apps          []*App          `json:"apps"`
	Spaces        []*Space        `json:"spaces"`
	Organizations []*Organization `json:"organizations"`
}

// RoutePolicyIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type RoutePolicyIncludeType int

const (
	RoutePolicyIncludeNone RoutePolicyIncludeType = iota
	RoutePolicyIncludeRoute
	RoutePolicyIncludeSource
	RoutePolicyIncludeRouteSource
)

func (r RoutePolicyIncludeType) String() string {
	switch r {
	case RoutePolicyIncludeRoute:
		return "route"
	case RoutePolicyIncludeSource:
		return "source"
	case RoutePolicyIncludeRouteSource:
		return "route,source"
	default:
		return ""
	}
}
