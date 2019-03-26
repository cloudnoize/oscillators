package oscillators

type Freq map[string]uint

func (f Freq) ToFreq(note string) uint {
	fe, ok := f[note]
	if !ok {
		println("Failed to convert ", note)
	}
	return fe
}

func NewFreq() Freq {
	f := make(Freq)
	f["a"] = 440
	f["b"] = 494
	f["c"] = 523
	f["d"] = 587
	f["e"] = 659
	f["f"] = 698
	f["g"] = 784
	return f
}
