package oscillators

import (
	"math"
)

type SinOsc struct {
	sr    uint
	step  float64
	freq  uint
	index uint
	cont  bool
}

func NewSinOsc(sr, freq uint) Oscillator {
	so := SinOsc{sr: sr, freq: freq, step: CalculateStep(freq, sr), cont: true}
	return &so
}

func (s *SinOsc) calculateSample() float32 {
	ret := math.Sin(float64(s.index) * s.step)
	s.index = (s.index + 1) % s.sr
	return float32(ret)
}

func (s *SinOsc) GetSample() float32 {
	return s.calculateSample()
}

func (s *SinOsc) Close() {
	s.GetSample()
	s.cont = false
}

func (s *SinOsc) SetFreq(freq uint) {
	println("Setting freq ", freq)
	s.freq = freq
	s.step = CalculateStep(freq, s.sr)
}

func CalculateStep(freq, sr uint) float64 {
	step := float64(freq) * ((2 * math.Pi) / float64(sr))
	return step
}
