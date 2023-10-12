package operation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

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
