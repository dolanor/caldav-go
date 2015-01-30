# caldav-go
A CalDAV ([rfc4791][1]) and iCalendar client for Go

Project Structure
---------------
This project contains several modules, all of which work together to to allow for
calendaring using CalDAV:

- An HTTP Client
- A WebDAV Client and Mapping Layer
- A CalDAV Client and Mapping Layer
- An iCalendar (Un)Marshaler


Getting Started
---------------
To install the library into your project, add caldev-go into your `GOPATH`:

```sh
$ go get github.com/taviti/caldav-go
```

Then, in your go application, include the caldav client and start making requests:

```go

import "github.com/taviti/caldav-go/caldav"

// create a reference to your CalDAV-compliant server
var server = caldav.NewServer("http://my-caldav-host.net:8008")

// create a CalDAV client to speak to the server
var client = caldav.NewClient(server, http.DefaultClient)

// start executing requests!
err := client.ValidateServer()
```

Testing
-------
To test the client, you must first have access to (or run your own) [caldav-compliant server][1]. On the machine
you wish to test on, ensure that the `CALDAV_SERVER_URL` environment variable is set to the host and path of the
account you wish to run tests on. Afterwords, the standard `go test` command will run the tests for the whole library.
For instance, if you have a server running locally on port 8008, you could run the tests in one command:

```sh
CALDAV_SERVER_URL='http://localhost:8008/calendars/users/admin/calendar/' go test ./...
```

HTTP Basic authentication can be baked into the URL as well:

```sh
CALDAV_SERVER_URL='http://admin:admin@localhost:8008/calendars/users/admin/calendar/' go test ./...
```

[1]:http://tools.ietf.org/html/rfc4791
[2]:http://calendarserver.org/