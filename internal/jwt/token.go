package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	grantTypeJwtBearer  = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	clientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
)

type tokenPayload struct {
	Expiration int64 `json:"exp"`
}

type JWTAssertionTokenSource struct { // revive:disable-line:exported
	Assertion       string
	ClientAssertion string
	ClientID        string
	ClientSecret    string
	GrantType       string
	Scopes          []string
	TokenURL        string
	HTTPClient      *http.Client
}

func AccessTokenExpiration(accessToken string) (time.Time, error) {
	tp := strings.Split(accessToken, ".")
	if len(tp) != 3 {
		return time.Time{}, errors.New("access token format is invalid")
	}

	// Decode the payload segment
	decoded, err := base64.RawURLEncoding.DecodeString(tp[1])
	if err != nil {
		return time.Time{}, errors.New("access token base64 encoding is invalid")
	}

	var t tokenPayload
	if err := json.Unmarshal(decoded, &t); err != nil {
		return time.Time{}, fmt.Errorf("access token is invalid: %w", err)
	}

	return time.Unix(t.Expiration, 0), nil
}

// ToOAuth2Token converts access token and refresh token to an oauth2.Token.
func ToOAuth2Token(accessToken, refreshToken string) (*oauth2.Token, error) {
	accessToken = strings.TrimSpace(accessToken)
	refreshToken = strings.TrimSpace(refreshToken)
	if accessToken == "" && refreshToken == "" {
		return nil, errors.New("expected a non-empty CF API access token or refresh token")
	}
	oAuthToken := &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "bearer", // Default token type
	}
	if accessToken != "" {
		tokens := strings.SplitN(accessToken, " ", 2)
		if len(tokens) > 1 {
			oAuthToken.TokenType = strings.ToLower(tokens[0])
		}

		token := tokens[len(tokens)-1]

		exp, err := AccessTokenExpiration(token)
		if err != nil {
			return nil, fmt.Errorf("error decoding token: %w", err)
		}

		oAuthToken.AccessToken = token
		oAuthToken.Expiry = exp
	}
	return oAuthToken, nil
}

func (s *JWTAssertionTokenSource) Token() (*oauth2.Token, error) {
	data := url.Values{}

	if s.TokenURL == "" {
		return nil, fmt.Errorf("token URL is required")
	}
	if s.GrantType == "" {
		data.Set("grant_type", grantTypeJwtBearer)
	} else {
		data.Set("grant_type", s.GrantType)
	}
	// Assertion is required for JWT Bearer grant type
	if s.Assertion == "" && (s.GrantType == grantTypeJwtBearer || s.GrantType == "") {
		return nil, fmt.Errorf("assertion is required for JWT Bearer grant type")
	}

	if s.Assertion != "" {
		if err := validateJWTTokenFormat(s.Assertion); err != nil {
			return nil, err
		}
		data.Set("assertion", s.Assertion)
	}

	// Optional client_id
	if s.ClientID != "" {
		data.Set("client_id", s.ClientID)
	}

	// Optional client_secret
	if s.ClientSecret != "" {
		if s.ClientID == "" {
			return nil, fmt.Errorf("client_id is required when using client secret")
		}
		data.Set("client_secret", s.ClientSecret)
	}

	// Optional client_assertion
	if s.ClientAssertion != "" {
		if s.ClientID == "" {
			return nil, fmt.Errorf("client_id is required when using client assertion")
		}
		if err := validateJWTTokenFormat(s.ClientAssertion); err != nil {
			return nil, err
		}
		data.Set("client_assertion_type", clientAssertionType)
		data.Set("client_assertion", s.ClientAssertion)
	}
	if len(s.Scopes) > 0 {
		data.Set("scope", strings.Join(s.Scopes, " "))
	}

	req, err := http.NewRequest("POST", s.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token request object creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := s.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed: %s", body)
	}

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, fmt.Errorf("token unmarshal error: %w", err)
	}
	return &token, nil
}

// validateJWTTokenFormat checks if the provided JWT token has a valid format.
func validateJWTTokenFormat(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("token must have three parts separated by '.'")
	}

	for i, part := range parts {
		if part == "" {
			return fmt.Errorf("token part is empty")
		}
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			return fmt.Errorf("invalid base64 encoding in part %d", i+1)
		}
	}

	return nil
}
