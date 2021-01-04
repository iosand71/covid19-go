package main

import "time"

const (
	timeFormat = "2006-01-02T15:04:05"
	dateFormat = "2006-01-02"
)

// DateTime custom date/time
type DateTime time.Time

//Set date using custom format
func (d *DateTime) Set(v string) error {
	parsed, err := time.Parse(dateFormat, v)
	if err != nil {
		return err
	}
	*d = DateTime(parsed)
	return nil
}

func (d *DateTime) isZero() bool {
	date := time.Time(*d)
	return date.IsZero()
}

func (d *DateTime) String() string {
	date := time.Time(*d)
	return date.Format(dateFormat)
}

// Config holder for cli configs
type Config struct {
	Region       string
	startDate    DateTime
	printRegions bool
	debug        bool
}
