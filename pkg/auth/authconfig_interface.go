package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuthState saves data about the current clients authentication
type OAuthState struct {
	OAuth            *oauth2.Config
	OauthStateString string
}

// Config sets the environment configs for Facebook auth
type Config interface {
	Create()
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetAuthState() *OAuthState
	GetAPIData(accessToken string) (*http.Response, error)
}
