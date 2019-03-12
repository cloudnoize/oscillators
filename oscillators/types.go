package oscillators

type Oscillator interface {
	GetSample() float32
	SetFreq(uint)
	Close()
}
