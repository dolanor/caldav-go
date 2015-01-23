package utils

import (
	"github.com/taviti/check"
	"net/url"
	"os"
)

func AssertServerUrl(c *check.C) *url.URL {
	urlstr := AssertEnvString("CALDAV_SERVER_URL", c)
	uri, err := url.Parse(urlstr)
	c.Assert(err, check.IsNil)
	return uri
}

func AssertEnvString(name string, c *check.C) string {
	value := os.Getenv(name)
	c.Assert(value, check.Not(check.HasLen), 0)
	return value
}
