package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/WestCoastOpenSource/GameStore/pkg/testutils"
)

func TestServerStatusReturnsOK(t *testing.T) {
	cli := Start()

	r, err := http.NewRequest("GET", ServerStatus, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	testutils.AuthenticateRoute(r)

	w := httptest.NewRecorder()
	cli.Handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}

	resp := Response{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf(err.Error())
	}
	if resp.Status != "Server OK" {
		t.Errorf("Expected %v, but got %v", "Status OK", resp.Status)
	}
}

// func TestServerReturnsUnauthorizedAccess(t *testing.T) {
// 	cli := Start()

// 	r, err := http.NewRequest("GET", ServerStatus, nil)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	w := httptest.NewRecorder()
// 	cli.Handler.ServeHTTP(w, r)

// 	if w.Code != http.StatusUnauthorized {
// 		t.Errorf("Expected %v, but got %v", http.StatusUnauthorized, w.Code)
// 	}
// }
