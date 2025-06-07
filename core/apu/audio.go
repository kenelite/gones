package apu

import (
	"io"
	"math"
	"sync"

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

func NewAudioOutput() *AudioOutput {
	ctx, ready, _ := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	<-ready
	buf := &audioBuffer{}
	player := ctx.NewPlayer(buf)
	player.Play()
	return &AudioOutput{player: player, buf: buf}
}

func (a *AudioOutput) EnqueueSample(sample float64) {
	s := int16(math.MaxInt16 * sample)
	a.buf.mu.Lock()
	a.buf.data = append(a.buf.data, byte(s), byte(s>>8))
	a.buf.mu.Unlock()
}
