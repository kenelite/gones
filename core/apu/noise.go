package apu

import (
	"log"
	"math/rand"
)

type NoiseChannel struct {
	enabled bool
	volume  float64
}

func NewNoiseChannel() *NoiseChannel {
	log.Println("[APU] Enter NewNoiseChannel")
	ch := &NoiseChannel{}
	log.Println("[APU] Exit NewNoiseChannel")
	return ch
}

func (n *NoiseChannel) Step() {
	// 可实现 LFSR 随机生成，简化为随机
}

func (n *NoiseChannel) Output() float64 {
	if !n.enabled {
		return 0
	}
	return (rand.Float64()*2 - 1) * n.volume
}
