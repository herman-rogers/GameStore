package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
)

func mockRouteEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestHMAC256CanWrapAndAuthRoutes(t *testing.T) {
	authHandler := SHA256Middleware{}
	mockRouter := http.NewServeMux()
	r, err := http.NewRequest("GET", "/mock/route", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	os.Setenv(config.ClientSecret, "secret")
	token := authHandler.CreateKey("GameStore Request")
	r.Header.Add("Authorization", token)

	w := httptest.NewRecorder()
	mockRouter.Handle("/mock/route", authHandler.Handler(mockRouteEndpoint()))
	mockRouter.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}
}

func TestHMAC256DeniesAccesWithWrongKey(t *testing.T) {
	authHandler := SHA256Middleware{}
	mockRouter := http.NewServeMux()
	r, err := http.NewRequest("GET", "/mock/route", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	os.Setenv(config.ClientSecret, "wrongkey")
	token := authHandler.CreateKey("GameStore Request")
	r.Header.Add("Authorization", token)

	w := httptest.NewRecorder()

	os.Setenv(config.ClientSecret, "secret")
	mockRouter.Handle("/mock/route", authHandler.Handler(mockRouteEndpoint()))
	mockRouter.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %v, but got %v", http.StatusUnauthorized, w.Code)
	}
}
