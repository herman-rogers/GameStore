package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
	"github.com/go-errors/errors"
	cache "github.com/patrickmn/go-cache"
)

// FacebookClient is the client to interact with facebook's api
type FacebookClient struct {
	AuthConfig      auth.Config
	Cache           *cache.Cache
	Logger          logger.Log
	SuccessTemplate string
	FailureTemplate string
}

// FacebookResponse is the data returned from facebook's api
type FacebookResponse struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

// HandleFacebookLoginSuccess renders a login success template
func (cli *FacebookClient) HandleFacebookLoginSuccess() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, cli.SuccessTemplate)
	})
}

// HandleFacebookLoginFailed renders a login failed template
func (cli *FacebookClient) HandleFacebookLoginFailed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, cli.FailureTemplate)
	})
}

// GetFacebookData endpoint called by clients to get their auth tokens and facebook user data
// Request requires an identity in the format ?identity=IDENTITY
func (cli *FacebookClient) GetFacebookData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity := r.FormValue("identity")
		if token, found := cli.Cache.Get(identity); found {
			w.WriteHeader(http.StatusOK)
			w.Write(token.([]byte))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})
}

// DeleteFacebookData will remove a specified users auth token from the servers cache
// Request requires an identity in the format ?identity=IDENTITY
func (cli *FacebookClient) DeleteFacebookData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity := r.FormValue("identity")
		cli.Cache.Delete(identity)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted"))
	})
}

// HandleFacebookLogin sends a login request with RPC parameter ?token=TOKEN_DATA
func (cli *FacebookClient) HandleFacebookLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cli.AuthConfig.GetAuthState().OAuth
		cli.AuthConfig.GetAuthState().OauthStateString = r.FormValue("token")

		URL, err := url.Parse(c.Endpoint.AuthURL)
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		parameters := url.Values{}
		parameters.Add("client_id", c.ClientID)
		parameters.Add("scope", strings.Join(c.Scopes, " "))
		parameters.Add("redirect_uri", c.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", cli.AuthConfig.GetAuthState().OauthStateString)
		URL.RawQuery = parameters.Encode()

		stringURL := URL.String()
		http.Redirect(w, r, stringURL, http.StatusTemporaryRedirect)
	})
}

// HandleFacebookCallback is the callback response after requesting authentication to facebook's api
func (cli *FacebookClient) HandleFacebookCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")

		if err := cli.checkOauthState(state); err != nil {
			cli.Logger.Error(err.Error())
			http.Redirect(w, r, FacebookLoginFailed, http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := cli.getTokenFromFacebookCode(code)

		if err != nil {
			cli.Logger.Error(err.Error())
			cli.Logger.Error("typically caused by user cancelled verification")
			http.Redirect(w, r, FacebookLoginFailed, http.StatusTemporaryRedirect)
			return
		}

		data, err := cli.getFacebookAPIData(token.AccessToken)

		response, err := json.Marshal(data)
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		// Set token in cache
		cli.Cache.Set(cli.AuthConfig.GetAuthState().OauthStateString, response, cache.DefaultExpiration)
		http.Redirect(w, r, FacebookLoginSuccess, http.StatusTemporaryRedirect)
	})
}

func (cli *FacebookClient) checkOauthState(state string) error {
	oauthState := cli.AuthConfig.GetAuthState().OauthStateString
	if state == "" || oauthState == "" {
		return errors.New("missing facebook oauth state")
	}
	if state != oauthState {
		return errors.New("facebook oauth state does not match client")
	}
	if state != oauthState {
		return errors.New("facebook oauth state does not match client")
	}
	return nil
}

func (cli *FacebookClient) getTokenFromFacebookCode(code string) (*oauth2.Token, error) {
	token, err := cli.AuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (cli *FacebookClient) getFacebookAPIData(accessToken string) (*FacebookResponse, error) {
	resp, err := cli.AuthConfig.GetAPIData(accessToken)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &FacebookResponse{}
	if err = json.Unmarshal(body, data); err != nil {
		return nil, err
	}
	data.Token = accessToken
	return data, nil
}

func (cli *FacebookClient) systemErrorRedirect(w http.ResponseWriter, r *http.Request, err error) {
	cli.Logger.ErrorStackTrace("Facebook Login Error: ", errors.Wrap(err, 1))
	http.Redirect(w, r, FacebookLoginFailed, http.StatusTemporaryRedirect)
}
