package ctime

import (
	"fmt"
	"time"
)

const (
	iso8601DATEFormat = "2006-01-02"
	iso8601Format     = "2006-01-02T15:04:05.000Z"
)

// ProviderFunc provides time
type ProviderFunc func() time.Time

// ISO8601 is an ISO-8601 formatted time
type ISO8601 struct {
	time time.Time
}

// Epoch is a representation of NANO seconds past epoch (unix) time
type Epoch int64

// EpochProviderFunc provides the functionality of providing an Epoch
type EpochProviderFunc func() Epoch

// NewISO8601 creates a new ISO8601 time
func NewISO8601(v string) (ISO8601, error) {
	t, err := time.Parse(iso8601DATEFormat, string(v))
	if err != nil {
		t, err := time.Parse(iso8601Format, string(v))
		if err != nil {
			return ISO8601{}, fmt.Errorf("time - unable to parse value. Format should be either :%s or %s ", iso8601DATEFormat, iso8601Format)
		}
		return ISO8601{time: t}, nil
	}
	return ISO8601{time: t}, nil
}

// New creates an ISO8601 from a Go Time
func New(t time.Time) (ISO8601, error) {
	isoTime, err := NewISO8601(t.Format(iso8601DATEFormat))
	if err != nil {
		isoTime, err := NewISO8601(t.Format(iso8601Format))
		if err != nil {
			return ISO8601{}, fmt.Errorf("time - unable to parse value. Format should be either :%s or %s ", iso8601DATEFormat, iso8601Format)
		}
		return isoTime, nil
	}
	return isoTime, nil
}

// String returns the string value of an ISO8601 time
// The output has 3 fractional seconds values and terminates with a Z
// An example output is 2006-01-02T15:04:05.629Z
func (t ISO8601) String() string {
	return t.time.Format(iso8601Format)
}

// DateString returns the date string value of an ISO8601 time
func (t ISO8601) DateString() string {
	return t.time.Format(iso8601DATEFormat)
}

// ToEpoch converts an ISO8601 time to epoch time
func (t ISO8601) ToEpoch() Epoch {
	return Epoch(t.time.UnixNano())
}

// Val returns the underlying value for an ISO time
func (t ISO8601) Val() time.Time {
	return t.time
}

func iso8601FromTime(t time.Time) ISO8601 {
	return ISO8601{time: t.UTC()}
}

// CurrentEpoch returns the number nano seconds past epoch time as of now
func CurrentEpoch() Epoch {
	return Epoch(time.Now().UnixNano())
}

// ToISO8601 converts an epoch time to ISO8601 formatted string
func (e Epoch) ToISO8601() ISO8601 {
	t := time.Unix(0, int64(e))
	return iso8601FromTime(t)
}

// String converts an epoch time to ISO8601 formatted string
func (e Epoch) String() string {
	return e.ToISO8601().String()
}
