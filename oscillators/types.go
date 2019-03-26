package oscillators

import (
	"net"
	"time"
)

type Oscillator interface {
	GetSample() float32
	SetFreq(uint)
	Close()
}

type ClientContext struct {
	Oscillator
	Addr  net.Addr
	Start time.Time
	Misc  string
}

type ClientContexts map[uint]*ClientContext
