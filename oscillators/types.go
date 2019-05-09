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
	Osc   Oscillators
}

type ClientContexts map[uint]*ClientContext

type Oscillators int

const (
	Sin  Oscillators = 0
	Test Oscillators = 1
)

func OscStr(osc int) string {
	switch osc {
	case 0:
		return "sin"
	case 1:
		return "test"
	}
	return "sin"
}

func GetOsc(osc Oscillators) Oscillator {
	switch osc {
	case 0:
		return NewSinOsc(44100, 440)
	case 1:
		return NewTestOsc(44100)
	}
	return NewSinOsc(44100, 440)
}
