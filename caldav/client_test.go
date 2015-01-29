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

func (s *ClientSuite) TestEventPutAndGet(c *C) {

	// select a timezone
	loc, err := time.LoadLocation("America/New_York")
	c.Assert(err, IsNil)

	// create the event object
	oneHourFromNow := time.Now().Add(time.Hour).Truncate(time.Hour).In(loc)
	uuid := fmt.Sprintf("test-event-%d", oneHourFromNow.Unix())
	putEvent := components.NewEventWithDuration(uuid, oneHourFromNow, time.Hour)
	putEvent.Summary = "This is a test event"

	// generate an ICS filepath
	path := fmt.Sprintf("%s%s.ics", s.path, uuid)

	// save the event to the server, then fetch it back out
	if err = s.client.PutEvent(path, putEvent); err != nil {
		c.Fatal(err.Error())
	} else if getEvent, err := s.client.GetEvent(path); err != nil {
		c.Fatal(err.Error())
	} else {
		// assert that the events match
		c.Assert(getEvent.UID, DeepEquals, putEvent.UID)
		c.Assert(getEvent.DateStart, DeepEquals, putEvent.DateStart)
		c.Assert(getEvent.DateStamp, DeepEquals, putEvent.DateStamp)
		c.Assert(getEvent.Summary, DeepEquals, putEvent.Summary)
		c.Assert(getEvent.Duration, DeepEquals, putEvent.Duration)
	}

}
