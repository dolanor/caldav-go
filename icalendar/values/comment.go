package values

// specifies non-processing information intended to provide a comment to the calendar user.
type Comment string

// encodes the comment value for the iCalendar specification
func (c Comment) EncodeICalValue() (string, error) {
	return string(c), nil
}

// decodes the comment value from the iCalendar specification
func (c Comment) DecodeICalValue(value string) error {
	c = Comment(value)
	return nil
}

// encodes the comment value for the iCalendar specification
func (c Comment) EncodeICalName() (string, error) {
	return "COMMENT", nil
}
