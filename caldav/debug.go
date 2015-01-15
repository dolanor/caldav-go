package caldav

import (
	"log"
)

const DEBUG_ENABLED = true

func logf(format string, args ...interface{}) {
	if DEBUG_ENABLED {
		log.Printf(format, args...)
	}
}
