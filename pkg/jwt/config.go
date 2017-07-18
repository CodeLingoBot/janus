package jwt

import (
	"net/http"
	"time"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/response"
)

// Config struct
type Config struct {
	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userID string, password string) (string, bool)

	// User can define own Unauthorized func.
	Unauthorized func(w http.ResponseWriter, r *http.Request, err error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, w http.ResponseWriter, r *http.Request) bool

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// SigningAlgorithm defines token signing algorithm.
	// Possible values are: HS256, HS384, HS512, RS256, RS384, RS512
	SigningAlgorithm string

	// SigningKey key used for signing. Required.
	SigningKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration
}

// NewConfig creates a new instance of Config
func NewConfig(algorithm, key string) Config {
	return Config{
		SigningAlgorithm: algorithm,
		SigningKey:       key,
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		TokenLookup:      "header:Authorization",
	}
}

// NewConfigWithHandlers creates a new instance of Config with default handlers
func NewConfigWithHandlers(cred config.Credentials) Config {
	jwtConfig := NewConfig(cred.Algorithm, cred.Secret)
	jwtConfig.Authenticator = func(userID string, password string) (string, bool) {
		if userID == cred.Username && password == cred.Password {
			return userID, true
		}

		return userID, false
	}
	jwtConfig.Authorizator = func(userID string, w http.ResponseWriter, r *http.Request) bool {
		return userID == cred.Username
	}

	jwtConfig.Unauthorized = func(w http.ResponseWriter, r *http.Request, err error) {
		response.JSON(w, http.StatusUnauthorized, response.H{
			"message": err.Error(),
		})
	}
	return jwtConfig
}
