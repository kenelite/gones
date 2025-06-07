package ui

import (
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
	// 获取 PPU framebuffer（RGBA 图像）
	frame := ppu.GetFrame()

	// 填充 image 数据
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			col := frame.At(x, y)
			s.image.Set(x, y, col)
		}
	}

	// 放大 2x 并绘制到 screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	target.DrawImage(s.image, op)
}
