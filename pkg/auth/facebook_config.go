package auth

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// FacebookAuthConfig is auth configs needed to communicate with facebook
type FacebookAuthConfig struct {
	OAuthState *OAuthState
}

// Create a new facebook client with oauth config
func (c *FacebookAuthConfig) Create() {
	oauthstate := &OAuthState{}
	oauthstate.OAuth = &oauth2.Config{
		ClientID:     os.Getenv(config.FBClientID),
		ClientSecret: os.Getenv(config.FBSecretKey),
		RedirectURL:  os.Getenv(config.FacebookRedirectURL) + "/api/facebook/callback",
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	c.OAuthState = oauthstate
}

// Exchange checks and validates access token returned from facebook
func (c *FacebookAuthConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.OAuthState.OAuth.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// GetAuthState returns the auth config struct
func (c *FacebookAuthConfig) GetAuthState() *OAuthState {
	return c.OAuthState
}

// GetAPIData sends a request to facebook and returns the authenticated users name and id
func (c *FacebookAuthConfig) GetAPIData(accessToken string) (*http.Response, error) {
	facebookAccessURL := "https://graph.facebook.com/me?access_token="
	return http.Get(facebookAccessURL + url.QueryEscape(accessToken))
}
