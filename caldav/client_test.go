package caldav

import (
	"fmt"
	calentities "github.com/dolanor/caldav-go/caldav/entities"
	"github.com/dolanor/caldav-go/icalendar/components"
	"github.com/dolanor/caldav-go/icalendar/properties"
	"github.com/dolanor/caldav-go/icalendar/values"
	"github.com/dolanor/caldav-go/webdav"
	webentities "github.com/dolanor/caldav-go/webdav/entities"
	. "gopkg.in/check.v1"
	"net/url"
	"os"
	"testing"
	"time"
)

type ClientSuite struct {
	client *Client
	server *Server
}

var _ = Suite(new(ClientSuite))

func Test(t *testing.T) { TestingT(t) }

func (s *ClientSuite) SetUpSuite(c *C) {
	var err error
	uri := AssertServerUrl(c)
	s.server, err = NewServer(uri.String())
	c.Assert(err, IsNil)
	s.client = NewDefaultClient(s.server)
}

func (s *ClientSuite) TestValidate(c *C) {
	c.Assert(s.client.ValidateServer("/"), IsNil)
}

func (s *ClientSuite) TestPropfind(c *C) {
	ms, err := s.client.WebDAV().Propfind("/", webdav.Depth0, webentities.NewAllPropsFind())
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
	path := fmt.Sprintf("/%s.ics", uuid)

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

	// create the master event object
	start := time.Now().Truncate(time.Hour).UTC()
	uid := fmt.Sprintf("test-recurring-event-%d", start.Unix())
	putEvent := components.NewEventWithDuration(uid, start, time.Hour)
	putEvent.Summary = "This is a test recurring event"
	rule := values.NewRecurrenceRule(values.DayRecurrenceFrequency)
	rule.Count = 14 // two weeks of events
	putEvent.AddRecurrenceRules(rule)

	// create an instance override at one week out
	nextWeek := start.AddDate(0, 0, 7)
	// start it an hour later, make it go for twice as long
	overrideEvent := components.NewEventWithDuration(uid, nextWeek.Add(time.Hour), 2*time.Hour)
	// mark it as an override of the recurrence at one week out
	overrideEvent.RecurrenceId = values.NewDateTime(nextWeek)
	overrideEvent.Summary = "This is a test override event"

	// generate an ICS filepath
	path := fmt.Sprintf("/%s.ics", uid)

	// save the events to the server
	if err := s.client.PutEvents(path, putEvent, overrideEvent); err != nil {
		c.Fatal(err.Error())
	}

	// create a query for all events between one week out + days in range
	daysInRange := 2
	nextWeekEnd := nextWeek.AddDate(0, 0, daysInRange)
	query, err := calentities.NewEventRangeQuery(nextWeek, nextWeekEnd)
	if err != nil {
		c.Fatal(err.Error())
	}

	// add in a filter for UID so that we don't get back unwanted results
	pf := calentities.NewPropertyMatcher(properties.UIDPropertyName, uid)
	query.Filter.ComponentFilter.ComponentFilter.PropertyFilter = pf

	// send the query to the server
	if events, err := s.client.QueryEvents("/", query); err != nil {
		c.Fatal(err.Error())
	} else {
		// since this is a daily recurring event, we should only get back one event for every day in our range
		c.Assert(events, HasLen, daysInRange)
		for i, event := range events {

			// all events should have the same UID
			c.Assert(event.UID, Equals, uid)

			// all events should have a consistant recurrence ID
			day := nextWeek.AddDate(0, 0, i)
			offset := values.NewDateTime(day)
			c.Assert(event.RecurrenceId, DeepEquals, offset)

			if i == 0 {
				// the first event is an override, it should start an hour later and be twice as long
				c.Assert(event.DateStart, DeepEquals, overrideEvent.DateStart)
				c.Assert(event.Duration, DeepEquals, overrideEvent.Duration)
				c.Assert(event.Summary, DeepEquals, overrideEvent.Summary)
			} else {
				// regular occurences should start at the same time as their recurrence ID
				c.Assert(event.DateStart, DeepEquals, event.RecurrenceId)
				c.Assert(event.Duration, DeepEquals, putEvent.Duration)
				c.Assert(event.Summary, DeepEquals, putEvent.Summary)
			}
		}
	}

}

func (s *ClientSuite) TestResetCalendar(c *C) {

	// only delete if the calendar exists
	if exists, err := s.client.WebDAV().Exists("/"); err != nil {
		c.Fatal(err.Error())
	} else if exists {
		c.Assert(s.client.WebDAV().Delete("/"), IsNil)
	}

	// now try to recreate the calendar
	c.Assert(s.client.MakeCalendar("/"), IsNil)

}

func AssertServerUrl(c *C) *url.URL {
	urlstr := AssertEnvString("CALDAV_SERVER_URL", c)
	uri, err := url.Parse(urlstr)
	c.Assert(err, IsNil)
	return uri
}

func AssertEnvString(name string, c *C) string {
	value := os.Getenv(name)
	c.Assert(value, Not(HasLen), 0)
	return value
}
