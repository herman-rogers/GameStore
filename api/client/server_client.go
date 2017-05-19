package client

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/WestCoastOpenSource/GameStore/pkg/clock"
	"github.com/WestCoastOpenSource/GameStore/pkg/logger"
)

// ServerClient is the client for interacting with the server
type ServerClient struct {
	Logger logger.Log
	Clock  clock.TimeInfo
}

// ServerPingResp is the api response from the server client
type ServerPingResp struct {
	Status string `json:"status"`
}

// ServerTimeResp sends the servers time in UTC to the client
type ServerTimeResp struct {
	Time string `json:"time"`
}

// HandleServerStatus gives information about the servers health
func (cli ServerClient) HandleServerStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(ServerPingResp{Status: "Server OK"})
		if err != nil {
			cli.Logger.Error(err.Error())
			return
		}
		w.Write(resp)
	})
}

// HandleServerTime returns the server time in UTC
func (cli ServerClient) HandleServerTime() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		serverTime := cli.Clock.ClockTime(time.Now().UTC())

		resp, err := json.Marshal(ServerTimeResp{Time: serverTime})
		if err != nil {
			cli.Logger.Error(err.Error())
			return
		}
		w.Write(resp)
	})
}
