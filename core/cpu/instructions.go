package cpu

type instruction struct {
	name    string
	cycles  int
	handler func(*CPU)
}

var instructionTable = [256]instruction{
	0xA9: {"LDA", 2, ldaImmediate},
	0x00: {"BRK", 7, brk},
	// ... 添加更多指令
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
