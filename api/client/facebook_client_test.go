package client

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"strings"

	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/testutils"
	cache "github.com/patrickmn/go-cache"
)

var mockState = "mockState"

func newFacebookClient(logs *testutils.MockLogger) *FacebookClient {
	facebookConfig := auth.MockAuthConfig{}
	facebookConfig.Create()
	facebookConfig.GetAuthState().OauthStateString = mockState

	facebook := FacebookClient{
		AuthConfig:      &facebookConfig,
		Cache:           cache.New(60*time.Minute, 30*time.Second),
		Logger:          logs,
		SuccessTemplate: "<p>Success</p>",
		FailureTemplate: "<p>Failure</p>",
	}
	return &facebook
}

func TestFBClientChecksForCorrectOAuthState(t *testing.T) {
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)

	if err := cli.checkOauthState(mockState); err != nil {
		t.Errorf("Expected state string to match oauth state, but got error %v", err.Error())
	}
}

func TestFBClientRetunsMissingOauthStateError(t *testing.T) {
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)

	if err := cli.checkOauthState(""); err == nil {
		t.Errorf("Expected state form string to return missing error")
	}

	cli.AuthConfig.GetAuthState().OauthStateString = ""
	if err := cli.checkOauthState(mockState); err == nil {
		t.Errorf("Expected oauth state string to return missing error")
	}
}

func TestFBClientReturnsNotMatchinOauthStateError(t *testing.T) {
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)

	if err := cli.checkOauthState("no-match"); err == nil {
		t.Errorf("Expected state form string to return not matching error")
	}
}

func TestFBClientReturnsTokenFromCode(t *testing.T) {
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)

	code := "mockCode"
	token, err := cli.getTokenFromFacebookCode(code)
	if err != nil {
		t.Fatalf("Expected code to create token, but got error %v", err.Error())
	}
	if token == nil {
		t.Errorf("Expected token to not be nil")
	}
	if token.AccessToken == "" {
		t.Errorf("Expected access token to not be empty")
	}
}

func TestFBClientReturnsCorrectResponseData(t *testing.T) {
	mockAccessToken := "mockAccessToken"
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)

	data, err := cli.getFacebookAPIData(mockAccessToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if data.ID == "" {
		t.Errorf("Expected facebook id to not by empty")
	}
	if data.Name == "" {
		t.Errorf("Expected facebook api name to not be empty")
	}
	if data.Token == "" {
		t.Errorf("Expected facebook api token to not be empty")
	}
	if strings.Compare(data.Token, mockAccessToken) != 0 {
		t.Errorf("Expected facebook api token to be %v, but got %v", mockAccessToken, data.Token)
	}
}

func TestServerReturnsSuccessfulFacebookCallback(t *testing.T) {
	mockLogger := testutils.MockLogger{}
	cli := newFacebookClient(&mockLogger)
	router := http.NewServeMux()

	router.Handle(FacebookCallback, cli.HandleFacebookCallback())

	data := url.Values{}
	data.Set("code", "test")
	data.Set("state", mockState)

	r, err := http.NewRequest("POST", FacebookCallback, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r.FormValue("code") == "" {
		t.Fatalf("Expected code form value to not be empty")
	}
	if r.FormValue("state") == "" {
		t.Fatalf("Expected state form value to not be empty")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if mockLogger.ErrorStackTraceCount != 0 {
		t.Errorf("Expected no stack trace errors, but got %v", mockLogger.ErrorCallCount)
	}
	if mockLogger.ErrorCallCount != 0 {
		t.Errorf("Expected no errors, but got %v", mockLogger.ErrorCallCount)
	}
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected %v, but got %v", http.StatusTemporaryRedirect, w.Code)
	}
}
