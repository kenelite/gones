package apu

import "log"

type TriangleChannel struct {
	phase   float64
	freq    float64
	enabled bool
}

func NewTriangleChannel() *TriangleChannel {
	log.Println("[APU] Enter NewTriangleChannel")
	ch := &TriangleChannel{}
	log.Println("[APU] Exit NewTriangleChannel")
	return ch
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
