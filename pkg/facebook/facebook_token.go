package facebook

import (
	"os"

	"github.com/WestCoastOpenSource/GameStore/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// AuthConfig sets the environment configs for Facebook auth
var AuthConfig = &oauth2.Config{
	ClientID:     os.Getenv(config.FBClientID),
	ClientSecret: os.Getenv(config.FBSecretKey),
	RedirectURL:  os.Getenv(config.FacebookRedirectURL) + "/api/facebookcallback",
	Scopes:       []string{"public_profile"},
	Endpoint:     facebook.Endpoint,
}

// OathStateString is the unique identifier from the client request
var OathStateString = ""
