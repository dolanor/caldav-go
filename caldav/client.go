package caldav

import (
	"fmt"
	cent "github.com/taviti/caldav-go/caldav/entities"
	"github.com/taviti/caldav-go/icalendar/components"
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav"
	"github.com/taviti/caldav-go/webdav/entities"
	"log"
	"net/http"
	"strings"
)

var _ = log.Print

// a client for making WebDAV requests
type Client webdav.Client

// downcasts the client to the WebDAV interface
func (c *Client) WebDAV() *webdav.Client {
	return (*webdav.Client)(c)
}

// returns the embedded CalDAV server reference
func (c *Client) Server() *Server {
	return (*Server)(c.WebDAV().Server())
}

// fetches a list of CalDAV features supported by the server
// returns an error if the server does not support DAV
func (c *Client) Features() ([]string, error) {
	var cfeatures []string
	if features, err := c.WebDAV().Features(); err != nil {
		return cfeatures, utils.NewError(c.Features, "unable to detect features", c, err)
	} else {
		for _, feature := range features {
			if strings.HasPrefix(feature, "calendar-") {
				cfeatures = append(cfeatures, feature)
			}
		}
		return cfeatures, nil
	}
}

// fetches a list of CalDAV features and checks if a certain one is supported by the server
// returns an error if the server does not support DAV
func (c *Client) SupportsFeature(name string) (bool, error) {
	if features, err := c.Features(); err != nil {
		return false, utils.NewError(c.SupportsFeature, "feature detection failed", c, err)
	} else {
		var test = fmt.Sprintf("calendar-%s", name)
		for _, feature := range features {
			if feature == test {
				return true, nil
			}
		}
		return false, nil
	}
}

// fetches a list of CalDAV features and checks if a certain one is supported by the server
// returns an error if the server does not support DAV
func (c *Client) ValidateServer() error {
	if found, err := c.SupportsFeature("access"); err != nil {
		return utils.NewError(c.SupportsFeature, "feature detection failed", c, err)
	} else if !found {
		return utils.NewError(c.SupportsFeature, "calendar access feature missing", c, nil)
	} else {
		return nil
	}
}

// creates or updates an event on the remote CalDAV server
func (c *Client) PutEvents(path string, events ...*components.Event) error {
	if len(events) <= 0 {
		return utils.NewError(c.PutEvents, "no calendar events provided", c, nil)
	} else if cal := components.NewCalendar(events...); events[0] == nil {
		return utils.NewError(c.PutEvents, "icalendar event must not be nil", c, nil)
	} else if err := cal.ValidateICalValue(); err != nil {
		return utils.NewError(c.PutEvents, "invalid icalendar event", c, err)
	} else if req, err := c.Server().NewRequest("PUT", path, cal); err != nil {
		return utils.NewError(c.PutEvents, "unable to encode request", c, err)
	} else if resp, err := c.Do(req); err != nil {
		return utils.NewError(c.PutEvents, "unable to execute request", c, err)
	} else if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		err := new(entities.Error)
		resp.WebDAV().Decode(err)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		return utils.NewError(c.PutEvents, msg, c, err)
	}
	return nil
}

// attempts to fetch an event on the remote CalDAV server
func (c *Client) GetEvents(path string) ([]*components.Event, error) {
	cal := new(components.Calendar)
	if req, err := c.Server().NewRequest("GET", path); err != nil {
		return nil, utils.NewError(c.GetEvents, "unable to create request", c, err)
	} else if resp, err := c.Do(req); err != nil {
		return nil, utils.NewError(c.GetEvents, "unable to execute request", c, err)
	} else if resp.StatusCode != http.StatusOK {
		err := new(entities.Error)
		resp.WebDAV().Decode(err)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		return nil, utils.NewError(c.GetEvents, msg, c, err)
	} else if err := resp.Decode(cal); err != nil {
		return nil, utils.NewError(c.GetEvents, "unable to decode response", c, err)
	} else {
		return cal.Events, nil
	}
}

// attempts to fetch an event on the remote CalDAV server
func (c *Client) QueryEvents(path string, query *cent.CalendarQuery) (events []*components.Event, oerr error) {
	ms := new(cent.Multistatus)
	if req, err := c.Server().WebDAV().NewRequest("REPORT", path, query); err != nil {
		oerr = utils.NewError(c.QueryEvents, "unable to create request", c, err)
	} else if resp, err := c.WebDAV().Do(req); err != nil {
		oerr = utils.NewError(c.QueryEvents, "unable to execute request", c, err)
	} else if resp.StatusCode != webdav.StatusMulti {
		err := new(entities.Error)
		msg := fmt.Sprintf("unexpected server response %s", resp.Status)
		resp.Decode(err)
		oerr = utils.NewError(c.QueryEvents, msg, c, err)
	} else if err := resp.Decode(ms); err != nil {
		msg := "unable to decode response"
		oerr = utils.NewError(c.QueryEvents, msg, c, err)
	} else {
		for i, r := range ms.Responses {
			for j, p := range r.PropStats {
				if p.Prop == nil || p.Prop.CalendarData == nil {
					continue
				} else if cal, err := p.Prop.CalendarData.CalendarComponent(); err != nil {
					msg := fmt.Sprintf("unable to decode property %d of response %d", j, i)
					oerr = utils.NewError(c.QueryEvents, msg, c, err)
					return
				} else {
					events = append(events, cal.Events...)
				}
			}
		}
	}
	return
}

// executes a CalDAV request
func (c *Client) Do(req *Request) (*Response, error) {
	if resp, err := c.WebDAV().Do((*webdav.Request)(req)); err != nil {
		return nil, utils.NewError(c.Do, "unable to execute CalDAV request", c, err)
	} else {
		return NewResponse(resp), nil
	}
}

// creates a new client for communicating with an WebDAV server
func NewClient(server *Server, native *http.Client) *Client {
	return (*Client)(webdav.NewClient((*webdav.Server)(server), native))
}

// creates a new client for communicating with a WebDAV server
// uses the default HTTP client from net/http
func NewDefaultClient(server *Server) *Client {
	return NewClient(server, http.DefaultClient)
}
