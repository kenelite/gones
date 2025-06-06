package audio

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type PulseChannel struct {
	frequency float64
	volume    float64
	phase     float64
	player    *audio.Player
}

// pulseWaveReader 用于生成方波音频流
// 实现 audio.Streamer 接口

type pulseWaveReader struct {
	freq   float64
	vol    float64
	phase  float64
	sample int
}

func (r *pulseWaveReader) Read(buf []byte) (int, error) {
	// 每次输出一帧 16-bit PCM (小端)
	for i := 0; i+1 < len(buf); i += 2 {
		period := float64(r.sample) / r.freq
		var v float64
		if int(r.phase) < int(period/2) {
			v = r.vol
		} else {
			v = -r.vol
		}
		// 16-bit PCM
		sample := int16(v * 32767)
		buf[i] = byte(sample)
		buf[i+1] = byte(sample >> 8)
		r.phase++
		if r.phase >= period {
			r.phase -= period
		}
	}
	return len(buf), nil
}

func (r *pulseWaveReader) Seek(offset int64, whence int) (int64, error) {
	r.phase = 0
	return 0, nil
}

func NewPulse(frequency float64, volume float64) *PulseChannel {
	p := &PulseChannel{
		frequency: frequency,
		volume:    volume,
	}
	p.init()
	return p
}

func (p *PulseChannel) init() {
	stream := audio.NewInfiniteLoop(&pulseWaveReader{
		freq:   p.frequency,
		vol:    p.volume,
		phase:  0,
		sample: SampleRate,
	}, int64(time.Second))
	p.player, _ = audio.NewPlayer(GetContext(), stream)
	p.player.Play()
}

// 停止播放
func (p *PulseChannel) Stop() {
	if p.player != nil {
		p.player.Pause()
	}
}

// 复写频率
func (p *PulseChannel) SetFrequency(freq float64) {
	p.frequency = freq
}
