package components

import (
	"fmt"
	"github.com/taviti/caldav-go/icalendar"
	"github.com/taviti/caldav-go/icalendar/values"
	. "github.com/taviti/check"
	"net/mail"
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
	c.Assert(err, ErrorMatches, "DateEnd or Duration must be set")
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
	uri, _ := url.Parse("http://taviti.com/some/attachment.ics")
	event.Attachment = values.NewUrl(*uri)
	event.Attendees = []*values.AttendeeAddress{
		values.NewAttendeeAddress(mail.Address{Name: "Jon Azoff", Address: "jon@taviti.com"}),
		values.NewAttendeeAddress(mail.Address{Name: "Matthew Davie", Address: "matthew@taviti.com"}),
	}
	event.Categories = values.CSV{"vinyasa", "level 1"}
	event.Comments = []values.Comment{"Great class, 5 stars!", "I love this class!"}
	event.ContactInfo = values.CSV{"Send us an email!", "<jon@taviti.com>"}
	event.Created = event.DateStart
	event.Description = "An all-levels class combining strength and flexibility with breath"
	ex1 := values.NewDateTime(now.Add(oneWeek))
	ex2 := values.NewDateTime(now.Add(oneWeek * 2))
	event.ExceptionDateTimes = values.ExceptionDateTimes{ex1, ex2}
	event.Geo = values.NewGeo(37.747643, -122.445400)
	event.LastModified = event.DateStart
	event.Location = "Dolores Park"
	event.Organizer = values.NewOrganizerAddress(mail.Address{Name: "Jon Azoff", Address: "jon@taviti.com"})
	event.Priority = 1
	event.RecurrenceId = event.DateStart
	r1 := values.NewDateTime(now.Add(oneWeek + oneDay))
	r2 := values.NewDateTime(now.Add(oneWeek*2 + oneDay))
	event.RecurrenceDateTimes = values.RecurrenceDateTimes{r1, r2}
	event.RecurrenceRule = values.NewRecurrenceRule(values.WeekRecurrenceFrequency)
	event.RelatedTo = values.NewRelationAddress(mail.Address{Address: "matthew@taviti.com"})
	event.Resources = values.CSV{"yoga mat", "towel"}
	event.Sequence = 1
	event.Status = values.TentativeEventStatus
	event.Summary = "Jon's Super-Sweaty Vinyasa 1"
	event.TimeTransparency = values.OpaqueTimeTransparency
	uri, _ = url.Parse("http://student.taviti.com/san-francisco/jonathan-azoff/vinyasa-1")
	event.Url = values.NewUrl(*uri)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:1\\:2\\:3\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nCREATED:%sZ\r\n" +
		"DESCRIPTION:An all-levels class combining strength and flexibility with breath\r\n" +
		"GEO:37.747643 -122.445400\r\nLAST-MODIFIED:%sZ\r\nLOCATION:Dolores Park\r\n" +
		"ORGANIZER;CN=Jon Azoff:MAILTO\\:jon@taviti.com\r\nPRIORITY:1\r\nSEQUENCE:1\r\nSTATUS:TENTATIVE\r\n" +
		"SUMMARY:Jon's Super-Sweaty Vinyasa 1\r\nTRANSP:OPAQUE\r\n" +
		"URL:http\\://student.taviti.com/san-francisco/jonathan-azoff/vinyasa-1\r\n" +
		"RECURRENCE-ID:%sZ\r\nRRULE:FREQ=WEEKLY\r\nATTACH:http\\://taviti.com/some/attachment.ics\r\n" +
		"ATTENDEE;CN=Jon Azoff:MAILTO\\:jon@taviti.com\r\nATTENDEE;CN=Matthew Davie:MAILTO\\:matthew@taviti.com\r\n" +
		"CATEGORIES:vinyasa,level 1\r\nCOMMENT:Great class, 5 stars!\r\nCOMMENT:I love this class!\r\n" +
		"CONTACT:Send us an email!,<jon@taviti.com>\r\nEXDATE:%s,%s\r\nRDATE:%s,%s\r\n" +
		"RELATED-TO:<matthew@taviti.com>\r\nRESOURCES:yoga mat,towel\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate, sdate, sdate, sdate, ex1, ex2, r1, r2))
}
