package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"encoding/json"
)

type mockClock struct{}

// ClockTime returns the server time in UTC formatted to YYYY-MM-DD HH:MM:SS
func (mockClock) ClockTime(serverTime time.Time) string {
	mockTime := "2006-01-02 03:04:05"
	return mockTime
}

func newServerClient() *ServerClient {
	server := ServerClient{}
	server.Clock = mockClock{}
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

	resp := ServerPingResp{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf(err.Error())
	}
	if resp.Status != "Server OK" {
		t.Errorf("Expected %v, but got %v", "Status OK", resp.Status)
	}
}

func TestServerReturnsUTCTime(t *testing.T) {
	expectedTime := "2006-01-02 03:04:05"
	router := http.NewServeMux()
	cli := newServerClient()

	router.Handle(ServerTime, cli.HandleServerTime())

	r, err := http.NewRequest("GET", ServerTime, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}

	resp := ServerTimeResp{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf(err.Error())
	}
	if resp.Time != expectedTime {
		t.Errorf("Expected %v, but got %v", expectedTime, resp.Time)
	}
}
