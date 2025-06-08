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
	case 0x2007: // PPUDATA (VRAM 读)
		return p.readPPUData()
	case 0x2004: // OAMDATA
		return p.readOAMData()
	}
	return 0
}

func (p *PPU) WriteRegister(addr uint16, value byte) {
	switch addr {
	case 0x2000: // PPUCTRL
		p.Ctrl.Reg = value
	case 0x2001: // PPUMASK
		p.Mask.Reg = value
	case 0x2003: // OAMADDR
		p.oamAddr = value
	case 0x2004: // OAMDATA
		p.writeOAMData(value)
	case 0x2005: // PPUSCROLL
		p.writePPUScroll(value)
	case 0x2006: // PPUADDR
		p.writePPUAddr(value)
	case 0x2007: // PPUDATA (VRAM 写)
		p.writePPUData(value)
	case 0x4014: // OAMDMA
		// OAM DMA 由 Bus 实现（CPU 写 $4014 时，Bus 负责将 $XX00-$XXFF 256 字节写入 PPU.OAM）
	}
}

// --- PPU 寄存器/VRAM 读写辅助实现 ---
func (p *PPU) readPPUData() byte {
	// 仅简化实现：直接读 VRAM
	return p.VRAM.Read(p.vramAddr)
}

func (p *PPU) writePPUData(val byte) {
	// 仅简化实现：直接写 VRAM
	p.VRAM.Write(p.vramAddr, val)
	p.vramAddr++
}

func (p *PPU) writePPUAddr(val byte) {
	// 仅简化实现：高低字节交替写入
	if !p.addrLatch {
		p.vramAddr = uint16(val) << 8
		p.addrLatch = true
	} else {
		p.vramAddr |= uint16(val)
		p.addrLatch = false
	}
}

func (p *PPU) writePPUScroll(val byte) {
	// 可根据需要实现
}

func (p *PPU) readOAMData() byte {
	val := p.OAM.ReadOAMByte(p.oamAddr)
	p.oamAddr++
	return val
}

func (p *PPU) writeOAMData(val byte) {
	p.OAM.WriteOAMByte(p.oamAddr, val)
	p.oamAddr++
}
