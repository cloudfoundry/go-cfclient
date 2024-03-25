package operation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
)

const fullSpringMusicYaml = `applications:
- name: spring-music
  buildpacks:
  - java_buildpack_offline
  env:
    SPRING_CLOUD_PROFILE: dev
  routes:
  - route: spring-music-egregious-porcupine-oa.apps.example.org
  services:
  - name: my-sql
  stack: cflinuxfs3
  command: java
  disk_quota: 1G
  health-check-type: http
  health-check-http-endpoint: /health
  instances: 2
  log-rate-limit-per-second: 100MB
  memory: 1G
  timeout: 60
`

const minimalSpringMusicYaml = `applications:
- name: spring-music
  buildpacks:
  - java_buildpack_offline
  no-route: true
  stack: cflinuxfs3
  health-check-type: port
  health-check-http-endpoint: /
  instances: 1
  memory: 1G
`

const fullSpringMusicYamlV2 = `applications:
  - name: spring-music
    buildpacks:
      - java_buildpack_offline
    env:
      SPRING_CLOUD_PROFILE: dev
    routes:
      - route: spring-music-egregious-porcupine-oa.apps.example.org
    services:
      - name: my-sql
        binding_name: mysql
        parameters:
          name: mysql
      - oauth2
    stack: cflinuxfs3
    command: java
    disk_quota: 1G
    health-check-type: http
    health-check-http-endpoint: /health
    instances: 2
    log-rate-limit-per-second: 100MB
    memory: 1G
    timeout: 60
`

func TestManifestMarshalling(t *testing.T) {
	m := &Manifest{
		Applications: []*AppManifest{
			{
				Name:       "spring-music",
				Buildpacks: []string{"java_buildpack_offline"},
				Env: map[string]string{
					"SPRING_CLOUD_PROFILE": "dev",
				},
				NoRoute: false,
				Routes: &AppManifestRoutes{
					{Route: "spring-music-egregious-porcupine-oa.apps.example.org"},
				},
				Services: &AppManifestServices{
					{Name: "my-sql"},
				},
				Stack: "cflinuxfs3",
				AppManifestProcess: AppManifestProcess{
					HealthCheckType:         "http",
					HealthCheckHTTPEndpoint: "/health",
					Instances:               2,
					LogRateLimitPerSecond:   "100MB",
					Memory:                  "1G",
					Timeout:                 60,
					Command:                 "java",
					DiskQuota:               "1G",
				},
			},
		},
	}
	b, err := yaml.Marshal(&m)
	require.NoError(t, err)
	require.Equal(t, fullSpringMusicYaml, string(b))

	a := NewAppManifest("spring-music")
	a.Buildpacks = []string{"java_buildpack_offline"}
	a.Memory = "1G"
	a.NoRoute = true
	a.Stack = "cflinuxfs3"

	m = &Manifest{
		Applications: []*AppManifest{a},
	}
	b, err = yaml.Marshal(&m)
	require.NoError(t, err)
	require.Equal(t, minimalSpringMusicYaml, string(b))
}

func TestManifestUnMarshalling(t *testing.T) {
	var m Manifest

	err := yaml.Unmarshal([]byte(fullSpringMusicYamlV2), &m)
	require.NoError(t, err)
	require.Equal(t, "my-sql", (*m.Applications[0].Services)[0].Name)
	require.Equal(t, "mysql", (*m.Applications[0].Services)[0].BindingName)
	require.Equal(t, "mysql", (*m.Applications[0].Services)[0].Parameters["name"])
	require.Equal(t, "oauth2", (*m.Applications[0].Services)[1].Name)
	require.Equal(t, 1, len(m.Applications))
	require.Equal(t, 2, len(*m.Applications[0].Services))

	err = yamlv3.Unmarshal([]byte(fullSpringMusicYamlV2), &m)
	require.NoError(t, err)
	require.Equal(t, "my-sql", (*m.Applications[0].Services)[0].Name)
	require.Equal(t, "mysql", (*m.Applications[0].Services)[0].BindingName)
	require.Equal(t, "mysql", (*m.Applications[0].Services)[0].Parameters["name"])
	require.Equal(t, "oauth2", (*m.Applications[0].Services)[1].Name)
	require.Equal(t, 1, len(m.Applications))
	require.Equal(t, 2, len(*m.Applications[0].Services))
}
