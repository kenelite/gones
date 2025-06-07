package cpu

import "fmt"

type CPU struct {
	A, X, Y byte   // 累加器、X寄存器、Y寄存器
	PC      uint16 // 程序计数器
	SP      byte   // 栈指针
	Status  byte   // 状态寄存器

	Bus    Bus // 与总线通信接口
	cycles int // 剩余的周期数
}

// Bus 是 CPU 与内存/PPU 等组件交互的接口
type Bus interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}

// New returns a new CPU instance
func New(bus Bus) *CPU {
	return &CPU{
		SP:     0xFD, // 初始栈指针
		Status: 0x24,
		Bus:    bus,
	}
}

// Step 执行一个指令周期
func (c *CPU) Step() {
	if c.cycles == 0 {
		opcode := c.fetchByte()
		instruction := instructionTable[opcode]
		c.cycles = instruction.cycles

		if instruction.handler == nil {
			panic(fmt.Sprintf("instructionTable[%02X].handler is nil", opcode))
		}
		instruction.handler(c)
	}

	c.cycles--
}

// fetchByte 从当前 PC 获取一个字节并前进
func (c *CPU) fetchByte() byte {
	data := c.Bus.Read(c.PC)
	c.PC++
	return data
}

// pushByte 将一个字节压入栈
func (c *CPU) pushByte(val byte) {
	c.Bus.Write(0x0100+uint16(c.SP), val) // 栈空间从0x0100开始
	c.SP--
}

// pushWord 将一个16位值低字节先压入栈，高字节后压入栈（6502栈是小端）
func (c *CPU) pushWord(val uint16) {
	high := byte(val >> 8)
	low := byte(val & 0xff)
	c.pushByte(high)
	c.pushByte(low)
}

// popByte 从栈顶弹出一个字节
func (c *CPU) popByte() byte {
	c.SP++
	return c.Bus.Read(0x0100 + uint16(c.SP))
}

// adc 实现加法并设置标志位（简化版，建议后续完善）
func (c *CPU) adc(value byte) {
	carry := byte(0)
	if c.getFlag(FlagCarry) {
		carry = 1
	}
	sum := uint16(c.A) + uint16(value) + uint16(carry)
	c.setFlag(FlagCarry, sum > 0xFF)
	c.setFlag(FlagZero, byte(sum) == 0)
	c.setFlag(FlagOverflow, (^uint16(c.A^value)&(uint16(c.A)^sum)&0x80) != 0)
	c.setFlag(FlagNegative, (sum&0x80) != 0)
	c.A = byte(sum)
}

// sbc 实现减法并设置标志位（简化版，建议后续完善）
func (c *CPU) sbc(value byte) {
	carry := byte(0)
	if c.getFlag(FlagCarry) {
		carry = 1
	}
	diff := int16(c.A) - int16(value) - int16(1-carry)
	c.setFlag(FlagCarry, diff >= 0)
	c.setFlag(FlagZero, byte(diff) == 0)
	c.setFlag(FlagOverflow, ((int16(c.A)^diff)&(^int16(value)^diff)&0x80) != 0)
	c.setFlag(FlagNegative, (diff&0x80) != 0)
	c.A = byte(diff)
}
