package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
	"github.com/go-errors/errors"
	"github.com/patrickmn/go-cache"
)

// Client creates a new web api client to recieve requests
type Client struct {
	Handler *http.ServeMux
	Cache   *cache.Cache
	Logger  logger.Log
}

// Response is a generic struct to format and return jsonified data
type Response struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	ID     string `json:"id"`
	Token  string `json:"token"`
}

const facebookAccessURL string = "https://graph.facebook.com/me?access_token="

// Start creates and returns a new server Client
func Start() *Client {
	logSystem := logger.FileSystemLogs{
		File:      "gamestore_client_",
		Directory: "/var/log/gamestore/",
	}
	client := Client{
		Handler: http.NewServeMux(),
		Cache:   cache.New(60*time.Minute, 30*time.Second),
		Logger:  &logSystem,
	}

	client.addRoutes()
	return &client
}

func (cli Client) addRoutes() {
	sha := auth.SHA256Middleware{}
	cli.Logger.Error(errors.Errorf("Test Error Format"))
	cli.Handler.Handle(ServerStatus, cli.serverStatusHandler())
	cli.Handler.Handle(FacebookCallback, cli.handleFacebookCallback())
	cli.Handler.Handle(FacebookLogin, cli.handleFacebookLogin())
	cli.Handler.Handle(FacebookLoginPage, cli.handleFacebookLoginPage())
	cli.Handler.Handle(GetFacebookData, sha.Handler(cli.getFacebookData()))
	cli.Handler.Handle(DeleteFacebookData, sha.Handler(cli.deleteFacebookData()))
}

func (cli Client) serverStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(Response{Status: "Server OK"})
		if err != nil {
			cli.Logger.Error(err)
			return
		}
		w.Write(resp)
	})
}
