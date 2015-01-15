package providers

import (
	"fmt"
	"github.com/taviti/caldav-go/caldav"
	"github.com/taviti/caldav-go/webdav/constants"
	"github.com/taviti/caldav-go/webdav/entities"
	"os"
	"testing"
)

var client caldav.Client

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

func TestValidate(t *testing.T) {

	if err := client.Validate("/"); err != nil {
		t.Error(err)
	}

	return

}

func TestPropfind(t *testing.T) {

	// lookup current user's calendar
	username := client.Provider().CurrentUsername()
	path := client.Provider().UserCalendarPath(username)
	ms, err := client.Propfind(path, constants.Depth0, entities.AllProps())

	if err != nil {
		t.Error(err)
	} else if len(ms.Responses) == 0 {
		t.Errorf("no responses returned!")
	} else if len(ms.Responses[0].PropStats) == 0 {
		t.Errorf("no property metadata returned!")
	} else if ms.Responses[0].PropStats[0].Prop == nil {
		t.Errorf("no property value returned!")
	} else if ms.Responses[0].PropStats[0].Prop.ResourceType == nil {
		t.Errorf("no property type returned!")
	} else if ms.Responses[0].PropStats[0].Prop.ResourceType.Calendar == nil {
		t.Errorf("expected property %s to be a calendar!", ms.Responses[0].PropStats[0].Prop.DisplayName)
	}

	return

}
