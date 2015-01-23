package webdav

import (
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav/entities"
	. "github.com/taviti/check"
	"testing"
)

type ClientSuite struct {
	client *Client
	server *Server
}

var _ = Suite(new(ClientSuite))

func Test(t *testing.T) { TestingT(t) }

func (s *ClientSuite) SetUpSuite(c *C) {
	var err error
	uri := utils.AssertServerUrl(c)
	s.server, err = NewServer(uri.String())
	c.Assert(err, IsNil)
	s.client = NewDefaultClient(s.server)
}

func (s *ClientSuite) TestValidate(c *C) {
	c.Assert(s.client.ValidateServer(), IsNil)
}

func (s *ClientSuite) TestPropfind(c *C) {
	ms, err := s.client.Propfind("/", Depth0, entities.AllProps())
	c.Assert(err, IsNil)
	c.Assert(ms.Responses, Not(HasLen), 0)
	c.Assert(ms.Responses[0].Href, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats, Not(HasLen), 0)
	c.Assert(ms.Responses[0].PropStats[0].Prop, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType, NotNil)
	c.Assert(ms.Responses[0].PropStats[0].Prop.ResourceType.Collection, NotNil)
}
