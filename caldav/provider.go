package caldav

import (
	"io"
	"net/http"
)

// A mechanism for defining, authenticating and connecting to a CalDav server
type Provider interface {

	// creates a generic request for the CalDAV provider
	NewRequest(method, path string, body io.ReadCloser) *http.Request

	// creates a request to verify the provider as a CalDAV server
	NewVerifyRequest() *http.Request
}
