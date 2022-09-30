package v3

// SecurityGroup implements the security group object. Security groups are collections of egress traffic rules that can be applied to the staging or running state of applications.
type SecurityGroup struct {
	Name            string                         `json:"name,omitempty"`
	GUID            string                         `json:"guid,omitempty"`
	CreatedAt       string                         `json:"created_at,omitempty"`
	UpdatedAt       string                         `json:"updated_at,omitempty"`
	GloballyEnabled GloballyEnabled                `json:"globally_enabled,omitempty"`
	Rules           []Rule                         `json:"rules,omitempty"`
	Relationships   map[string]ToManyRelationships `json:"relationships,omitempty"`
	Links           map[string]Link                `json:"links,omitempty"`
}

// GloballyEnabled object controls if the group is applied globally to the lifecycle of all applications
type GloballyEnabled struct {
	Running bool `json:"running,omitempty"`
	Staging bool `json:"staging,omitempty"`
}

// Rule is an object that provide a rule that will be applied by a security group
type Rule struct {
	Protocol    string `json:"protocol,omitempty"`
	Destination string `json:"destination,omitempty"`
	Ports       string `json:"ports,omitempty"`
	Type        *int   `json:"type,omitempty"`
	Code        *int   `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Log         bool   `json:"log,omitempty"`
}

type ListSecurityGroupResponse struct {
	Pagination Pagination      `json:"pagination,omitempty"`
	Resources  []SecurityGroup `json:"resources,omitempty"`
}

// CreateSecurityGroupRequest implements an object that is passed to CreateSecurityGroup method
type CreateSecurityGroupRequest struct {
	Name            string                         `json:"name"`
	GloballyEnabled *GloballyEnabled               `json:"globally_enabled,omitempty"`
	Rules           []*Rule                        `json:"rules,omitempty"`
	Relationships   map[string]ToManyRelationships `json:"relationships,omitempty"`
}

// UpdateSecurityGroupRequest implements an object that is passed to UpdateSecurityGroup method
type UpdateSecurityGroupRequest struct {
	Name            string           `json:"name,omitempty"`
	GloballyEnabled *GloballyEnabled `json:"globally_enabled,omitempty"`
	Rules           []*Rule          `json:"rules,omitempty"`
}
