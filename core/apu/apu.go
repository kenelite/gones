package apu

import "log"

type APU struct {
	Pulse1   *PulseChannel
	Pulse2   *PulseChannel
	Triangle *TriangleChannel
	Noise    *NoiseChannel
	Mixer    *Mixer
	Output   *AudioOutput
}

func NewAPU() *APU {
	log.Println("[APU] Initializing...")
	pulse1 := NewPulseChannel()
	if pulse1 == nil {
		log.Println("[APU] NewPulseChannel failed!")
	}
	pulse2 := NewPulseChannel()
	if pulse2 == nil {
		log.Println("[APU] NewPulseChannel (2) failed!")
	}
	triangle := NewTriangleChannel()
	if triangle == nil {
		log.Println("[APU] NewTriangleChannel failed!")
	}
	noise := NewNoiseChannel()
	if noise == nil {
		log.Println("[APU] NewNoiseChannel failed!")
	}
	mixer := NewMixer()
	if mixer == nil {
		log.Println("[APU] NewMixer failed!")
	}
	output := NewAudioOutput()
	if output == nil {
		log.Println("[APU] NewAudioOutput failed!")
	}
	log.Println("[APU] All submodules initialized.")
	return &APU{
		Pulse1:   pulse1,
		Pulse2:   pulse2,
		Triangle: triangle,
		Noise:    noise,
		Mixer:    mixer,
		Output:   output,
	}
}

// Step 模拟每一 CPU 周期，通常 APU 运行频率为 1.789773 MHz
func (a *APU) Step() {
	a.Pulse1.Step()
	a.Pulse2.Step()
	a.Triangle.Step()
	a.Noise.Step()

	sample := a.Mixer.Mix(a.Pulse1, a.Pulse2, a.Triangle, a.Noise)
	a.Output.EnqueueSample(sample)
}
