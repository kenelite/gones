package ppu

type VRAM struct {
	PatternTables [0x2000]byte // 0x0000-0x1FFF
	NameTables    [0x1000]byte // 0x2000-0x2FFF (including mirroring)
	Palette       [32]byte     // 0x3F00-0x3F1F
}

func NewVRAM() *VRAM {
	return &VRAM{}
}

func (v *VRAM) Read(addr uint16) byte {
	// 实现 Pattern Table、Name Table、Palette 区域的读取
	switch {
	case addr < 0x2000:
		return v.PatternTables[addr]
	case addr >= 0x2000 && addr < 0x3000:
		return v.NameTables[addr-0x2000]
	case addr >= 0x3F00 && addr < 0x3F20:
		return v.Palette[addr-0x3F00]
	default:
		return 0
	}
}

func (v *VRAM) Write(addr uint16, val byte) {
	// 实现 Pattern Table、Name Table、Palette 区域的写入
	switch {
	case addr < 0x2000:
		v.PatternTables[addr] = val
	case addr >= 0x2000 && addr < 0x3000:
		v.NameTables[addr-0x2000] = val
	case addr >= 0x3F00 && addr < 0x3F20:
		v.Palette[addr-0x3F00] = val
	}
}
