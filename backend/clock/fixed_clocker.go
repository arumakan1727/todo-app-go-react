package clock

import "time"

type fixedClocker struct {
	t time.Time
}

func GetFixedClocker() Clocker {
	return fixedClocker{
		t: time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC),
	}
}

func GetFixedClockerOf(t time.Time) Clocker {
	return &fixedClocker{
		t: t,
	}
}

func (c fixedClocker) Now() time.Time {
	return c.t
}

func (c fixedClocker) Location() *time.Location {
	return c.t.Location()
}

func (c fixedClocker) In(loc *time.Location) Clocker {
	return GetFixedClockerOf(c.t.In(loc))
}
