package actors

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestManifestMarshalling(t *testing.T) {
	m := &Manifest{
		Applications: []AppManifest{
			{
				Name:       "spring-music",
				Buildpacks: []string{"java_buildpack_offline"},
				Command:    "java",
				DiskQuota:  "1G",
				Env: map[string]string{
					"SPRING_CLOUD_PROFILE": "dev",
				},
				HealthCheckType:         "http",
				HealthCheckHTTPEndpoint: "/health",
				Instances:               2,
				LogRateLimit:            "100MB",
				Memory:                  "1G",
				NoRoute:                 false,
				Routes: []AppManifestRoutes{
					{"spring-music-egregious-porcupine-oa.apps.example.org"},
				},
				Services: []string{
					"my-sql",
				},
				Stack:   "cflinuxfs3",
				Timeout: 60,
			},
		},
	}
	b, err := yaml.Marshal(&m)
	require.NoError(t, err)
	require.Equal(t, fullSpringMusicYaml, string(b))

	m = &Manifest{
		Applications: []AppManifest{
			{
				Name:       "spring-music",
				Buildpacks: []string{"java_buildpack_offline"},
				Memory:     "1G",
				NoRoute:    true,
				Stack:      "cflinuxfs3",
			},
		},
	}
	b, err = yaml.Marshal(&m)
	require.NoError(t, err)
	require.Equal(t, minimalSpringMusicYaml, string(b))
}

const fullSpringMusicYaml = `applications:
- name: spring-music
  buildpacks:
  - java_buildpack_offline
  command: java
  disk_quota: 1G
  env:
    SPRING_CLOUD_PROFILE: dev
  health-check-type: http
  health-check-http-endpoint: /health
  instances: 2
  log-rate-limit: 100MB
  memory: 1G
  routes:
  - route: spring-music-egregious-porcupine-oa.apps.example.org
  services:
  - my-sql
  stack: cflinuxfs3
  timeout: 60
`

const minimalSpringMusicYaml = `applications:
- name: spring-music
  buildpacks:
  - java_buildpack_offline
  memory: 1G
  no-route: true
  stack: cflinuxfs3
`
