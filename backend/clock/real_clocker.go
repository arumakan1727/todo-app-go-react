package clock

import "time"

type realClocker struct{}

func GetRealClocker() Clocker {
	return realClocker{}
}

func (_ realClocker) Now() time.Time {
	return time.Now()
}
