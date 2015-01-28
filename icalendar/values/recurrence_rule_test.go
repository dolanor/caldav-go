package values

import (
	"fmt"
	. "github.com/taviti/check"
	"testing"
	"time"
)

type RecurrenceRuleSuite struct{}

var _ = Suite(new(RecurrenceRuleSuite))

func TestRecurrenceRule(t *testing.T) { TestingT(t) }

func (s *RecurrenceRuleSuite) TestEncode(c *C) {
	date, err := time.ParseInLocation(DateTimeFormatString, DateTimeFormatString, time.UTC)
	c.Assert(err, IsNil)
	r := NewRecurrenceRule(WeekRecurrenceFrequency)
	r.Until = NewDateTime(date)
	r.Count = 1
	r.Interval = 2
	r.BySecond = []int{3}
	r.ByMinute = []int{4}
	r.ByHour = []int{5, 6}
	r.ByDay = append(r.ByDay, MondayRecurrenceWeekday)
	r.ByDay = append(r.ByDay, TuesdayRecurrenceWeekday)
	r.ByMonthDay = []int{7, 8}
	r.ByYearDay = []int{9, 10, 11}
	r.ByWeekNumber = []int{12}
	r.ByMonth = []int{3}
	r.BySetPosition = []int{1}
	r.WeekStart = SundayRecurrenceWeekday
	fs := "FREQ=WEEKLY;UNTIL=%sZ;COUNT=1;INTERVAL=2;BYSECOND=3;BYMINUTE=4;BYHOUR=5,6;BYDAY=MO,TU;BYMONTHDAY=7,8;BYYEARDAY=9,10,11;BYWEEKNO=12;BYMONTH=3;BYSETPOS=1;WKST=SU"
	expected := fmt.Sprintf(fs, date.Format(DateTimeFormatString))
	actual, err := r.EncodeICalValue()
	c.Assert(err, IsNil)
	c.Assert(actual, Equals, expected)
}
