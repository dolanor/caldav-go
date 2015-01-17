package components

type Calendar struct {
	Version string `ical:",2.0"`
	ProdId  string `ical:",-//taviti/caldav-go//NONSGML v1.0.0//EN"`
	*Event
}
