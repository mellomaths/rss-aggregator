package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts the API key from the HTTP Request Headers
// Example: Authorization: ApiKey <api_key>
func GetApiKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("no authorization header provided")
	}
	values := strings.Split(authorization, " ")
	if len(values) != 2 {
		return "", errors.New("malformed authorization header")
	}
	if values[0] != "ApiKey" {
		return "", errors.New("invalid authorization header format")
	}
	if len(values[1]) != 64 {
		return "", errors.New("invalid api key format")
	}
	return values[1], nil
}
