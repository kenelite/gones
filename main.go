package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/kenelite/gones/core/audio"
	"github.com/kenelite/gones/core/bus"
	"github.com/kenelite/gones/core/input"
	"github.com/kenelite/gones/internal/ui"
	"log"
)

func main() {
	audio.InitAudio()
	controller := input.NewController()

	b := bus.NewBus(controller)
	g := ui.NewGame(b)

	// 加载 ROM
	data, _, err := ui.OpenROM()
	if err != nil {
		log.Fatal("No ROM loaded")
	}
	b.LoadROM(data)

	// 启动主循环
	ebiten.SetWindowTitle("Gones - NES Emulator")
	ebiten.SetWindowSize(512, 480) // 256x240 x2
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
