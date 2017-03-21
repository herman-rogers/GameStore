package auth

import (
	"net/http"
)

// CredentialMiddleware is the interface for
// implementing multiple authentication types
type CredentialMiddleware interface {
	CreateKey(message string)
	VerifyKey(message string)
	Handler(h http.Handler) http.Handler
}
