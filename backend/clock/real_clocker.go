package clock

import "time"

type realClocker struct {
	loc *time.Location
}

func GetRealClocker(loc *time.Location) Clocker {
	return realClocker{
		loc: loc,
	}
}

func (c realClocker) Now() time.Time {
	return time.Now().In(c.loc)
}

func (c realClocker) Location() *time.Location {
	return c.loc
}

func (c realClocker) In(loc *time.Location) Clocker {
	return GetRealClocker(loc)
}
