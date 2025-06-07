package cpu

const (
	FlagCarry     = 1 << 0
	FlagZero      = 1 << 1
	FlagInterrupt = 1 << 2
	FlagDecimal   = 1 << 3
	FlagBreak     = 1 << 4
	FlagUnused    = 1 << 5
	FlagOverflow  = 1 << 6
	FlagNegative  = 1 << 7
)

func (c *CPU) setFlag(flag byte, on bool) {
	if on {
		c.Status |= flag
	} else {
		c.Status &^= flag
	}
}

func (c *CPU) getFlag(flag byte) bool {
	return c.Status&flag != 0
}

func (c *CPU) updateZN(value byte) {
	c.setFlag(FlagZero, value == 0)
	c.setFlag(FlagNegative, value&0x80 != 0)
}
