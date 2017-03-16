package testutils

import (
	"net/http"
	"os"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
	jwt "github.com/dgrijalva/jwt-go"
)

// GetMockToken setups an test environment for mocking auth api calls
func GetMockToken(tokenValue string) (string, error) {
	os.Setenv(config.JWTClientSecret, "secret")

	secretKey := []byte(tokenValue)
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AuthenticateRoute will setup a route with token authentification
func AuthenticateRoute(r *http.Request) error {
	token, err := GetMockToken("secret")
	if err != nil {
		return err
	}
	bearer := "Bearer " + token
	r.Header.Add("Authorization", bearer)
	return nil
}
