package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const SampleRate = 44100

var context *audio.Context

func InitAudio() {
	context = audio.NewContext(SampleRate)
}

func GetContext() *audio.Context {
	return context
}
