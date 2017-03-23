package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
)

func newServerClient() *ServerClient {
	server := ServerClient{}
	return &server
}

func TestServerStatusReturnsOK(t *testing.T) {
	router := http.NewServeMux()
	cli := newServerClient()

	router.Handle(ServerStatus, cli.HandleServerStatus())

	r, err := http.NewRequest("GET", ServerStatus, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}

	resp := ServerResponse{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf(err.Error())
	}
	if resp.Status != "Server OK" {
		t.Errorf("Expected %v, but got %v", "Status OK", resp.Status)
	}
}
