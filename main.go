package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WestCoastOpenSource/GameStore/api"
	"github.com/WestCoastOpenSource/GameStore/api/client"
	"github.com/WestCoastOpenSource/GameStore/pkg/auth"
	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
	"github.com/WestCoastOpenSource/GameStore/pkg/storage"
	cache "github.com/patrickmn/go-cache"
)

func main() {
	saveSystem := storage.LocalDisk{
		File:      "gamestore_client_",
		Directory: "/var/log/gamestore/",
	}
	logSystem := logger.FileSystemLogs{Save: &saveSystem}

	router := api.GameStoreRouter{
		HTTPHandler:    http.NewServeMux(),
		FacebookClient: newFacebookClient(logSystem),
		ServerClient:   newServerClient(logSystem),
	}
	router.CreateRoutes()

	fmt.Println("GameStore running on port :3000")
	if err := http.ListenAndServe(":3000", router.HTTPHandler); err != nil {
		fmt.Printf(err.Error())
	}
}

func newServerClient(logs logger.FileSystemLogs) *client.ServerClient {
	server := client.ServerClient{Logger: logs}
	return &server
}

func newFacebookClient(logs logger.FileSystemLogs) *client.FacebookClient {
	facebookConfig := auth.FacebookAuthConfig{}
	facebookConfig.Create()

	facebook := client.FacebookClient{
		AuthConfig:      &facebookConfig,
		Cache:           cache.New(60*time.Minute, 30*time.Second),
		Logger:          logs,
		SuccessTemplate: "./templates/login_success.html",
		FailureTemplate: "./templates/login_failure.html",
	}
	return &facebook
}
