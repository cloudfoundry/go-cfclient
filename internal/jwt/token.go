package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type tokenPayload struct {
	Expiration int64 `json:"exp"`
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
