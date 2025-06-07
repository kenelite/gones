package input

type Controller struct {
	state     [8]bool // 当前按钮状态
	latched   bool
	shiftReg  byte
	readIndex int
}

func NewController() *Controller {
	return &Controller{}
}

// 写入控制位（0x4016）
// latch == 1：锁存当前状态
// latch == 0：允许读取 shift 寄存器
func (c *Controller) Write(val byte) {
	c.latched = val&1 == 1
	if c.latched {
		c.shiftReg = c.serialize()
		c.readIndex = 0
	}
}

// 读取控制器数据（0x4016）
func (c *Controller) Read() byte {
	if c.latched {
		c.shiftReg = c.serialize()
	}
	if c.readIndex >= 8 {
		return 1
	}
	bit := (c.shiftReg >> c.readIndex) & 1
	c.readIndex++
	return bit | 0x40 // NES 高位填 1
}

// 传入新的输入状态（例如按键/手柄）
func (c *Controller) Update(state [8]bool) {
	c.state = state
	if c.latched {
		c.shiftReg = c.serialize()
	}
}

func (c *Controller) serialize() byte {
	var result byte
	for i := 0; i < 8; i++ {
		if c.state[i] {
			result |= 1 << i
		}
	}
	return result
}
