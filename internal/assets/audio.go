package assets

import (
	"embed"
	"log"
)

//go:embed embed/audio/*.wav
var audioFiles embed.FS

// GetAudio returns raw []byte audio file by name (e.g. "startup.wav")
func GetAudio(name string) []byte {
	data, err := audioFiles.ReadFile("embed/audio/" + name)
	if err != nil {
		log.Fatalf("failed to load audio: %s: %v", name, err)
	}
	return data
}

func LoadAudio() {
	// Placeholder: preload or validate files if needed
}
