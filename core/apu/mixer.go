package apu

import "log"

type Mixer struct{}

func NewMixer() *Mixer {
	log.Println("[APU] Enter NewMixer")
	m := &Mixer{}
	log.Println("[APU] Exit NewMixer")
	return m
}

func (m *Mixer) Mix(p1, p2 *PulseChannel, t *TriangleChannel, n *NoiseChannel) float64 {
	return 0.4*p1.Output() + 0.4*p2.Output() + 0.2*t.Output() + 0.1*n.Output()
}
