package utility

import "time"

type ClockI interface {
	Now_ms() uint64
}

type clockRealtime struct {
}

func NewClockRealtime() ClockI {
	return &clockRealtime{}
}

func (c *clockRealtime) Now_ms() uint64 {
	return uint64(time.Now().UnixMilli())
}
