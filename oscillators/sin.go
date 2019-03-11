package oscillators

import "math"

type SinOsc struct {
	sr    uint
	step  float64
	freq  uint
	index uint
}

func NewSinOsc(sr, freq uint) Oscillator {
	return &SinOsc{sr: sr, freq: freq, step: CalculateStep(freq, sr)}
}

func (s *SinOsc) GetSample() float32 {
	ret := math.Sin(float64(s.index) * s.step)
	s.index = (s.index + 1) % s.sr
	return float32(ret)
}

func (s *SinOsc) SetFreq(freq uint) {
	s.freq = freq
}

func CalculateStep(freq, sr uint) float64 {
	step := float64(freq) * ((2 * math.Pi) / float64(sr))
	return step
}
