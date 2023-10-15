package values

import (
	"log"
	"strings"
)

var _ = log.Print

type CSV []string

func (csv *CSV) EncodeICalValue() (string, error) {
	log.Println("in EncodeICalValue:", *csv)
	return strings.Join(*csv, ","), nil
}

func (csv *CSV) DecodeICalValue(value string) error {
	log.Println("in DecodeICalValue:", *csv, value)
	value = strings.TrimSpace(value)
	*csv = CSV(strings.Split(value, ","))
	log.Println("in DecodeICalValue, out:", *csv)
	return nil
}

func NewCSV(items ...string) *CSV {
	log.Println("in NewCSV", items)
	return (*CSV)(&items)
}
