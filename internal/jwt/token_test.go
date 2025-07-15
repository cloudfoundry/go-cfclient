package jwt

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	validToken = "a.test.token"

	validAssertionToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            // revive:disable:line-length-limit
	accessToken         = `bearer ignored.eyJqdGkiOiJhOGE5YTJjNDY5MzY0YzU3YmI2M2QxMWFiYzdhNjgzOSIsInN1YiI6IjJiNmMzM2ZlLTExZTItNGQwMi05OTNhLTdiNjQ5ZjhhMmI5YyIsInNjb3BlIjpbIm9wZW5pZCIsInJvdXRpbmcucm91dGVyX2dyb3Vwcy53cml0ZSIsIm5ldHdvcmsud3JpdGUiLCJzY2ltLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLmFkbWluIiwidWFhLnVzZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIucmVhZCIsInBhc3N3b3JkLndyaXRlIiwiY2xvdWRfY29udHJvbGxlci53cml0ZSIsIm5ldHdvcmsuYWRtaW4iLCJkb3BwbGVyLmZpcmVob3NlIiwic2NpbS53cml0ZSJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiMmI2YzMzZmUtMTFlMi00ZDAyLTk5M2EtN2I2NDlmOGEyYjljIiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoiYWRtaW4iLCJlbWFpbCI6ImFkbWluIiwiYXV0aF90aW1lIjoxNjk4MDk2Mzc2LCJyZXZfc2lnIjoiZmNlMmY2MDAiLCJpYXQiOjE2OTgwOTY0MDgsImV4cCI6MTY5ODA5NjQ2OCwiaXNzIjoiaHR0cHM6Ly91YWEuc3lzLmgyby0yLTE5MTQ5Lmgyby52bXdhcmUuY29tL29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImRvcHBsZXIiLCJyb3V0aW5nLnJvdXRlcl9ncm91cHMiLCJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJzY2ltIiwidWFhIiwibmV0d29yayIsImNmIl19.ignored` // revive:disable:line-length-limit
)

func mockServer(statusCode int, responseBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if _, err := fmt.Fprint(w, responseBody); err != nil {
			log.Fatalf("failed to write response body: %v", err)
		}
	}))
}

func TestToken(t *testing.T) {
	t.Run("Test AccessTokenExpiration", func(t *testing.T) {
		expireTime, err := AccessTokenExpiration(accessToken)
		require.NoError(t, err)
		expected := time.Date(2023, 10, 23, 14, 27, 48, 0, time.FixedZone("UTC-7", -7*60*60))
		require.Equal(t, expected.Unix(), expireTime.Unix())

		_, err = AccessTokenExpiration("")
		require.EqualError(t, err, "access token format is invalid")

		_, err = AccessTokenExpiration("not base.64encoded.token")
		require.EqualError(t, err, "access token base64 encoding is invalid")

		_, err = AccessTokenExpiration("bearer ignored.eyJqdGkiOiJhOGE5YTJjNDY5MzY0YzU3YmI2M2QxMW.ignored")
		require.EqualError(t, err, "access token is invalid: unexpected end of JSON input")
	})

	t.Run("Test ToOAuth2Token", func(t *testing.T) {
		_, err := ToOAuth2Token("", "")
		require.EqualError(t, err, "expected a non-empty CF API access token or refresh token")

		token, err := ToOAuth2Token(accessToken, "")
		require.NoError(t, err)
		require.NotNil(t, token)
		require.Equal(t, "bearer", token.TokenType)
	})

}

func TestTokenSource(t *testing.T) {
	t.Run("Test JWTAssertionTokenSource", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			Assertion:       validAssertionToken,
			ClientAssertion: validAssertionToken,
			ClientID:        "test-client-id",
			ClientSecret:    "test-client-secret",
			GrantType:       grantTypeJwtBearer,
			TokenURL:        ts.URL,
			HTTPClient:      http.DefaultClient,
		}

		token, err := src.Token()
		require.NoError(t, err)
		require.NotNil(t, token)
	})
	t.Run("Test JWTAssertionTokenSource minimal", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			Assertion: validAssertionToken,
			TokenURL:  ts.URL,
		}

		token, err := src.Token()
		require.NoError(t, err)
		require.NotNil(t, token)
	})
	t.Run("Test JWTAssertionTokenSource invalid asssertion format", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			Assertion: validAssertionToken + `additional.invalid.part`,
			TokenURL:  ts.URL,
		}

		_, err := src.Token()
		require.EqualError(t, err, "token must have three parts separated by '.'")

	})
	t.Run("Test JWTAssertionTokenSource invalid encoding", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			Assertion: validAssertionToken + ` `,
			TokenURL:  ts.URL,
		}

		_, err := src.Token()
		require.EqualError(t, err, "invalid base64 encoding in part 3")

	})
	t.Run("Test JWTAssertionTokenSource without assertion and GrantType JWTBearer", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			TokenURL: ts.URL,
		}
		_, err := src.Token()
		require.EqualError(t, err, "assertion is required for JWT Bearer grant type")
	})
	t.Run("Test JWTAssertionTokenSource without Token URL", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			Assertion: validAssertionToken + ` `,
		}
		_, err := src.Token()
		require.EqualError(t, err, "token URL is required")
	})
	t.Run("Test client_credential with assertion flow", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			ClientAssertion: validAssertionToken,
			ClientID:        "test-client-id",
			GrantType:       `client_credentials`,
			TokenURL:        ts.URL,
		}
		token, err := src.Token()
		require.NoError(t, err)
		require.NotNil(t, token)
	})
	t.Run("Test client_credential assertion without client id", func(t *testing.T) {
		ts := mockServer(200, `{"access_token":"`+validToken+`","token_type":"Bearer","expires_in":3600}`)
		defer ts.Close()
		src := &JWTAssertionTokenSource{
			ClientAssertion: validAssertionToken,
			GrantType:       `client_credentials`,
			TokenURL:        ts.URL,
		}
		_, err := src.Token()
		require.EqualError(t, err, "client_id is required when using client assertion")
	})
}
