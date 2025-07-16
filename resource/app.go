package resource

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type App struct {
	Name          string           `json:"name"`
	State         string           `json:"state"`
	Lifecycle     Lifecycle        `json:"lifecycle"`
	Relationships AppRelationships `json:"relationships"`
	Metadata      *Metadata        `json:"metadata"`
	Resource      `json:",inline"`
}

type AppCreate struct {
	Name                 string            `json:"name"`
	Relationships        SpaceRelationship `json:"relationships"`
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	Lifecycle            *Lifecycle        `json:"lifecycle,omitempty"`
	Metadata             *Metadata         `json:"metadata,omitempty"`
}

type AppUpdate struct {
	Name      string     `json:"name"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	Metadata  *Metadata  `json:"metadata,omitempty"`
}

type AppList struct {
	Pagination Pagination   `json:"pagination"`
	Resources  []*App       `json:"resources"`
	Included   *AppIncluded `json:"included"`
}

type AppSSHEnabled struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason"`
}

type AppPermissions struct {
	ReadBasicData     bool `json:"read_basic_data"`
	ReadSensitiveData bool `json:"read_sensitive_data"`
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

// Lifecycle represents the Cloud Foundry V3 application lifecycle block.
// It supports buildpack, docker, and cnb (Cloud Native Buildpacks) types.
// Custom marshaling/unmarshaling is used to ensure correct JSON structure for each type.
// To add support for new lifecycle types, extend the marshaler/unmarshaler logic below.
type Lifecycle struct {
	Type string      `json:"type,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type BuildpackLifecycle struct {
	Buildpacks []string `json:"buildpacks,omitempty"`
	Stack      string   `json:"stack,omitempty"`
}

type DockerLifecycle struct {
	Image    string `json:"image,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type CNBLifecycle struct {
	Buildpacks []string `json:"buildpacks,omitempty"`
	Stack      string   `json:"stack,omitempty"`
}

type AppWithIncluded struct {
	App
	Included *AppIncluded `json:"included"`
}

type AppIncluded struct {
	Organizations []*Organization `json:"organizations"`
	Spaces        []*Space        `json:"spaces"`
}

// LifecycleType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#list-apps
type LifecycleType int

const (
	LifecycleNone LifecycleType = iota
	LifecycleBuildpack
	LifecycleDocker
)

func (l LifecycleType) String() string {
	switch l {
	case LifecycleBuildpack:
		return "buildpack"
	case LifecycleDocker:
		return "docker"
	default:
		return ""
	}
}

// AppIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type AppIncludeType int

const (
	AppIncludeNone AppIncludeType = iota
	AppIncludeSpace
	AppIncludeSpaceOrganization
)

func (a AppIncludeType) String() string {
	switch a {
	case AppIncludeSpace:
		return IncludeSpace
	case AppIncludeSpaceOrganization:
		return IncludeSpaceOrganization
	default:
		return IncludeNone
	}
}

func NewAppCreate(name, spaceGUID string) *AppCreate {
	return &AppCreate{
		Name: name,
		Relationships: SpaceRelationship{
			Space: ToOneRelationship{
				Data: &Relationship{
					GUID: spaceGUID,
				},
			},
		},
	}
}

type InterfaceEnvVar struct {
	Var map[string]interface{} `json:"var"`
}

func (f *EnvVar) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "" || s == "\"\"" {
		return nil
	}

	var interfaceEnvVar InterfaceEnvVar

	err := json.Unmarshal(b, &interfaceEnvVar)

	if err != nil {
		return err
	}

	varMap := make(map[string]*string)
	for key, value := range interfaceEnvVar.Var {
		switch v := value.(type) {
		default:
			fmt.Printf("unexpected env var type, skipping env var: %T", v)
			continue
		case int:
			stringVal := strconv.Itoa(v)
			varMap[key] = &stringVal
		case string:
			varMap[key] = &v
		case bool:
			stringBool := strconv.FormatBool(v)
			varMap[key] = &stringBool
		case float64:
			stringFloat := strconv.FormatFloat(v, 'f', -1, 64)
			varMap[key] = &stringFloat
		}
	}
	f.Var = varMap

	return nil
}

type InterfaceAppEnvironment struct {
	EnvVars       map[string]interface{}     `json:"environment_variables,omitempty"`
	StagingEnv    map[string]interface{}     `json:"staging_env_json,omitempty"`
	RunningEnv    map[string]interface{}     `json:"running_env_json,omitempty"`
	SystemEnvVars map[string]json.RawMessage `json:"system_env_json,omitempty"`      // VCAP_SERVICES
	AppEnvVars    map[string]json.RawMessage `json:"application_env_json,omitempty"` // VCAP_APPLICATION
}

func (f *AppEnvironment) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "" || s == "\"\"" {
		return nil
	}

	var interfaceAppEnvironment InterfaceAppEnvironment

	err := json.Unmarshal(b, &interfaceAppEnvironment)

	if err != nil {
		return err
	}

	f.EnvVars = convertEnvVars(interfaceAppEnvironment.EnvVars)
	f.StagingEnv = convertEnvVars(interfaceAppEnvironment.StagingEnv)
	f.RunningEnv = convertEnvVars(interfaceAppEnvironment.RunningEnv)
	f.AppEnvVars = interfaceAppEnvironment.AppEnvVars
	f.SystemEnvVars = interfaceAppEnvironment.SystemEnvVars

	return nil
}

func convertEnvVars(envVars map[string]interface{}) map[string]string {
	envVarMap := make(map[string]string)
	for key, value := range envVars {
		switch v := value.(type) {
		default:
			fmt.Printf("unexpected env var type, skipping env var: %T", v)
			continue
		case int:
			envVarMap[key] = strconv.Itoa(v)
		case string:
			envVarMap[key] = v
		case bool:
			envVarMap[key] = strconv.FormatBool(v)
		case float64:
			envVarMap[key] = strconv.FormatFloat(v, 'f', -1, 64)
		}
	}
	return envVarMap
}

// Update marshaling/unmarshaling logic for Lifecycle to handle multiple types
func (l *Lifecycle) MarshalJSON() ([]byte, error) {
	var data interface{}
	switch l.Type {
	case "buildpack":
		data = l.Data
	case "docker":
		data = l.Data
	case "cnb":
		data = l.Data
	default:
		data = nil
	}
	return json.Marshal(&struct {
		Type string      `json:"type,omitempty"`
		Data interface{} `json:"data,omitempty"`
	}{
		Type: l.Type,
		Data: data,
	})
}

func (l *Lifecycle) UnmarshalJSON(b []byte) error {
	var aux struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	l.Type = aux.Type
	switch aux.Type {
	case "buildpack":
		var bp BuildpackLifecycle
		if err := json.Unmarshal(aux.Data, &bp); err != nil {
			return err
		}
		l.Data = &bp
	case "docker":
		var d DockerLifecycle
		if err := json.Unmarshal(aux.Data, &d); err != nil {
			return err
		}
		l.Data = &d
	case "cnb":
		var cnb CNBLifecycle
		if err := json.Unmarshal(aux.Data, &cnb); err != nil {
			return err
		}
		l.Data = &cnb
	}
	return nil
}
