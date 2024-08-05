package time

import "time"

type Time interface {
	Now() time.Time
}

type timeImpl struct {
	currentTime time.Time
}

func (n timeImpl) Now() time.Time {
	return n.currentTime
}

func NewTimeImpl(time time.Time) timeImpl {
	return timeImpl{
		currentTime: time,
	}
}
