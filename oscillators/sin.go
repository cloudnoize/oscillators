package oscillators

import (
	"math"
	"time"
)

type SinOsc struct {
	sr    uint
	step  float64
	freq  uint
	index uint
	ch    chan float32
	cont  bool
}

func NewSinOsc(sr, freq uint) Oscillator {
	so := SinOsc{sr: sr, freq: freq, step: CalculateStep(freq, sr), ch: make(chan float32, sr), cont: true}
	go so.sampleGen()
	return &so
}

func (s *SinOsc) calculateSample() float32 {
	ret := math.Sin(float64(s.index) * s.step)
	s.index = (s.index + 1) % s.sr
	return float32(ret)
}
func (s *SinOsc) sampleGen() {
	dur := time.Duration(uint(time.Second) / (s.sr + 1250))
	println("tikcer dur is ", dur.String())
	t := time.NewTicker(dur)
	s.ch <- s.calculateSample()
	for range t.C {
		s.ch <- s.calculateSample()
		if !s.cont {
			break
		}
	}
}

func (s *SinOsc) GetSample() float32 {
	return <-s.ch
}

func (s *SinOsc) Close() {
	s.GetSample()
	s.cont = false
}

func (s *SinOsc) SetFreq(freq uint) {
	s.freq = freq
}

func CalculateStep(freq, sr uint) float64 {
	step := float64(freq) * ((2 * math.Pi) / float64(sr))
	return step
}
