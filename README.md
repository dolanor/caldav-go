# caldav-go
A CalDAV ([rfc4791][1]) Client for Go

Getting Started
---------------
To install the library into your project, add caldev-go into your `GOPATH`:

```sh
$ go get github.com/taviti/caldav-go
```

Then, in your go application, include the caldav client and start making requests:

```go

import "github.com/taviti/caldav-go/caldav"

// select from any provider in the providers folder, or create your own!
var provider caldav.Provider

// create the caldav client
var client = caldav.NewClient(provider, http.DefaultClient)
```

Testing
-------
To test the client, you must first have access to (or run your own) [caldav-compliant server][1]. On the machine you wish to test
from, ensure that the `CALDAV_SERVER_URL` environment variable is set. Afterwords, the standard `go test` command will run the tests
for the package. For instance, if you have a server running locally on port 8008, you could run the tests in one command:

```sh
CALDAV_SERVER_URL='http://localhost:8008' go test github.com/taviti/caldav-go/...
```

HTTP Basic authentication can be baked into the URL as well:

```sh
CALDAV_SERVER_URL='http://user:pass@localhost:8008' go test github.com/taviti/caldav-go/...
```

[1]:http://tools.ietf.org/html/rfc4791
[2]:http://calendarserver.org/