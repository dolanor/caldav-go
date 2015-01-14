package caldav

// a client for communicating with a CalDAV server
type Client interface {

	// verifies that the CalDAV provider can, in fact, speak CalDAV
	// you may supply an optional list of extra capabilities to test the server for
	// returns an error if the provider does not support CalDAV or any provided extra capabilities
	Verify(capabilities ...string) error
}
