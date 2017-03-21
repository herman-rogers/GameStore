package testutils

import (
	"net/http"
	"os"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
	jwt "github.com/dgrijalva/jwt-go"
)

// GetJWTToken setups an test environment for mocking auth api calls
func GetJWTToken(tokenValue string) (string, error) {
	os.Setenv(config.ClientSecret, "secret")

	secretKey := []byte(tokenValue)
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AuthenticateJWTRoute will setup a route with token authentification
func AuthenticateJWTRoute(r *http.Request) error {
	token, err := GetJWTToken("secret")
	if err != nil {
		return err
	}
	bearer := "Bearer " + token
	r.Header.Add("Authorization", bearer)
	return nil
}
