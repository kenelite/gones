package apu

import "log"

type PulseChannel struct {
	frequency float64
	volume    float64
	phase     float64
	duty      float64
	enabled   bool
}

func NewPulseChannel() *PulseChannel {
	log.Println("[APU] Enter NewPulseChannel")
	ch := &PulseChannel{
		duty: 0.5,
	}
	log.Println("[APU] Exit NewPulseChannel")
	return ch
}

func (p *PulseChannel) Step() {
	if !p.enabled {
		return
	}
	p.phase += p.frequency
	if p.phase > 1.0 {
		p.phase -= 1.0
	}
}

func (p *PulseChannel) Output() float64 {
	if !p.enabled {
		return 0
	}
	if p.phase < p.duty {
		return p.volume
	}
	return -p.volume
}
