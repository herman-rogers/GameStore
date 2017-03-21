package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"

	"fmt"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
)

// SHA256Middleware is SHA256 auth middleware to wrap routes and auth each request
type SHA256Middleware struct{}

// Handler creates a http.Handler to authenticate traffic requests
func (sha *SHA256Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let secure process the request. If it returns false,
		// that indicates the request should not continue.
		authenticated := sha.VerifyKey(r)
		if !authenticated {
			fmt.Printf("Missing Authentication Token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// CreateKey creates a new sha256 key from a secretkey and message
func (sha *SHA256Middleware) CreateKey(message string) string {
	key := []byte(os.Getenv(config.ClientSecret))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifyKey compares an generated expected key with a key recieved from a request
func (sha *SHA256Middleware) VerifyKey(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Printf("HEADER MISSING")
		return false
	}
	decodedKey, _ := base64.StdEncoding.DecodeString(authHeader)

	key := []byte(os.Getenv(config.ClientSecret))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("GameStore Request"))
	expectedMac := mac.Sum(nil)

	return hmac.Equal(decodedKey, expectedMac)
}
