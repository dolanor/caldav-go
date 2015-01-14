package providers

import (
	"fmt"
	"github.com/taviti/caldav-go/caldav"
	"os"
	"testing"
)

var client caldav.Client

// sets up the client for the tests
func TestMain(m *testing.M) {
	if rawurl := os.Getenv("CALDAV_SERVER_URL"); rawurl == "" {
		panic("unable to find value for environment variable CALDAV_SERVER_URL")
	} else if p, err := NewCalendarServerProvider(rawurl); err != nil {
		panic(fmt.Sprintf("unable to parse environment variable CALDAV_SERVER_URL, %s", err))
	} else {
		client = caldav.NewDefaultClient(p)
	}
	os.Exit(m.Run())
}

// ensures that the test client is using a valid CalDAV provider
func TestVerify(t *testing.T) {

	if err := client.Verify(); err != nil {
		t.Error(err)
	}

	return

}
