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

	// 调试：显示一些NameTable数据
	if ppu.Frame%60 == 0 { // 每秒显示一次
		fmt.Printf("[gones] 渲染帧 %d: NameTable[0-15] = ", ppu.Frame)
		for i := 0; i < 16; i++ {
			fmt.Printf("%d ", ppu.VRAM.Read(uint16(nametableBase+i)))
		}
		fmt.Println()

		// 显示一些PatternTable数据
		fmt.Printf("[gones] 渲染帧 %d: PatternTable[0-15] = ", ppu.Frame)
		for i := 0; i < 16; i++ {
			fmt.Printf("%d ", ppu.VRAM.Read(uint16(patternTableBase+i)))
		}
		fmt.Println()

		// 显示调色板数据
		fmt.Printf("[gones] 渲染帧 %d: Palette[0-15] = ", ppu.Frame)
		for i := 0; i < 16; i++ {
			fmt.Printf("%d ", ppu.VRAM.Palette[i])
		}
		fmt.Println()
	}

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
				shift := ((tx%4)/2)*2 + ((ty%4)/2)*4
				paletteIndex := (attrByte >> shift) & 0x03

				for col := 0; col < 8; col++ {
					bit0 := (low >> (7 - col)) & 0x01
					bit1 := (high >> (7 - col)) & 0x01
					colorIndex := bit1<<1 | bit0

					if colorIndex != 0 { // 透明色
						finalColorIndex := uint8(uint8(paletteIndex)*4 + uint8(colorIndex))
						yPos := ty*8 + row
						xPos := tx*8 + col
						if yPos < 240 && xPos < 256 {
							r.Framebuffer[yPos][xPos] = finalColorIndex
						}
					}
				}
			}
		}
	}

	// 精灵渲染（简化版）
	for i := 0; i < 64; i++ {
		y := ppu.OAM.ReadOAMByte(byte(i * 4))
		tileIndex := ppu.OAM.ReadOAMByte(byte(i*4 + 1))
		attributes := ppu.OAM.ReadOAMByte(byte(i*4 + 2))
		x := ppu.OAM.ReadOAMByte(byte(i*4 + 3))

		if y < 240 && tileIndex > 0 {
			paletteIndex := (attributes >> 1) & 0x03

			for row := 0; row < 8; row++ {
				low := ppu.VRAM.Read(uint16(patternTableBase + int(tileIndex)*16 + row))
				high := ppu.VRAM.Read(uint16(patternTableBase + int(tileIndex)*16 + row + 8))

				for col := 0; col < 8; col++ {
					bit0 := (low >> (7 - col)) & 0x01
					bit1 := (high >> (7 - col)) & 0x01
					colorIndex := bit1<<1 | bit0

					if colorIndex != 0 {
						finalColorIndex := uint8(uint8(paletteIndex)*4 + uint8(colorIndex) + 16) // +16 for sprite palettes
						yPos := int(y) + row
						xPos := int(x) + col
						if yPos < 240 && xPos < 256 {
							r.Framebuffer[yPos][xPos] = finalColorIndex
						}
					}
				}
			}
		}
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
