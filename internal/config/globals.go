package config

import (
	"fmt"
	"plex_monitor/internal/secrets"

	"github.com/go-chi/jwtauth/v5"
)

// Globals holds globals for the application. It is a singleton and should be accessed via the GetGlobals() function.
// It is initialized in the main function.
type Globals struct {
	// JWTAuth is the JWT authentication
	JWTAuth *jwtauth.JWTAuth
}

var globals *Globals

func init() {
	if globals == nil {
		populateGlobals()
	}
}

func populateGlobals() {
	secret, err := secrets.GetSecret("SECRET_KEY")
	if err != nil {
		panic(fmt.Errorf("failed to get secret key: %v", err))
	}

	SetGlobals(&Globals{
		JWTAuth: jwtauth.New("HS256", []byte(secret), nil),
	})
}

// GetGlobals returns the globals for the application
func GetGlobals() *Globals {
	if globals == nil {
		populateGlobals()
	}
	return globals
}

// SetGlobals sets the globals for the application
func SetGlobals(g *Globals) {
	globals = g
}
