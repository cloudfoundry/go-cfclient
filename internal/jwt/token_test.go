package jwt_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/jwt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const accessToken = `bearer ignored.eyJqdGkiOiJhOGE5YTJjNDY5MzY0YzU3YmI2M2QxMWFiYzdhNjgzOSIsInN1YiI6IjJiNmMzM2ZlLTExZTItNGQwMi05OTNhLTdiNjQ5ZjhhMmI5YyIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiMmI2YzMzZmUtMTFlMi00ZDAyLTk5M2EtN2I2NDlmOGEyYjljIiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aF90aW1lIjoxNjk4MDk2Mzc2LCJyZXZfc2lnIjoiZmNlMmY2MDAiLCJpYXQiOjE2OTgwOTY0MDgsImV4cCI6MTY5ODA5NjQ2OCwiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmgyby0yLTE5MTQ5Lmgyby52bXdhcmUuY29tL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJzY2ltIiwidWFhIiwibmV0d29yayIsImNmIl19.ignored`

func TestParseAccessToken(t *testing.T) {
	expireTime, err := jwt.AccessTokenExpiration(accessToken)
	require.NoError(t, err)
	expected := time.Date(2023, 10, 23, 14, 27, 48, 0, time.FixedZone("UTC-7", -7*60*60))
	require.Equal(t, expected.Unix(), expireTime.Unix())

	_, err = jwt.AccessTokenExpiration("")
	require.EqualError(t, err, "access token format is invalid")

	_, err = jwt.AccessTokenExpiration("not base.64encoded.token")
	require.EqualError(t, err, "access token base64 encoding is invalid")

	_, err = jwt.AccessTokenExpiration("bearer ignored.eyJqdGkiOiJhOGE5YTJjNDY5MzY0YzU3YmI2M2QxMW.ignored")
	require.EqualError(t, err, "access token is invalid: unexpected end of JSON input")
}
