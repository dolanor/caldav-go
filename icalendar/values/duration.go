package values

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// a representation of duration for iCalendar
type Duration struct {
	d time.Duration
}

// breaks apart the duration into its component time parts
func (d *Duration) Decompose() (weeks, days, hours, minutes, seconds int64) {

	// chip away at this
	rem := time.Duration(math.Abs(float64(d.d)))

	div := time.Hour * 24 * 7
	weeks = int64(rem / div)
	rem = rem % div
	div = div / 7
	days = int64(rem / div)
	rem = rem % div
	div = div / 24
	hours = int64(rem / div)
	rem = rem % div
	div = div / 60
	minutes = int64(rem / div)
	rem = rem % div
	div = div / 60
	seconds = int64(rem / div)

	return

}

// returns true if the duration is negative
func (d *Duration) IsPast() bool {
	return d.d < 0
}

// encodes the duration of time into iCalendar format
func (d *Duration) EncodeICalValue() string {
	var parts []string
	weeks, days, hours, minutes, seconds := d.Decompose()
	if d.IsPast() {
		parts = append(parts, "-")
	}
	parts = append(parts, "P")
	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%dW", weeks))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dD", days))
	}
	if hours > 0 || minutes > 0 || seconds > 0 {
		parts = append(parts, "T")
		if hours > 0 {
			parts = append(parts, fmt.Sprintf("%dH", hours))
		}
		if minutes > 0 {
			parts = append(parts, fmt.Sprintf("%dM", minutes))
		}
		if seconds > 0 {
			parts = append(parts, fmt.Sprintf("%dS", seconds))
		}
	}
	return strings.Join(parts, "")
}

func (d *Duration) String() string {
	return d.EncodeICalValue()
}

// creates a new iCalendar duration representation
func NewDuration(d time.Duration) *Duration {
	return &Duration{d: d}
}
