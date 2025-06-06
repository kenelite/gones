package apu

type APU struct {
	Pulse1   *PulseChannel
	Pulse2   *PulseChannel
	Triangle *TriangleChannel
	Noise    *NoiseChannel
	Mixer    *Mixer
	Output   *AudioOutput
}

func NewAPU() *APU {
	return &APU{
		Pulse1:   NewPulseChannel(),
		Pulse2:   NewPulseChannel(),
		Triangle: NewTriangleChannel(),
		Noise:    NewNoiseChannel(),
		Mixer:    NewMixer(),
		Output:   NewAudioOutput(),
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
