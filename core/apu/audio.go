package apu

import (
	"io"
	"math"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate = 44100
)

type AudioOutput struct {
	player oto.Player
	buf    []byte
	idx    int
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
	ctx, ready, _ := oto.NewContext(sampleRate, 1, 2)
	<-ready
	buf := &audioBuffer{}
	player := ctx.NewPlayer(buf)
	player.Play()
	return &AudioOutput{
		player: player,
		buf:    make([]byte, 2048),
	}
}

func (a *AudioOutput) EnqueueSample(sample float64) {
	// 将 float64 转为 16-bit PCM
	s := int16(math.MaxInt16 * sample)
	a.buf[a.idx] = byte(s)
	a.buf[a.idx+1] = byte(s >> 8)
	a.idx += 2
	if a.idx >= len(a.buf) {
		// 写入 audioBuffer
		ab := a.player.(interface{ Reader() io.Reader })
		if buf, ok := ab.Reader().(*audioBuffer); ok {
			buf.mu.Lock()
			buf.data = append(buf.data, a.buf...)
			buf.mu.Unlock()
		}
		a.idx = 0
	}
}
