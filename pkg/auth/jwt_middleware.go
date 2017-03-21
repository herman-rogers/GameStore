package auth

import (
	"errors"
	"os"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

const missingSecretKey string = "Missing Auth0 Client Secret"

// JWTAuthMiddleware is the route middleware for jwt authentication
var JWTAuthMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv(config.ClientSecret))
		if len(secret) == 0 {
			return nil, errors.New(missingSecretKey)
		}
		return secret, nil
	},
})
