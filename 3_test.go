package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/dolanor/caldav-go/icalendar"
	"github.com/dolanor/caldav-go/icalendar/components"
)

const icsData = `BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART;VALUE=DATE:20120114
DTEND;VALUE=DATE:20120115
DTSTAMP:20150209T114140Z
UID:REDACTED_IDENTIFIER@google.com
CATEGORIES:http://schemas.google.com/g/2005
END:VEVENT
END:VCALENDAR`

const icsData3 = `BEGIN:VCALENDAR
CALSCALE:GREGORIAN
VERSION:2.0
PRODID:-//SabreDAV//SabreDAV//EN
BEGIN:VEVENT

DTSTART;VALUE=DATE:20120114
DTEND;VALUE=DATE:20120115
DTSTAMP:20150209T114140Z
UID:REDACTED_IDENTIFIER@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;X-NUM-GUESTS=0:mailto:REDACTED@REDACTED.COM
CREATED:20111214T180041Z
DESCRIPTION:
LAST-MODIFIED:20111214T180041Z
LOCATION:Thuis
SEQUENCE:1
STATUS:TENTATIVE
SUMMARY:REDACTED BECAUSE PRIVACY
TRANSP:OPAQUE
CATEGORIES:http://schemas.google.com/g/2005
END:VEVENT
END:VCALENDAR`

const icsData2 = `BEGIN:VCALENDAR
CALSCALE:GREGORIAN
VERSION:2.0
PRODID:-//SabreDAV//SabreDAV//EN
BEGIN:VEVENT
CONTACT:REDACTED@REDACTED.COM
DTSTART;VALUE=DATE:20120114
DTEND;VALUE=DATE:20120115
DTSTAMP:20150209T114140Z
UID:REDACTED_IDENTIFIER@google.com
ATTENDEE;CUTYPE=INDIVIDUAL;ROLE=REQ-PARTICIPANT;PARTSTAT=ACCEPTED;X-NUM-GUESTS=0:mailto:REDACTED@REDACTED.COM
CREATED:20111214T180041Z
DESCRIPTION:
LAST-MODIFIED:20111214T180041Z
LOCATION:Thuis
SEQUENCE:1
STATUS:TENTATIVE
SUMMARY:REDACTED BECAUSE PRIVACY
TRANSP:OPAQUE
CATEGORIES:http://schemas.google.com/g/2005
END:VEVENT
END:VCALENDAR`

func TestHydrate(t *testing.T) {
	var cal components.Calendar
	err := icalendar.Unmarshal(icsData, &cal)
	if err != nil {
		t.Fatal(err)
	}
	ical, err := icalendar.Marshal(cal)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("final version:", ical)
}

type csv []string

func TestAssign(t *testing.T) {
	s := "ma string"
	_ = s
	k := csv{s}
	c := csv{}

	vs := reflect.ValueOf(s)
	vk := reflect.ValueOf(&k)
	vc := reflect.ValueOf(&c)

	lens := vs.Len()
	lenk := vk.Elem().Len()
	lenc := vc.Elem().Len()
	if vc.Elem().Kind() != reflect.Slice || vk.Elem().Kind() != reflect.Slice {
		t.Fatal("lenc or lenk is not a slice")
	}
	for i := 0; i < lenk; i++ {
		if !vc.Elem().CanSet() {
			t.Fatal("vc can't be set")
		}
		vkelem := vk.Elem()
		log.Println("vkelem kind", vkelem.Kind())
		vc.Elem().Set(reflect.Append(vc.Elem(), vkelem.Index(i)))
	}
	fmt.Println("lens", lens, "lenk", lenk, "lenc", lenc, vc)
	fmt.Println(vc)
	//c = append(c, k)

}
