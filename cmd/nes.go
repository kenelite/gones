package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/kenelite/gones/core/bus"
	"github.com/kenelite/gones/core/input"
	"github.com/kenelite/gones/internal/ui"
)

func main() {
	// 命令行参数支持
	romPath := flag.String("rom", "", "Path to NES ROM file (.nes)")
	windowScale := flag.Int("scale", 2, "Window scale factor (default 2)")
	flag.Parse()

	if *romPath == "" {
		fmt.Println("Usage: nes -rom <path_to_rom_file.nes> [-scale 2]")
		os.Exit(1)
	}

	// 读取 ROM 文件
	romData, err := os.ReadFile(*romPath)
	if err != nil {
		log.Fatalf("Failed to read ROM file: %v", err)
	}

	// 初始化控制器输入
	controller := input.NewController()

	// 初始化总线，载入 ROM
	bus := bus.NewBus(controller)
	if err := bus.LoadROM(romData); err != nil {
		log.Fatalf("Failed to load ROM into bus: %v", err)
	}

	// 初始化 UI Game 对象
	game := ui.NewGame(bus)

	// 设置窗口标题和尺寸
	ebiten.SetWindowTitle("Gones - NES Emulator")
	ebiten.SetWindowSize(256**windowScale, 240**windowScale)

	// 运行游戏主循环
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game loop exited with error: %v", err)
	}
}
