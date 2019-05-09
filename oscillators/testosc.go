package oscillators

type TestOsc struct {
	sr   uint
	ch   chan float32
	cont bool
}

func NewTestOsc(sr uint) Oscillator {
	to := TestOsc{sr: sr, ch: make(chan float32, 100), cont: true}
	go to.sampleGen()
	return &to
}

func (s *TestOsc) sampleGen() {
	s.ch <- float32(1)
}

func (s *TestOsc) GetSample() float32 {
	return float32(1)
}

func (s *TestOsc) Close() {
	s.GetSample()
	s.cont = false
}

func (s *TestOsc) SetFreq(freq uint) {

}
