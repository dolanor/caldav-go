package caldav

import (
	"fmt"
	calentities "github.com/taviti/caldav-go/caldav/entities"
	"github.com/taviti/caldav-go/icalendar/components"
	"github.com/taviti/caldav-go/icalendar/properties"
	"github.com/taviti/caldav-go/icalendar/values"
	"github.com/taviti/caldav-go/utils"
	"github.com/taviti/caldav-go/webdav"
	webentities "github.com/taviti/caldav-go/webdav/entities"
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
	ms, err := s.client.WebDAV().Propfind(s.path, webdav.Depth0, webentities.NewAllPropsFind())
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
	uuid := fmt.Sprintf("test-single-event-%d", oneHourFromNow.Unix())
	putEvent := components.NewEventWithDuration(uuid, oneHourFromNow, time.Hour)
	putEvent.Summary = "This is a test single event"

	// generate an ICS filepath
	path := fmt.Sprintf("%s%s.ics", s.path, uuid)

	// save the event to the server, then fetch it back out
	if err = s.client.PutEvents(path, putEvent); err != nil {
		c.Fatal(err.Error())
	} else if getEvents, err := s.client.GetEvents(path); err != nil {
		c.Fatal(err.Error())
	} else {
		// assert that the events match
		c.Assert(getEvents, HasLen, 1)
		c.Assert(getEvents[0], DeepEquals, putEvent)
	}

}

func (s *ClientSuite) TestRecurringEventQuery(c *C) {

	// create the event object
	start := time.Now().Truncate(time.Hour).UTC()
	uuid := fmt.Sprintf("test-recurring-event-%d", start.Unix())
	putEvent := components.NewEventWithDuration(uuid, start, time.Hour)
	putEvent.Summary = "This is a test recurring event"
	putEvent.RecurrenceRule = values.NewRecurrenceRule(values.DayRecurrenceFrequency)
	putEvent.RecurrenceRule.Count = 14 // two weeks of events

	// generate an ICS filepath
	path := fmt.Sprintf("%s%s.ics", s.path, uuid)

	// save the event to the server
	if err := s.client.PutEvents(path, putEvent); err != nil {
		c.Fatal(err.Error())
	}

	// now create a query for one week out
	days := 2
	nextWeek := start.AddDate(0, 0, 7)
	nextWeekEnd := nextWeek.AddDate(0, 0, days)
	query, err := calentities.NewEventRangeQuery(nextWeek, nextWeekEnd)
	if err != nil {
		c.Fatal(err.Error())
	}

	// add in a filter for UID so that we don't get back unwanted results
	pf := calentities.NewPropertyMatcher(properties.UIDPropertyName, uuid)
	query.Filter.ComponentFilter.ComponentFilter.PropertyFilter = pf

	if events, err := s.client.QueryEvents(s.path, query); err != nil {
		c.Fatal(err.Error())
	} else {
		c.Assert(events, HasLen, days)
		for i, event := range events {
			day := nextWeek.AddDate(0, 0, i)
			offset := values.NewDateTime(day)
			c.Assert(event.UID, Equals, uuid)
			c.Assert(event.DateStart, DeepEquals, offset)
			c.Assert(event.RecurrenceId, DeepEquals, offset)
		}
	}

}
