package clock

import "time"

type Clocker interface {
	Now() time.Time
	Location() *time.Location
	In(*time.Location) Clocker
}
