package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/WestCoastOpenSource/GameStore/pkg/facebook"
	"github.com/go-errors/errors"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

func (cli Client) handleFacebookLoginSuccess() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./templates/login_success.html")
	})
}

func (cli Client) handleFacebookLoginFailed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "./templates/login_failure.html")
	})
}

// should be called with RPC value of ?token=TOKEN_DATA
func (cli Client) handleFacebookLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := facebook.AuthConfig
		facebook.OathStateString = r.FormValue("token")

		URL, err := url.Parse(config.Endpoint.AuthURL)
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		parameters := url.Values{}
		parameters.Add("client_id", config.ClientID)
		parameters.Add("scope", strings.Join(config.Scopes, " "))
		parameters.Add("redirect_uri", config.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", facebook.OathStateString)
		URL.RawQuery = parameters.Encode()

		stringURL := URL.String()

		http.Redirect(w, r, stringURL, http.StatusTemporaryRedirect)
	})
}

func (cli Client) handleFacebookCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")

		if state != facebook.OathStateString {
			cli.systemErrorRedirect(w, r, errors.Errorf("Invalid Oath State String: "))
			return
		}

		code := r.FormValue("code")
		token, err := facebook.AuthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			cli.Logger.Error("User Cancelled Facebook Verification. Callback Data: " + err.Error())
			http.Redirect(w, r, FacebookLoginFailed, http.StatusTemporaryRedirect)
			return
		}

		resp, err := http.Get(facebookAccessURL + url.QueryEscape(token.AccessToken))
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		data := &Response{}
		if err = json.Unmarshal(body, data); err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		data.Token = token.AccessToken

		response, err := json.Marshal(data)
		if err != nil {
			cli.systemErrorRedirect(w, r, err)
			return
		}

		// Set token in cache
		cli.Cache.Set(facebook.OathStateString, response, cache.DefaultExpiration)
		http.Redirect(w, r, FacebookLoginSuccess, http.StatusTemporaryRedirect)
	})
}

// endpoint called by clients to get their auth tokens in the format ?identity=IDENTITY
func (cli Client) getFacebookData() http.Handler {
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

func (cli Client) deleteFacebookData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity := r.FormValue("identity")
		cli.Cache.Delete(identity)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted"))
	})
}

func (cli Client) systemErrorRedirect(w http.ResponseWriter, r *http.Request, err error) {
	cli.Logger.ErrorStackTrace("Facebook Login Error: ", errors.Wrap(err, 1))
	http.Redirect(w, r, FacebookLoginFailed, http.StatusTemporaryRedirect)
}
