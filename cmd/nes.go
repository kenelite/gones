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
	fmt.Println("[gones] 启动参数解析...")
	romPath := flag.String("rom", "", "Path to NES ROM file (.nes)")
	windowScale := flag.Int("scale", 2, "Window scale factor (default 2)")
	flag.Parse()

	var romData []byte
	var err error

	if *romPath == "" {
		fmt.Println("[gones] 未指定 -rom 参数，弹窗选择 ROM 文件...")
		romData, _, err = ui.OpenROM()
		if err != nil {
			log.Fatalf("No ROM loaded: %v", err)
		}
	} else {
		fmt.Printf("[gones] 通过命令行参数加载 ROM: %s\n", *romPath)
		romData, err = os.ReadFile(*romPath)
		if err != nil {
			log.Fatalf("Failed to read ROM file: %v", err)
		}
	}

	fmt.Println("[gones] 初始化控制器...")
	controller := input.NewController()

	fmt.Println("[gones] 初始化总线...")
	bus := bus.NewBus(controller)
	fmt.Println("[gones] 加载 ROM 到总线...")
	if err := bus.LoadROM(romData); err != nil {
		log.Fatalf("Failed to load ROM into bus: %v", err)
	}

	fmt.Println("[gones] 初始化 UI Game...")
	game := ui.NewGame(bus)

	fmt.Println("[gones] 设置窗口参数...")
	ebiten.SetWindowTitle("Gones - NES Emulator")
	ebiten.SetWindowSize(256**windowScale, 240**windowScale)

	fmt.Println("[gones] 启动主循环...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game loop exited with error: %v", err)
	}
}
