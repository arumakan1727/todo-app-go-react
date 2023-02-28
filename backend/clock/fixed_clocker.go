package clock

import "time"

type fixedClocker struct {
	t time.Time
}

func GetFixedClocker(t *time.Time) Clocker {
	c := fixedClocker{}
	if t == nil {
		c.t = time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
	}
	return c
}

func (c fixedClocker) Now() time.Time {
	return c.t
}
