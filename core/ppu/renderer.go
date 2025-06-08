package ppu

type Renderer struct {
	Framebuffer [240][256]byte // 每帧图像缓冲区
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// 渲染一帧：实现完整 NES 背景+精灵渲染（简化，支持 8x8 sprite，无优先级/遮挡/翻转/透明等）
func (r *Renderer) RenderFrame(ppu *PPU) {
	println("[gones] Renderer.RenderFrame: called, ppu.Frame=", ppu.Frame)
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
				attrX := tx / 4
				attrY := ty / 4
				attrByte := ppu.VRAM.Read(uint16(attributeTableBase + attrY*8 + attrX))
				shift := ((ty%4)/2)*4 + ((tx%4)/2)*2
				paletteIndex := (attrByte >> shift) & 0x03
				for col := 0; col < 8; col++ {
					bit := 7 - col
					pixel := ((high>>bit)&1)<<1 | ((low >> bit) & 1)
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
				colorIndex := ppu.VRAM.Read(0x3F10+uint16((paletteIndex-4)*4)+uint16(pixel)) & 0x3F
				xPix := x + col
				yPix := y + row
				if xPix < 256 && yPix < 240 {
					r.Framebuffer[yPix][xPix] = colorIndex
				}
			}
		}
	}
	println("[gones] Renderer.RenderFrame: VRAM patternTable[0:16]=",
		ppu.VRAM.Read(0x0000), ppu.VRAM.Read(0x0001), ppu.VRAM.Read(0x0002), ppu.VRAM.Read(0x0003),
		ppu.VRAM.Read(0x0004), ppu.VRAM.Read(0x0005), ppu.VRAM.Read(0x0006), ppu.VRAM.Read(0x0007),
		ppu.VRAM.Read(0x0008), ppu.VRAM.Read(0x0009), ppu.VRAM.Read(0x000A), ppu.VRAM.Read(0x000B),
		ppu.VRAM.Read(0x000C), ppu.VRAM.Read(0x000D), ppu.VRAM.Read(0x000E), ppu.VRAM.Read(0x000F))
	println("[gones] Renderer.RenderFrame: Palette[0x3F00:0x3F08]=",
		ppu.VRAM.Read(0x3F00), ppu.VRAM.Read(0x3F01), ppu.VRAM.Read(0x3F02), ppu.VRAM.Read(0x3F03),
		ppu.VRAM.Read(0x3F04), ppu.VRAM.Read(0x3F05), ppu.VRAM.Read(0x3F06), ppu.VRAM.Read(0x3F07))
	println("[gones] Renderer.RenderFrame: NameTable[0:16]=",
		ppu.VRAM.Read(0x2000), ppu.VRAM.Read(0x2001), ppu.VRAM.Read(0x2002), ppu.VRAM.Read(0x2003),
		ppu.VRAM.Read(0x2004), ppu.VRAM.Read(0x2005), ppu.VRAM.Read(0x2006), ppu.VRAM.Read(0x2007),
		ppu.VRAM.Read(0x2008), ppu.VRAM.Read(0x2009), ppu.VRAM.Read(0x200A), ppu.VRAM.Read(0x200B),
		ppu.VRAM.Read(0x200C), ppu.VRAM.Read(0x200D), ppu.VRAM.Read(0x200E), ppu.VRAM.Read(0x200F))
	println("[gones] Renderer.RenderFrame: Framebuffer[0:16]=",
		r.Framebuffer[0][0], r.Framebuffer[0][1], r.Framebuffer[0][2], r.Framebuffer[0][3],
		r.Framebuffer[0][4], r.Framebuffer[0][5], r.Framebuffer[0][6], r.Framebuffer[0][7],
		r.Framebuffer[0][8], r.Framebuffer[0][9], r.Framebuffer[0][10], r.Framebuffer[0][11],
		r.Framebuffer[0][12], r.Framebuffer[0][13], r.Framebuffer[0][14], r.Framebuffer[0][15])
	println("[gones] Renderer.RenderFrame: finished")
}
