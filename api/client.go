package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"strings"

	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/facebook"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
)

type Client struct {
	Handler *http.ServeMux
	Cache   *cache.Cache
}

type Response struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

const facebookAccessURL string = "https://graph.facebook.com/me?access_token="

// Start creates and returns a new server Client
func Start() Client {
	client := Client{
		Handler: http.NewServeMux(),
		Cache:   cache.New(60*time.Minute, 30*time.Second),
	}

	client.addRoutes()
	return client
}

func (cli Client) addRoutes() {
	cli.Handler.Handle(ServerStatus, cli.serverStatusHandler())
	cli.Handler.Handle(FacebookLogin, cli.handleFacebookLogin())
	cli.Handler.Handle(FacebookCallback, cli.handleFacebookCallback())
	cli.Handler.Handle(FacebookLoginPage, cli.handleFacebookLoginPage())
	cli.Handler.Handle(GetFacebookData, auth.JWTAuthMiddleware.Handler(cli.getFacebookData()))
}

func (cli Client) serverStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(Response{Status: "Server OK"})
		if err != nil {
			fmt.Printf(err.Error())
		}
		w.Write(resp)
	})
}

// Remove and Replace
const htmlIndex = `<html><body>
Logged in with <a href="/login">facebook</a>
</body></html>
`

func (cli Client) handleFacebookLoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(htmlIndex))
		if err != nil {
			fmt.Printf(err.Error())
		}
	})
}

// should be called with RPC value of ?token=TOKEN_DATA
func (cli Client) handleFacebookLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := facebook.AuthConfig
		facebook.OathStateString = r.FormValue("token")

		URL, err := url.Parse(config.Endpoint.AuthURL)
		if err != nil {
			log.Fatal("Facebook Login Parse URL: ", err)
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
			cli.errorRedirect(w, r, "Invalid Auth State")
		}

		code := r.FormValue("code")
		token, err := facebook.AuthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			cli.errorRedirect(w, r, err.Error())
			return
		}

		resp, err := http.Get(facebookAccessURL + url.QueryEscape(token.AccessToken))
		if err != nil {
			cli.errorRedirect(w, r, err.Error())
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			cli.errorRedirect(w, r, err.Error())
			return
		}

		data := &Response{}
		if err = json.Unmarshal(body, data); err != nil {
			cli.errorRedirect(w, r, err.Error())
			return
		}

		data.Token = token.AccessToken

		response, err := json.Marshal(data)
		if err != nil {
			cli.errorRedirect(w, r, err.Error())
			return
		}

		// Set token in cache
		cli.Cache.Set(facebook.OathStateString, response, cache.DefaultExpiration)
		http.Redirect(w, r, FacebookLoginPage, http.StatusTemporaryRedirect)
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

func (cli Client) errorRedirect(w http.ResponseWriter, r *http.Request, err string) {
	fmt.Printf("Facebook Callback Error: %v", err)
	http.Redirect(w, r, FacebookLoginPage, http.StatusTemporaryRedirect)
}
