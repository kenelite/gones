package ui

import (
	"image/color"
	"log"
	"os"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kenelite/gones/core/bus"
	"golang.org/x/image/font/basicfont"
)

var basicFont = basicfont.Face7x13

// Menu 结构体包含模拟器总线和菜单状态
type Menu struct {
	Bus         *bus.Bus
	Visible     bool
	menuItems   []string
	selectedIdx int
}

// NewMenu 构造函数
func NewMenu(bus *bus.Bus) *Menu {
	return &Menu{
		Bus:       bus,
		Visible:   false,
		menuItems: []string{"Load ROM", "Exit"},
	}
}

// Toggle 显示/隐藏菜单
func (m *Menu) Toggle() {
	m.Visible = !m.Visible
}

// Update 处理菜单输入和操作
func (m *Menu) Update() {
	if !m.Visible {
		return
	}

	// 使用上下键选择菜单项
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		m.selectedIdx = (m.selectedIdx + 1) % len(m.menuItems)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		m.selectedIdx = (m.selectedIdx - 1 + len(m.menuItems)) % len(m.menuItems)
	}

	// 确认键执行菜单项操作
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch m.menuItems[m.selectedIdx] {
		case "Load ROM":
			data, _, err := OpenROM()
			if err != nil {
				log.Println("Failed to load ROM:", err)
			} else {
				m.Bus.LoadROM(data)
				m.Visible = false
			}
		case "Exit":
			// 退出程序
			os.Exit(0)
		}
	}
}

// Draw 绘制菜单
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.Visible {
		return
	}

	const (
		menuWidth  = 200
		menuHeight = 150
		x          = 20
		y          = 20
		itemHeight = 30
	)

	// 半透明背景
	bgColor := color.RGBA{0, 0, 0, 160}
	screen.Fill(bgColor)

	// 菜单背景框
	menuRect := ebiten.NewImage(menuWidth, menuHeight)
	menuRect.Fill(color.RGBA{30, 30, 30, 230})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(menuRect, op)

	// 绘制菜单项
	for i, item := range m.menuItems {
		var col color.Color = color.White
		if i == m.selectedIdx {
			col = color.RGBA{255, 215, 0, 255} // 金色高亮
		}
		text.Draw(screen, item, basicFont, x+20, y+40+i*itemHeight, col)
	}
}
