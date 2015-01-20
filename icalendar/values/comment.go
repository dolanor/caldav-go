package values

// specifies non-processing information intended to provide a comment to the calendar user.
type Comment string

// encodes the comment value for the iCalendar specification
func (c Comment) EncodeICalValue() string {
	return string(c)
}

// encodes the comment value for the iCalendar specification
func (c Comment) EncodeICalName() string {
	return "COMMENT"
}
