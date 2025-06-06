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
	// TODO: 实现地址映射和镜像逻辑
	return 0
}

func (v *VRAM) Write(addr uint16, val byte) {
	// TODO: 实现地址映射和镜像逻辑
}
