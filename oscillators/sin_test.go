package oscillators

import (
	"math"
	"testing"
)

func TestCalculateStep(t *testing.T) {
	freq := uint(1)
	sr := uint(44100)
	st := CalculateStep(freq, sr)

	if freq != uint((st*44100)/(2*math.Pi)) {
		t.Errorf("Got %v multiple 44100 should have gotten %v but got %v", st, freq, (st*44100)/(2*math.Pi))
	}

	freq = uint(10)
	st = CalculateStep(freq, sr)

	if freq != uint((st*44100)/(2*math.Pi)) {
		t.Errorf("Got %v multiple 44100 should have gotten %v but got %v", st, freq, (st*44100)/(2*math.Pi))
	}

	freq = uint(440)
	st = CalculateStep(freq, sr)

	if freq != uint((st*44100)/(2*math.Pi)) {
		t.Errorf("Got %v multiple 44100 should have gotten %v but got %v", st, freq, (st*44100)/(2*math.Pi))
	}
}
