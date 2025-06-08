package bus

import (
	"fmt"

	"github.com/kenelite/gones/core/apu"
	"github.com/kenelite/gones/core/cpu"
	"github.com/kenelite/gones/core/input"
	"github.com/kenelite/gones/core/ppu"
	"github.com/kenelite/gones/core/rom"
)

type Bus struct {
	RAM         [0x0800]byte // 2KB 内存
	CPU         *cpu.CPU
	PPU         *ppu.PPU
	APU         *apu.APU
	Cartridge   *rom.Cartridge
	Controller1 *input.Controller // 新增 Controller1 字段
}

// 实现 CPU.Bus 接口
func (b *Bus) Read(addr uint16) byte {
	// 自动加日志，排查 nil 指针
	switch {
	case addr < 0x2000:
		return b.RAM[addr%0x0800] // 内存镜像
	case addr < 0x4000:
		if b.PPU == nil {
			panic("Bus.PPU is nil in Read")
		}
		return b.PPU.ReadRegister(0x2000 + addr%8)
	case addr == 0x4015:
		// APU 状态寄存器
		return 0 // 可扩展
	case addr == 0x4016:
		if b.Controller1 == nil {
			panic("Bus.Controller1 is nil in Read")
		}
		return b.Controller1.Read()
	case addr >= 0x8000:
		if b.Cartridge == nil {
			panic("Bus.Cartridge is nil in Read")
		}
		return b.Cartridge.Read(addr)
	default:
		// 未实现的地址
		return 0
	}
}

func (b *Bus) Write(addr uint16, val byte) {
	// 自动加日志，排查 nil 指针
	switch {
	case addr < 0x2000:
		b.RAM[addr%0x0800] = val
	case addr < 0x4000:
		if b.PPU == nil {
			panic("Bus.PPU is nil in Write")
		}
		b.PPU.WriteRegister(0x2000+addr%8, val)
	case addr == 0x4015:
		// APU 控制寄存器
	case addr == 0x4016:
		if b.Controller1 == nil {
			panic("Bus.Controller1 is nil in Write")
		}
		b.Controller1.Write(val)
	case addr == 0x4014:
		if b.PPU == nil {
			panic("Bus.PPU is nil in Write (OAMDMA)")
		}
		page := uint16(val) << 8
		println("[gones] OAMDMA: page=", val)
		for i := 0; i < 256; i++ {
			data := b.Read(page + uint16(i))
			b.PPU.OAM.WriteOAMByte(byte(i), data)
			if i < 8 { // 只打印前8项，避免刷屏
				println("[gones] OAMDMA Write: i=", i, "data=", data)
			}
		}
		// DMA 期间 CPU 会卡住 513/514 cycles，简化版可忽略
	case addr >= 0x8000:
		if b.Cartridge == nil {
			panic("Bus.Cartridge is nil in Write")
		}
		b.Cartridge.Write(addr, val)
	}
}

// RunFrame 执行一帧主循环
func (b *Bus) RunFrame() error {
	const ppuCyclesPerFrame = 261 * 341 // 88901
	for i := 0; i < ppuCyclesPerFrame; i++ {
		b.PPU.Step()
		if i%3 == 0 {
			b.CPU.Step()
			b.APU.Step()
		}
	}
	return nil
}

// 加载 ROM 数据
func (b *Bus) LoadROM(data []byte) error {
	fmt.Println("[gones] LoadROM: 开始")
	fmt.Printf("[gones] LoadROM: data 长度 = %d\n", len(data))
	fmt.Println("[gones] LoadROM: 解析 ROM 数据...")
	cartridge, err := rom.LoadCartridgeFromData(data)
	if err != nil {
		fmt.Println("[gones] LoadROM: 解析 ROM 失败")
		return fmt.Errorf("LoadROM: %w", err)
	}
	fmt.Println("[gones] LoadROM: 替换 Cartridge ...")
	b.Cartridge = cartridge
	// 可选：重建 APU
	if b.APU != nil {
		fmt.Println("[gones] LoadROM: 重建 APU ...")
		b.APU = apu.NewAPU()
	}
	fmt.Println("[gones] LoadROM: 写入 CHR 数据到 PPU VRAM ...")
	if b.PPU != nil && b.Cartridge != nil && len(b.Cartridge.CHR) > 0 {
		for i, v := range b.Cartridge.CHR {
			b.PPU.VRAM.PatternTables[i] = v
		}
	}
	fmt.Println("[gones] LoadROM: 完成")
	return nil
}

func NewBus(controller *input.Controller) *Bus {
	fmt.Println("[gones] Bus: 初始化 PPU...")
	ppuObj := ppu.NewPPU()
	fmt.Println("[gones] Bus: 初始化 APU...")
	apuObj := apu.NewAPU()
	fmt.Println("[gones] Bus: 构造 Bus 实例...")
	b := &Bus{
		PPU:         ppuObj,
		APU:         apuObj,
		Controller1: controller,
	}
	fmt.Println("[gones] Bus: 初始化 CPU...")
	b.CPU = cpu.New(b)
	fmt.Println("[gones] Bus: 完成")
	return b
}
