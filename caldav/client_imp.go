package caldav

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/taviti/caldav-go/webdav"
	"github.com/taviti/caldav-go/webdav/entities"
	"io/ioutil"
	"net/http"
	"strings"
)

// a standards compliant interface for talking over CalDAV
type client struct {
	provider   Provider
	httpClient *http.Client
}

func (c *client) newRequest(method string, path string, xmlData ...interface{}) (r *http.Request, err error) {
	r = new(http.Request)
	r.Method = method
	r.URL = c.provider.AbsURL(path)
	r.Proto = "HTTP/1.1"
	r.ProtoMajor = 1
	r.ProtoMinor = 1
	r.Header = make(http.Header)
	r.Header.Set("Content-Type", "text/xml; charset=UTF-8")
	r.Host = r.URL.Host
	if len(xmlData) > 0 && xmlData[0] != nil {
		if buffer, err := xml.Marshal(xmlData); err != nil {
			return nil, fmt.Errorf("unable to encode %s request, %s", method, err)
		} else {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
		}
	}
	return
}

func (c *client) do(r *http.Request, expectStatus int, xmlData ...interface{}) (resp *http.Response, err error) {
	if resp, err = c.httpClient.Do(r); err != nil {
		return nil, fmt.Errorf("unable to connect to the CalDAV provider, %s", err)
	} else if expectStatus > 0 && resp.StatusCode != expectStatus {
		return nil, fmt.Errorf("unexpected status %s returned from provider, expected %d", resp.Status, expectStatus)
	} else if len(xmlData) > 0 && xmlData[0] != nil {
		decoder := xml.NewDecoder(resp.Body)
		if err = decoder.Decode(xmlData[0]); err != nil {
			return nil, fmt.Errorf("unable to decode response, %s", err)
		}
	}
	return
}

func (c *client) Provider() Provider {
	return c.provider
}

func (c *client) Validate(path string, capabilities ...string) error {

	var dav string
	var req, _ = c.newRequest("OPTIONS", path)

	if resp, err := c.do(req, 200); err != nil {
		return err
	} else if dav = resp.Header.Get("DAV"); dav == "" {
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

func (c *client) Propfind(path string, depth webdav.Depth, pf *entities.Propfind) (*entities.Multistatus, error) {

	req, err := c.newRequest(webdav.PropfindMethod, path, pf)
	if err == nil {
		req.Header.Set("Depth", string(depth))
	} else {
		return nil, err
	}

	ms := new(entities.Multistatus)
	if _, err := c.do(req, 207, ms); err != nil {
		return nil, err
	}

	return ms, nil

}

// creates a new client for communicating with a CalDAV server
func NewClient(provider Provider, httpClient *http.Client) Client {
	c := new(client)
	c.provider = provider
	c.httpClient = httpClient
	return c
}

// creates a new client for communicating with a CalDAV server
// uses the default HTTP client from net/http
func NewDefaultClient(provider Provider) Client {
	return NewClient(provider, http.DefaultClient)
}
