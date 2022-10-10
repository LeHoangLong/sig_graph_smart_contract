package utility

import "time"

//go:generate mockgen -source=$GOFILE -destination ../testutils/clock.go -package mock
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
