package cpu

// import "fmt"

// func nilOp(c *CPU) {
// 	// 兼容非法/未定义指令，直接 NOP，不再 panic
// 	// panic(fmt.Sprintf("Unimplemented opcode: %02X", c.fetchByte()-1))
// }

// 所有未实现/非法指令默认映射为 nop，保证兼容性
var instructionTable [256]struct {
	name    string
	cycles  int
	handler func(*CPU)
}

func init() {
	for i := 0; i < 256; i++ {
		instructionTable[i] = struct {
			name    string
			cycles  int
			handler func(*CPU)
		}{"NOP*", 2, nop}
	}
	// 6502常用指令补全（部分寻址模式，便于NES主循环运行）
	instructionTable[0xA9] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"LDA", 2, ldaImmediate}
	instructionTable[0xA5] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"LDA", 3, ldaZeroPage}
	instructionTable[0xAD] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"LDA", 4, ldaAbsolute}
	instructionTable[0x85] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"STA", 3, staZeroPage}
	instructionTable[0x8D] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"STA", 4, staAbsolute}
	instructionTable[0xA2] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"LDX", 2, ldxImmediate}
	instructionTable[0xA0] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"LDY", 2, ldyImmediate}
	instructionTable[0xAA] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TAX", 2, tax}
	instructionTable[0xA8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TAY", 2, tay}
	instructionTable[0xE8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"INX", 2, inx}
	instructionTable[0xC8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"INY", 2, iny}
	instructionTable[0xCA] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"DEX", 2, dex}
	instructionTable[0x88] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"DEY", 2, dey}
	instructionTable[0x8A] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TXA", 2, txa}
	instructionTable[0x98] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TYA", 2, tya}
	instructionTable[0x9A] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TXS", 2, txs}
	instructionTable[0xBA] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"TSX", 2, tsx}
	instructionTable[0x48] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"PHA", 3, pha}
	instructionTable[0x68] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"PLA", 4, pla}
	instructionTable[0x08] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"PHP", 3, php}
	instructionTable[0x28] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"PLP", 4, plp}
	instructionTable[0x00] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"BRK", 7, brk}
	instructionTable[0xEA] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"NOP", 2, nop}
	instructionTable[0x78] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"SEI", 2, sei}
	instructionTable[0x58] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"CLI", 2, cli}
	instructionTable[0x18] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"CLC", 2, clc}
	instructionTable[0x38] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"SEC", 2, sec}
	instructionTable[0xD8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"CLD", 2, cld}
	instructionTable[0xF8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"SED", 2, sed}
	instructionTable[0xB8] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"CLV", 2, clv}
	instructionTable[0x69] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"ADC", 2, adcImmediate}
	instructionTable[0xE9] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"SBC", 2, sbcImmediate}
	instructionTable[0x29] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"AND", 2, andImmediate}
	instructionTable[0x09] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"ORA", 2, oraImmediate}
	instructionTable[0x49] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"EOR", 2, eorImmediate}
	instructionTable[0x4C] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"JMP", 3, jmpAbsolute}
	instructionTable[0x6C] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"JMP", 5, jmpIndirect}
	instructionTable[0x20] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"JSR", 6, jsr}
	instructionTable[0x60] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"RTS", 6, rts}
	instructionTable[0x40] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"RTI", 6, rti}
	instructionTable[0xFA] = struct {
		name    string
		cycles  int
		handler func(*CPU)
	}{"NOP*", 2, nop}
	// ...可继续补充其它常用指令...
}

func ldaImmediate(c *CPU) {
	value := c.fetchByte()
	c.A = value
	c.updateZN(c.A)
}

func brk(c *CPU) {
	c.PC++
	c.pushWord(c.PC)
	c.setFlag(FlagBreak, true)
	c.pushByte(c.Status)
	c.setFlag(FlagInterrupt, true)
	c.PC = uint16(c.Bus.Read(0xFFFE)) | uint16(c.Bus.Read(0xFFFF))<<8
}

func tay(c *CPU) {
	c.Y = c.A
	c.updateZN(c.Y)
}

func ldaZeroPage(c *CPU) {
	addr := c.fetchByte()
	c.A = c.Bus.Read(uint16(addr))
	c.updateZN(c.A)
}

