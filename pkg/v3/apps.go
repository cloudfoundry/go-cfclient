package v3

import (
	"encoding/json"
)

type App struct {
	Name          string                       `json:"name,omitempty"`
	State         string                       `json:"state,omitempty"`
	Lifecycle     Lifecycle                    `json:"lifecycle,omitempty"`
	GUID          string                       `json:"guid,omitempty"`
	CreatedAt     string                       `json:"created_at,omitempty"`
	UpdatedAt     string                       `json:"updated_at,omitempty"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
	Metadata      Metadata                     `json:"metadata,omitempty"`
}

type Lifecycle struct {
	Type          string             `json:"type,omitempty"`
	BuildpackData BuildpackLifecycle `json:"data,omitempty"`
}

type BuildpackLifecycle struct {
	Buildpacks []string `json:"buildpacks,omitempty"`
	Stack      string   `json:"stack,omitempty"`
}

type CreateAppRequest struct {
	Name                 string
	SpaceGUID            string
	EnvironmentVariables map[string]string
	Lifecycle            *Lifecycle
	Metadata             *Metadata
}

type UpdateAppRequest struct {
	Name      string     `json:"name"`
	Lifecycle *Lifecycle `json:"lifecycle"`
	Metadata  *Metadata  `json:"metadata"`
}

type AppEnvironment struct {
	EnvVars       map[string]string          `json:"environment_variables,omitempty"`
	StagingEnv    map[string]string          `json:"staging_env_json,omitempty"`
	RunningEnv    map[string]string          `json:"running_env_json,omitempty"`
	SystemEnvVars map[string]json.RawMessage `json:"system_env_json,omitempty"`      // VCAP_SERVICES
	AppEnvVars    map[string]json.RawMessage `json:"application_env_json,omitempty"` // VCAP_APPLICATION
}

type EnvVar struct {
	Var map[string]*string `json:"var"`
}

type EnvVarResponse struct {
	EnvVar
	Links map[string]Link `json:"links"`
}

type ListAppsResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []App      `json:"resources,omitempty"`
}
