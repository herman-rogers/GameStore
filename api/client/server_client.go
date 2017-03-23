package client

import (
	"encoding/json"
	"net/http"

	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
)

// ServerClient is the client for interacting with the server
type ServerClient struct {
	Logger logger.Log
}

// ServerResponse is the api response from the server client
type ServerResponse struct {
	Status string `json:"status"`
}

// HandleServerStatus gives information about the servers health
func (cli ServerClient) HandleServerStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(ServerResponse{Status: "Server OK"})
		if err != nil {
			cli.Logger.Error(err.Error())
			return
		}
		w.Write(resp)
	})
}
