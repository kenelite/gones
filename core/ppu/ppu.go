package ppu

import (
	"image"
	"image/color"
)

type PPU struct {
	Cycle    int
	Scanline int
	Frame    int

	VRAM     *VRAM
	OAM      *OAM
	Renderer *Renderer
	Status   *StatusRegister
	Ctrl     *ControlRegister
	Mask     *MaskRegister

	framebuffer [256 * 240]uint8

	NMIOutput bool

	// 新增: PPUADDR/PPUDATA 寄存器相关
	vramAddr  uint16 // 当前 VRAM 地址
	addrLatch bool   // PPUADDR 写入高低字节切换
	oamAddr   byte   // OAMADDR
}

const (
	ScreenWidth  = 256
	ScreenHeight = 240
)

var paletteRGB = [64][3]uint8{
	{84, 84, 84}, {0, 30, 116}, {8, 16, 144}, {48, 0, 136},
	{68, 0, 100}, {92, 0, 48}, {84, 4, 0}, {60, 24, 0},
	{32, 42, 0}, {8, 58, 0}, {0, 64, 0}, {0, 60, 0},
	{0, 50, 60}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0},
	{152, 150, 152}, {8, 76, 196}, {48, 50, 236}, {92, 30, 228},
	{136, 20, 176}, {160, 20, 100}, {152, 34, 32}, {120, 60, 0},
	{84, 90, 0}, {40, 114, 0}, {8, 124, 0}, {0, 118, 40},
	{0, 102, 120}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0},
	{236, 238, 236}, {76, 154, 236}, {120, 124, 236}, {176, 98, 236},
	{228, 84, 236}, {236, 88, 180}, {236, 106, 100}, {212, 136, 32},
	{160, 170, 0}, {116, 196, 0}, {76, 208, 32}, {56, 204, 108},
	{56, 180, 204}, {60, 60, 60}, {0, 0, 0}, {0, 0, 0},
	{236, 238, 236}, {168, 204, 236}, {188, 188, 236}, {212, 178, 236},
	{236, 174, 236}, {236, 174, 212}, {236, 180, 176}, {228, 196, 144},
	{204, 210, 120}, {180, 222, 120}, {168, 226, 144}, {152, 226, 180},
	{160, 214, 228}, {160, 162, 160}, {0, 0, 0}, {0, 0, 0},
}

func NewPPU() *PPU {
	ppu := &PPU{
		VRAM:     NewVRAM(),
		OAM:      NewOAM(),
		Renderer: NewRenderer(),
		Status:   &StatusRegister{},
		Ctrl:     &ControlRegister{},
		Mask:     &MaskRegister{},
	}

	// 初始化默认调色板（调试用）
	// 调色板地址从 0x3F00 开始
	defaultPalette := [32]byte{
		// 背景调色板 (0x3F00-0x3F0F)
		0x0F, 0x01, 0x21, 0x31, // 调色板 0
		0x0F, 0x06, 0x16, 0x26, // 调色板 1
		0x0F, 0x0B, 0x1B, 0x2B, // 调色板 2
		0x0F, 0x10, 0x20, 0x30, // 调色板 3
		// 精灵调色板 (0x3F10-0x3F1F)
		0x0F, 0x01, 0x21, 0x31, // 调色板 0
		0x0F, 0x06, 0x16, 0x26, // 调色板 1
		0x0F, 0x0B, 0x1B, 0x2B, // 调色板 2
		0x0F, 0x10, 0x20, 0x30, // 调色板 3
	}

	// 写入调色板数据到 VRAM
	for i, v := range defaultPalette {
		ppu.VRAM.Palette[i] = v
	}

	// 写入测试 PatternTable（创建一些简单的图案）
	for i := 0; i < 0x2000; i++ {
		ppu.VRAM.PatternTables[i] = byte(i % 256)
	}

	// 写入测试 NameTable（创建棋盘格图案）
	for ty := 0; ty < 30; ty++ {
		for tx := 0; tx < 32; tx++ {
			index := ty*32 + tx
			if (ty/4+tx/4)%2 == 0 {
				ppu.VRAM.NameTables[index] = 0x01 // 使用图案1
			} else {
				ppu.VRAM.NameTables[index] = 0x02 // 使用图案2
			}
		}
	}

	// 写入测试 AttributeTable（设置调色板）
	for i := 0; i < 64; i++ {
		ppu.VRAM.NameTables[960+i] = byte(i % 4) // 使用不同的调色板
	}

	return ppu
}

func (p *PPU) Step() {
	p.Cycle++
	if p.Cycle >= 341 {
		p.Cycle = 0
		p.Scanline++
		if p.Scanline >= 261 {
			p.Scanline = 0
			p.Frame++
			p.Renderer.RenderFrame(p)
		}
	}
}

// GetFrame returns the current framebuffer as an RGBA image.
func (p *PPU) GetFrame() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))

	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			// 从 Renderer 的 Framebuffer 获取颜色索引
			colorIndex := p.Renderer.Framebuffer[y][x] & 0x3F // 6-bit index

			rgb := paletteRGB[colorIndex]
			col := color.RGBA{
				R: rgb[0],
				G: rgb[1],
				B: rgb[2],
				A: 255,
			}
			img.Set(x, y, col)
		}
	}

	return img
}

// GetFrameBuffer returns the current frame as a 2D array of color.RGBA for rendering
func (p *PPU) GetFrameBuffer() [256][240]color.RGBA {
	var buf [256][240]color.RGBA
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			// 从 Renderer 的 Framebuffer 获取颜色索引
			colorIndex := p.Renderer.Framebuffer[y][x] & 0x3F
			rgb := paletteRGB[colorIndex]
			buf[x][y] = color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255}
		}
	}
	return buf
}

func (p *PPU) ClearFrame(colorIndex uint8) {
	for i := range p.framebuffer {
		p.framebuffer[i] = colorIndex & 0x3F
	}
}
