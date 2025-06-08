package ui

import (
	"fmt"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/kenelite/gones/core/ppu"
)

type Screen struct {
	image *ebiten.Image
}

func NewScreen() *Screen {
	img := ebiten.NewImage(256, 240)
	return &Screen{image: img}
}

func (s *Screen) Render(ppu *ppu.PPU, target *ebiten.Image) {
	fmt.Println("[gones] Screen.Render: begin")
	// 获取 PPU 帧缓冲区
	frame := ppu.GetFrameBuffer() // [256][240]color.RGBA
	// 创建一个 256x240 的 ebiten.Image
	img := ebiten.NewImage(256, 240)
	count := 0
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			img.Set(x, y, frame[x][y])
			c := frame[x][y]
			if c.R != 0 || c.G != 0 || c.B != 0 {
				count++
			}
		}
	}
	fmt.Printf("[gones] 非黑色像素数量: %d\n", count)
	// 放大 2x 绘制到 screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	target.DrawImage(img, op)
	fmt.Println("[gones] Screen.Render: end")
}
