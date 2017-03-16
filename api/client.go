package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Handler *http.ServeMux
}

type Response struct {
	Status string
}

// Start creates and returns a new server Client
func Start() Client {
	client := Client{Handler: http.NewServeMux()}
	client.addRoutes()

	return client
}

func (cli Client) addRoutes() {
	cli.Handler.Handle(ServerStatus, cli.serverStatusHandler())
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
