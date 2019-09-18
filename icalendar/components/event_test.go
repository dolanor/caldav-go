package components

import (
	"fmt"
	"github.com/dolanor/caldav-go/icalendar"
	"github.com/dolanor/caldav-go/icalendar/values"
	. "gopkg.in/check.v1"
	"net/url"
	"testing"
	"time"
)

type EventSuite struct{}

var _ = Suite(new(EventSuite))

func TestEvent(t *testing.T) { TestingT(t) }

func (s *EventSuite) TestMissingEndMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEvent("test", now)
	_, err := icalendar.Marshal(event)
	c.Assert(err, ErrorMatches, "(?s).*end date or duration must be set.*")
}

func (s *EventSuite) TestBasicWithDurationMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDURATION:PT1H\r\nEND:VEVENT"
	fdate := now.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, fdate, fdate))
}

func (s *EventSuite) TestBasicWithEndMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	event := NewEventWithEnd("test", now, end)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate))
}

func (s *EventSuite) TestFullEventMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	oneDay := time.Hour * 24
	oneWeek := oneDay * 7
	event := NewEventWithEnd("1:2:3", now, end)
	uri, _ := url.Parse("http://dolanor.com/some/attachment.ics")
	event.Attachment = values.NewUrl(*uri)
	event.Attendees = []*values.AttendeeContact{
		values.NewAttendeeContact("Jon Azoff", "jon@dolanor.com"),
		values.NewAttendeeContact("Matthew Davie", "matthew@dolanor.com"),
	}
	event.Categories = values.NewCSV("vinyasa", "level 1")
	event.Comments = values.NewComments("Great class, 5 stars!", "I love this class!")
	event.ContactInfo = values.NewCSV("Send us an email!", "<jon@dolanor.com>")
	event.Created = event.DateStart
	event.Description = "An all-levels class combining strength and flexibility with breath"
	ex1 := values.NewDateTime(now.Add(oneWeek))
	ex2 := values.NewDateTime(now.Add(oneWeek * 2))
	event.ExceptionDateTimes = values.NewExceptionDateTimes(ex1, ex2)
	event.Geo = values.NewGeo(37.747643, -122.445400)
	event.LastModified = event.DateStart
	event.Location = values.NewLocation("Dolores Park")
	event.Organizer = values.NewOrganizerContact("Jon Azoff", "jon@dolanor.com")
	event.Priority = 1
	event.RecurrenceId = event.DateStart
	r1 := values.NewDateTime(now.Add(oneWeek + oneDay))
	r2 := values.NewDateTime(now.Add(oneWeek*2 + oneDay))
	event.RecurrenceDateTimes = values.NewRecurrenceDateTimes(r1, r2)
	event.AddRecurrenceRules(values.NewRecurrenceRule(values.WeekRecurrenceFrequency))
	uri, _ = url.Parse("matthew@dolanor.com")
	event.RelatedTo = values.NewUrl(*uri)
	event.Resources = values.NewCSV("yoga mat", "towel")
	event.Sequence = 1
	event.Status = values.TentativeEventStatus
	event.Summary = "Jon's Super-Sweaty Vinyasa 1"
	event.TimeTransparency = values.OpaqueTimeTransparency
	uri, _ = url.Parse("http://student.dolanor.com/san-francisco/jonathan-azoff/vinyasa-1")
	event.Url = values.NewUrl(*uri)
	enc, err := icalendar.Marshal(event)
	if err != nil {
		c.Fatal(err.Error())
	}
	tmpl := "BEGIN:VEVENT\r\nUID:1:2:3\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nCREATED:%sZ\r\n" +
		"DESCRIPTION:An all-levels class combining strength and flexibility with breath\r\n" +
		"GEO:37.747643 -122.445400\r\nLAST-MODIFIED:%sZ\r\nLOCATION:Dolores Park\r\n" +
		"ORGANIZER;CN=\"Jon Azoff\":MAILTO:jon@dolanor.com\r\nPRIORITY:1\r\nSEQUENCE:1\r\nSTATUS:TENTATIVE\r\n" +
		"SUMMARY:Jon's Super-Sweaty Vinyasa 1\r\nTRANSP:OPAQUE\r\n" +
		"URL;VALUE=URI:http://student.dolanor.com/san-francisco/jonathan-azoff/vinyasa-1\r\n" +
		"RECURRENCE-ID:%sZ\r\nRRULE:FREQ=WEEKLY\r\nATTACH;VALUE=URI:http://dolanor.com/some/attachment.ics\r\n" +
		"ATTENDEE;CN=\"Jon Azoff\":MAILTO:jon@dolanor.com\r\nATTENDEE;CN=\"Matthew Davie\":MAILTO:matthew@dolanor.com\r\n" +
		"CATEGORIES:vinyasa,level 1\r\nCOMMENT:Great class, 5 stars!\r\nCOMMENT:I love this class!\r\n" +
		"CONTACT:Send us an email!,<jon@dolanor.com>\r\nEXDATE:%s,%s\r\nRDATE:%s,%s\r\n" +
		"RELATED-TO;VALUE=URI:matthew@dolanor.com\r\nRESOURCES:yoga mat,towel\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate, sdate, sdate, sdate, ex1, ex2, r1, r2))
}

func (s *EventSuite) TestQualifiers(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	c.Assert(event.IsRecurrence(), Equals, false)
	event.RecurrenceId = values.NewDateTime(now)
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, false)
	event.DateStart = values.NewDateTime(now.Add(time.Hour))
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, true)
}

func (s *EventSuite) TestUnmarshalMultipleLines(c *C) {
	// Event that has an ATTENDEE that spans 3 lines
	raw := `BEGIN:VEVENT
DTSTART;TZID=America/Los_Angeles:20150511T140000
DTEND;TZID=America/Los_Angeles:20150511T150000
DTSTAMP:20150511T204516Z
ORGANIZER;CN=Fakebiz Shared:mailto:fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikm
 aeq8@group.calendar.google.com
UID:na9njgloe10sch3h0uootli104@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;CN=Fakebiz
  Shared;X-NUM-GUESTS=0:mailto:fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikm
 aeq8@group.calendar.google.com
CREATED:20150504T173946Z
DESCRIPTION:
LAST-MODIFIED:20150511T204516Z
LOCATION:Outer space
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:Brand Presentation
TRANSP:OPAQUE
END:VEVENT`

	e := Event{}
	err := icalendar.Unmarshal(raw, &e)
	c.Assert(err, IsNil)
	c.Assert(len(e.Attendees), Equals, 1)
	c.Assert(e.Attendees[0].Entry.Address, Equals, "fakemcfakebiz.com_b3a0grbjdr4dcje2fc4ikmaeq8@group.calendar.google.com")
	c.Assert(e.Attendees[0].Entry.Name, Equals, "Fakebiz Shared")
}
