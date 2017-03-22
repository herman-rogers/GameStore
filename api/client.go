package api

import (
	"net/http"
	"time"

	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
	"github.com/WestCoastOpenSource/GameStore/pkg/storage"
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
	saveSystem := storage.LocalDisk{
		File:      "gamestore_client_",
		Directory: "/var/log/gamestore/",
	}
	logSystem := logger.FileSystemLogs{Save: &saveSystem}
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

	// API Routes
	cli.Handler.Handle(ServerStatus, cli.serverStatusHandler())
	cli.Handler.Handle(FacebookCallback, cli.handleFacebookCallback())
	cli.Handler.Handle(FacebookLogin, cli.handleFacebookLogin())
	cli.Handler.Handle(GetFacebookData, sha.Handler(cli.getFacebookData()))
	cli.Handler.Handle(DeleteFacebookData, sha.Handler(cli.deleteFacebookData()))

	// Template Routes
	cli.Handler.Handle(FacebookLoginSuccess, cli.handleFacebookLoginSuccess())
	cli.Handler.Handle(FacebookLoginFailed, cli.handleFacebookLoginFailed())

	// Resources
	cli.Handler.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
}
