package ppu

type ControlRegister struct {
	Reg byte
}

type MaskRegister struct {
	Reg byte
}

type StatusRegister struct {
	Reg byte
}

func (p *PPU) ReadRegister(addr uint16) byte {
	switch addr {
	case 0x2002: // PPUSTATUS
		return p.Status.Reg
	}
	return 0
}

func (p *PPU) WriteRegister(addr uint16, value byte) {
	switch addr {
	case 0x2000: // PPUCTRL
		p.Ctrl.Reg = value
	}
}
