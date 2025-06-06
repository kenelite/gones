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
	switch {
	case addr < 0x2000:
		return b.RAM[addr%0x0800] // 内存镜像
	case addr < 0x4000:
		return b.PPU.ReadRegister(0x2000 + addr%8)
	case addr == 0x4015:
		// APU 状态寄存器
		return 0 // 可扩展
	case addr == 0x4016:
		if b.Controller1 != nil {
			return b.Controller1.Read()
		}
		return 0
	case addr >= 0x8000:
		return b.Cartridge.Read(addr)
	default:
		// 未实现的地址
		return 0
	}
}

func (b *Bus) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		b.RAM[addr%0x0800] = val
	case addr < 0x4000:
		b.PPU.WriteRegister(0x2000+addr%8, val)
	case addr == 0x4015:
		// APU 控制寄存器
	case addr == 0x4016:
		if b.Controller1 != nil {
			b.Controller1.Write(val)
		}
	case addr >= 0x8000:
		b.Cartridge.Write(addr, val)
	}
}

// RunFrame 执行一帧主循环
func (b *Bus) RunFrame() error {
	const cpuCyclesPerFrame = 29780 // NTSC: 1.789773 MHz / 60Hz
	for i := 0; i < cpuCyclesPerFrame; i++ {
		b.CPU.Step()
		b.PPU.Step()
		b.APU.Step()
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
