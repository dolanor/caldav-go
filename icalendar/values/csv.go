package values

import "strings"

type CSV []string

func (csv CSV) EncodeICalValue() string {
	return strings.Join(csv, ",")
}
