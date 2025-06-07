package apu

type TriangleChannel struct {
	phase   float64
	freq    float64
	enabled bool
}

func NewTriangleChannel() *TriangleChannel {
	return &TriangleChannel{}
}

func (t *TriangleChannel) Step() {
	if t.enabled {
		t.phase += t.freq
		if t.phase > 1.0 {
			t.phase -= 1.0
		}
	}
}

func (t *TriangleChannel) Output() float64 {
	if !t.enabled {
		return 0
	}
	// 线性三角波 -1.0 到 1.0
	if t.phase < 0.5 {
		return 4.0*t.phase - 1.0
	}
	return 3.0 - 4.0*t.phase
}
