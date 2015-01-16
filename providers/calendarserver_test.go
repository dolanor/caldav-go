package providers

import (
	. "github.com/go-check/check"
	"github.com/taviti/caldav-go/caldav"
	"github.com/taviti/caldav-go/webdav/constants"
	"github.com/taviti/caldav-go/webdav/entities"
	"os"
	"testing"
)

type suite struct{ client caldav.Client }

var _ = Suite(new(suite))

func Test(t *testing.T) { TestingT(t) }

func (s *suite) SetUpSuite(c *C) {
	rawurl := os.Getenv("CALDAV_SERVER_URL")
	c.Assert(rawurl, Not(HasLen), 0)
	p, err := NewCalendarServerProvider(rawurl)
	c.Assert(err, Not(NotNil))
	s.client = caldav.NewDefaultClient(p)
}

// tests the current server for CalDAV support
func (s *suite) TestValidate(c *C) {
	err := s.client.Validate("/")
	c.Assert(err, Not(NotNil))
}

// tests the current server for a calendar collection
func (s *suite) TestPropfind(c *C) {
	username := s.client.Provider().CurrentUsername()
	path := s.client.Provider().UserCalendarPath(username)
	ms, err := s.client.Propfind(path, constants.Depth0, entities.AllProps())
	c.Assert(err, Not(NotNil))
	c.Assert(ms.Responses, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats[0].Prop, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType.Calendar, NotNil)
}