func ldaAbsolute(c *CPU) {
	lo := c.fetchByte()
	hi := c.fetchByte()
	addr := uint16(lo) | uint16(hi)<<8
	c.A = c.Bus.Read(addr)
	c.updateZN(c.A)
}

func staZeroPage(c *CPU) {
	addr := c.fetchByte()
	c.Bus.Write(uint16(addr), c.A)
}

func staAbsolute(c *CPU) {
	lo := c.fetchByte()
	hi := c.fetchByte()
	addr := uint16(lo) | uint16(hi)<<8
	c.Bus.Write(addr, c.A)
}

func ldxImmediate(c *CPU) {
	c.X = c.fetchByte()
	c.updateZN(c.X)
}

func ldyImmediate(c *CPU) {
	c.Y = c.fetchByte()
	c.updateZN(c.Y)
}

func tax(c *CPU) {
	c.X = c.A
	c.updateZN(c.X)
}

func inx(c *CPU) {
	c.X++
	c.updateZN(c.X)
}

func iny(c *CPU) {
	c.Y++
	c.updateZN(c.Y)
}

func dex(c *CPU) {
	c.X--
	c.updateZN(c.X)
}

func dey(c *CPU) {
	c.Y--
	c.updateZN(c.Y)
}

func nop(c *CPU) {}

func sei(c *CPU) {
	c.setFlag(FlagInterrupt, true)
}

func cli(c *CPU) {
	c.setFlag(FlagInterrupt, false)
}

func clc(c *CPU) {
	c.setFlag(FlagCarry, false)
}

func sec(c *CPU) {
	c.setFlag(FlagCarry, true)
}

func cld(c *CPU) {
	c.setFlag(FlagDecimal, false)
}

func sed(c *CPU) {
	c.setFlag(FlagDecimal, true)
}

func clv(c *CPU) {
	c.setFlag(FlagOverflow, false)
}

func adcImmediate(c *CPU) {
	value := c.fetchByte()
	c.adc(value)
}

func sbcImmediate(c *CPU) {
	value := c.fetchByte()
	c.sbc(value)
}

func andImmediate(c *CPU) {
	value := c.fetchByte()
	c.A &= value
	c.updateZN(c.A)
}

func oraImmediate(c *CPU) {
	value := c.fetchByte()
	c.A |= value
	c.updateZN(c.A)
}

func eorImmediate(c *CPU) {
	value := c.fetchByte()
	c.A ^= value
	c.updateZN(c.A)
}

func jmpAbsolute(c *CPU) {
	lo := c.fetchByte()
	hi := c.fetchByte()
	c.PC = uint16(lo) | uint16(hi)<<8
}

func jmpIndirect(c *CPU) {
	lo := c.fetchByte()
	hi := c.fetchByte()
	ptr := uint16(lo) | uint16(hi)<<8
	c.PC = uint16(c.Bus.Read(ptr)) | uint16(c.Bus.Read((ptr&0xFF00)|((ptr+1)&0x00FF)))<<8
}

func jsr(c *CPU) {
	lo := c.fetchByte()
	hi := c.fetchByte()
	addr := uint16(lo) | uint16(hi)<<8
	c.pushWord(c.PC - 1)
	c.PC = addr
}

func rts(c *CPU) {
	lo := c.popByte()
	hi := c.popByte()
	c.PC = (uint16(lo) | uint16(hi)<<8) + 1
}

func rti(c *CPU) {
	c.Status = c.popByte()
	lo := c.popByte()
	hi := c.popByte()
	c.PC = uint16(lo) | uint16(hi)<<8
}

// 新增常用指令handler
func txa(c *CPU) { c.A = c.X; c.updateZN(c.A) }
func tya(c *CPU) { c.A = c.Y; c.updateZN(c.A) }
func txs(c *CPU) { c.SP = c.X }
func tsx(c *CPU) { c.X = c.SP; c.updateZN(c.X) }
func pha(c *CPU) { c.pushByte(c.A) }
func pla(c *CPU) { c.A = c.popByte(); c.updateZN(c.A) }
func php(c *CPU) { c.pushByte(c.Status | 0x10) }
func plp(c *CPU) { c.Status = c.popByte() }
