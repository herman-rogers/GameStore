package api

import (
	"encoding/json"
	"net/http"
)

func (cli Client) serverStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(Response{Status: "Server OK"})
		if err != nil {
			cli.Logger.Error(err.Error())
			return
		}
		w.Write(resp)
	})
}
