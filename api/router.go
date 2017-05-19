package api

import (
	"net/http"

	"github.com/WestCoastOpenSource/GameStore/api/client"
	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
)

// GameStoreRouter is the main route handler for GameStore
type GameStoreRouter struct {
	HTTPHandler    *http.ServeMux
	FacebookClient *client.FacebookClient
	ServerClient   *client.ServerClient
}

// CreateRoutes will create all client api endpoints
func (router GameStoreRouter) CreateRoutes() {
	router.facebookRoutes()
	router.resourceRoutes()
	router.serverRoutes()
}

func (router GameStoreRouter) facebookRoutes() {
	sha := auth.SHA256Middleware{}

	// Facebook Templates
	router.HTTPHandler.Handle(client.FacebookLoginSuccess, router.FacebookClient.HandleFacebookLoginSuccess())
	router.HTTPHandler.Handle(client.FacebookLoginFailed, router.FacebookClient.HandleFacebookLoginFailed())

	// API Routes
	router.HTTPHandler.Handle(client.FacebookCallback, router.FacebookClient.HandleFacebookCallback())
	router.HTTPHandler.Handle(client.FacebookLogin, router.FacebookClient.HandleFacebookLogin())
	router.HTTPHandler.Handle(client.GetFacebookData, sha.Handler(router.FacebookClient.GetFacebookData()))
	router.HTTPHandler.Handle(client.DeleteFacebookData, sha.Handler(router.FacebookClient.DeleteFacebookData()))
}

func (router GameStoreRouter) serverRoutes() {
	sha := auth.SHA256Middleware{}
	router.HTTPHandler.Handle(client.ServerStatus, router.ServerClient.HandleServerStatus())
	router.HTTPHandler.Handle(client.ServerTime, sha.Handler(router.ServerClient.HandleServerTime()))
}

func (router GameStoreRouter) resourceRoutes() {
	router.HTTPHandler.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
}
