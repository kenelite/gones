package apu

import (
	"io"
	"log"
	"math"
	"sync"
	"time"

	oto "github.com/ebitengine/oto/v3"
)

const (
	sampleRate = 44100
)

type AudioOutput struct {
	player *oto.Player
	buf    *audioBuffer
}

// 实现一个 io.Reader 用于 PCM 数据流
// audioBuffer 实现 io.Reader，供 oto.Player 读取
// 这里我们用一个简单的 ring buffer

type audioBuffer struct {
	data []byte
	mu   sync.Mutex
}

func (b *audioBuffer) Read(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}

var (
	globalOtoCtx   *oto.Context
	globalOtoReady <-chan struct{}
	globalOtoErr   error
	otoOnce        sync.Once
)

func getGlobalOtoContext() (*oto.Context, <-chan struct{}, error) {
	otoOnce.Do(func() {
		globalOtoCtx, globalOtoReady, globalOtoErr = oto.NewContext(&oto.NewContextOptions{
			SampleRate:   sampleRate,
			ChannelCount: 1,
			Format:       oto.FormatSignedInt16LE,
		})
	})
	return globalOtoCtx, globalOtoReady, globalOtoErr
}

func NewAudioOutput() *AudioOutput {
	log.Println("[APU] Enter NewAudioOutput")
	ctx, ready, err := getGlobalOtoContext()
	if err != nil {
		log.Fatalf("[APU] oto.NewContext error: %v", err)
	}
	select {
	case <-ready:
		// 正常
	case <-time.After(3 * time.Second):
		log.Fatalf("[APU] oto.NewContext ready timeout")
	}
	buf := &audioBuffer{}
	player := ctx.NewPlayer(buf)
	player.Play()
	log.Println("[APU] Exit NewAudioOutput")
	return &AudioOutput{player: player, buf: buf}
}

func (a *AudioOutput) EnqueueSample(sample float64) {
	s := int16(math.MaxInt16 * sample)
	a.buf.mu.Lock()
	a.buf.data = append(a.buf.data, byte(s), byte(s>>8))
	a.buf.mu.Unlock()
}
