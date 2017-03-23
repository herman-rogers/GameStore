package auth

import (
	"context"
	"io/ioutil"
	"net/http"

	"bytes"

	"golang.org/x/oauth2"
)

// MockAuthConfig can be used for testing clients
type MockAuthConfig struct {
	OAuthState *OAuthState
}

// Create set up the OAuthState struct
func (c *MockAuthConfig) Create() {
	c.OAuthState = &OAuthState{}
}

// Exchange stubs out checking for oauth access and returns a mock access token
func (c *MockAuthConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token := oauth2.Token{}
	token.AccessToken = "MockAccessToken"
	return &token, nil
}

// GetAuthState returns mocked OauthStateString
func (c *MockAuthConfig) GetAuthState() *OAuthState {
	return c.OAuthState
}

// GetAPIData returns mock facebook data for ID and name
func (c *MockAuthConfig) GetAPIData(accessToken string) (*http.Response, error) {
	buffer := bytes.NewBuffer([]byte(`{"ID":"12345","name":"GameStore"}`))

	resp := http.Response{
		Body: ioutil.NopCloser(buffer),
	}
	return &resp, nil
}
