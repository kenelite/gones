package ppu

import "fmt"

type Renderer struct {
	Framebuffer [240][256]byte // 每帧图像缓冲区
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// 渲染一帧：实现完整 NES 背景+精灵渲染（简化，支持 8x8 sprite，无优先级/遮挡/翻转/透明等）
func (r *Renderer) RenderFrame(ppu *PPU) {
	// 清空帧缓冲区
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			r.Framebuffer[y][x] = 0x0F // 默认背景色
		}
	}

	// 背景渲染
	nametableBase := 0x2000
	patternTableBase := 0x0000
	attributeTableBase := nametableBase + 0x3C0

	for ty := 0; ty < 30; ty++ {
		for tx := 0; tx < 32; tx++ {
			tileIndex := ppu.VRAM.Read(uint16(nametableBase + ty*32 + tx))

			for row := 0; row < 8; row++ {
				low := ppu.VRAM.Read(uint16(patternTableBase + int(tileIndex)*16 + row))
				high := ppu.VRAM.Read(uint16(patternTableBase + int(tileIndex)*16 + row + 8))

				// 计算属性表索引
				attrX := tx / 4
				attrY := ty / 4
				attrByte := ppu.VRAM.Read(uint16(attributeTableBase + attrY*8 + attrX))

				// 计算调色板索引
				shift := ((ty%4)/2)*4 + ((tx%4)/2)*2
				paletteIndex := (attrByte >> shift) & 0x03

				for col := 0; col < 8; col++ {
					bit := 7 - col
					pixel := ((high>>bit)&1)<<1 | ((low >> bit) & 1)

					if pixel > 0 { // 只渲染非透明像素
						// 背景调色板从 0x3F00 开始
						colorIndex := ppu.VRAM.Read(0x3F00+uint16(paletteIndex)*4+uint16(pixel)) & 0x3F
						xPix := tx*8 + col
						yPix := ty*8 + row
						if xPix < 256 && yPix < 240 {
							r.Framebuffer[yPix][xPix] = colorIndex
						}
					}
				}
			}
		}
	}

	// 精灵渲染（OAM，8x8 sprite，优先级/遮挡/翻转/透明等未实现）
	for i := 0; i < 64; i++ {
		spr := ppu.OAM.Sprites[i]
		y := int(spr.Y)
		tileIndex := spr.Tile
		attr := spr.Attr
		x := int(spr.X)
		paletteIndex := (attr & 0x3) + 4 // sprite palette 起始于 0x3F10

		for row := 0; row < 8; row++ {
			low := ppu.VRAM.Read(uint16(0x0000 + int(tileIndex)*16 + row))
			high := ppu.VRAM.Read(uint16(0x0000 + int(tileIndex)*16 + row + 8))

			for col := 0; col < 8; col++ {
				bit := 7 - col
				pixel := ((high>>bit)&1)<<1 | ((low >> bit) & 1)

				if pixel == 0 {
					continue // 透明
				}

				// 精灵调色板从 0x3F10 开始
				colorIndex := ppu.VRAM.Read(0x3F10+uint16((paletteIndex-4)*4)+uint16(pixel)) & 0x3F
				xPix := x + col
				yPix := y + row

				if xPix >= 0 && xPix < 256 && yPix >= 0 && yPix < 240 {
					r.Framebuffer[yPix][xPix] = colorIndex
				}
			}
		}
	}

	// 只在调试模式下显示信息
	if ppu.Frame%60 == 0 { // 每秒显示一次
		fmt.Printf("[gones] 渲染帧 %d: 非黑色像素数量 = %d\n", ppu.Frame, countNonBlackPixels(r.Framebuffer))
	}
}

// 计算非黑色像素数量
func countNonBlackPixels(fb [240][256]byte) int {
	count := 0
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			if fb[y][x] != 0x0F { // 0x0F 是背景色
				count++
			}
		}
	}
	return count
}
