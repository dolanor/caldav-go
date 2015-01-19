package values

import (
	"time"
)

type Duration struct {
	d time.Duration
}

func (d *Duration) EncodeICalValue() (string, error) {
	return "", nil
}
