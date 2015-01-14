package caldav

import (
	"fmt"
	"net/http"
	"strings"
)

// a standards compliant interface for talking over CalDAV
type client struct {
	Provider
	*http.Client
}

// executes a HTTP request using the inner http-client
func (c *client) do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

// verifies that the CalDAV provider can, in fact, speak CalDAV
func (c *client) Verify(capabilities ...string) error {

	var dav string

	if r, err := c.do(c.Provider.NewVerifyRequest()); err != nil {
		return fmt.Errorf("unable to conntect to the CalDAV provider, %s", err)
	} else if r.StatusCode != 200 {
		return fmt.Errorf("invalid status from CalDAV server: %s", r.Status)
	} else if dav = r.Header.Get("DAV"); dav == "" {
		return fmt.Errorf("server does not support WebDAV: missing DAV header")
	}

	capabilities = append(capabilities, "calendar-access")
	for _, capability := range capabilities {
		if !strings.Contains(dav, capability) {
			return fmt.Errorf("server does not support required CalDAV capability: %s", capability)
		}
	}

	return nil

}

// creates a new client for communicating with a CalDAV server
func NewClient(provider Provider, httpClient *http.Client) Client {
	c := new(client)
	c.Provider = provider
	c.Client = httpClient
	return c
}

// creates a new client for communicating with a CalDAV server
// uses the default HTTP client from net/http
func NewDefaultClient(provider Provider) Client {
	return NewClient(provider, http.DefaultClient)
}
