package ui

import (
	"fmt"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/kenelite/gones/core/bus"
	"github.com/kenelite/gones/core/input"
)

type Game struct {
	Bus   *bus.Bus
	Frame *Screen
	Menu  *Menu
}

func NewGame(bus *bus.Bus) *Game {
	return &Game{
		Bus:   bus,
		Frame: NewScreen(),
		Menu:  NewMenu(bus),
	}
}

func (g *Game) Update() error {
	// 按 ESC 切换菜单
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.Menu.Toggle()
	}

	if g.Menu.Visible {
		g.Menu.Update()
		return nil
	}

	// 普通游戏逻辑
	g.Bus.Controller1.Update(input.PollInput())
	return g.Bus.RunFrame()
}

func (g *Game) Draw(screen *ebiten.Image) {
	fmt.Println("[gones] Draw called")
	fmt.Println("[gones] Frame.Render 调用")
	g.Frame.Render(g.Bus.PPU, screen)
	fmt.Println("[gones] Menu.Draw 调用")
	g.Menu.Draw(screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 256 * 2, 240 * 2 // 默认放大 2 倍
}

func formatFPS() string {
	return fmt.Sprintf("%.0f", ebiten.CurrentFPS())
}
