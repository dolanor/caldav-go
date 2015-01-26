package caldav

import (
	"fmt"
	"github.com/taviti/caldav-go/icalendar/components"
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav"
	"github.com/taviti/caldav-go/webdav/entities"
	. "github.com/taviti/check"
	"testing"
	"time"
)

type ClientSuite struct {
	client *Client
	server *Server
	path   string
}

var _ = Suite(new(ClientSuite))

func Test(t *testing.T) { TestingT(t) }

func (s *ClientSuite) SetUpSuite(c *C) {
	var err error
	uri := utils.AssertServerUrl(c)
	s.path = uri.Path
	s.server, err = NewServer(uri.String())
	c.Assert(err, IsNil)
	s.client = NewDefaultClient(s.server)
}

func (s *ClientSuite) TestValidate(c *C) {
	c.Assert(s.client.ValidateServer(), IsNil)
}

func (s *ClientSuite) TestPropfind(c *C) {
	ms, err := s.client.WebDAV().Propfind(s.path, webdav.Depth0, entities.AllProps())
	c.Assert(err, IsNil)
	c.Assert(ms.Responses, Not(HasLen), 0)
	c.Assert(ms.Responses[0].Href, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats[0].Prop, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType.Calendar, NotNil)
}

func (s *ClientSuite) TestPutEventBasic(c *C) {

	// select a timezone
	loc, err := time.LoadLocation("America/New_York")
	c.Assert(err, IsNil)

	// create the event object
	oneHourFromNow := time.Now().Add(time.Hour).Truncate(time.Hour).In(loc)
	uuid := fmt.Sprintf("test-event-%d", oneHourFromNow.Unix())
	event := components.NewEventWithDuration(uuid, oneHourFromNow, time.Hour)
	event.Summary = "This is a test event"

	// generate an ICS filepath
	path := fmt.Sprintf("%s%s.ics", s.path, uuid)

	// ship it!
	if err = s.client.PutEvent(path, event); err != nil {
		c.Fatal(err.Error())
	}

}
