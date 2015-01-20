package providers

import (
	"github.com/taviti/caldav-go/caldav"
	"github.com/taviti/caldav-go/webdav"
	"github.com/taviti/caldav-go/webdav/entities"
	. "github.com/taviti/check"
	"os"
	"testing"
)

type CalendarServerSuite struct{ client caldav.Client }

var _ = Suite(new(CalendarServerSuite))

func Test(t *testing.T) { TestingT(t) }

func (s *CalendarServerSuite) SetUpSuite(c *C) {
	rawurl := os.Getenv("CALDAV_SERVER_URL")
	c.Assert(rawurl, Not(HasLen), 0)
	p, err := NewCalendarServerProvider(rawurl)
	c.Assert(err, IsNil)
	s.client = caldav.NewDefaultClient(p)
}

// tests the current server for CalDAV support
func (s *CalendarServerSuite) TestValidate(c *C) {
	err := s.client.Validate("/")
	c.Assert(err, IsNil)
}

// tests the current server for a calendar collection
func (s *CalendarServerSuite) TestPropfind(c *C) {
	username := s.client.Provider().CurrentUsername()
	path := s.client.Provider().UserCalendarPath(username)
	ms, err := s.client.Propfind(path, webdav.Depth0, entities.AllProps())
	c.Assert(err, IsNil)
	c.Assert(ms.Responses, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats[0].Prop, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType.Calendar, NotNil)
}
