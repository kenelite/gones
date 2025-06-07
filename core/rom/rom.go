package rom

import (
	"fmt"
	"os"
)

type Cartridge struct {
	PRG       []byte // PRG ROM 数据
	CHR       []byte // CHR ROM 数据
	MapperID  byte   // Mapper 编号（本项目初期仅支持 0）
	Mirroring byte   // 镜像模式（0: Horizontal, 1: Vertical）
}

func LoadCartridge(path string) (*Cartridge, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rom: %v", err)
	}

	if len(data) < 16 {
		return nil, fmt.Errorf("invalid rom: too short")
	}

	header := ParseINESHeader(data[:16])
	if string(header.Magic[:]) != "NES\x1A" {
		return nil, fmt.Errorf("invalid iNES header")
	}

	// PRG/CHR 区域
	prgSize := int(header.PRGROMSize) * 0x4000
	chrSize := int(header.CHRROMSize) * 0x2000

	prgStart := 16
	chrStart := prgStart + prgSize

	if len(data) < chrStart+chrSize {
		return nil, fmt.Errorf("invalid rom: file too small")
	}

	return &Cartridge{
		PRG:       data[prgStart:chrStart],
		CHR:       data[chrStart : chrStart+chrSize],
		MapperID:  (header.Flags6 >> 4) | (header.Flags7 & 0xF0),
		Mirroring: header.Flags6 & 0x01,
	}, nil
}

// 只读 PRG 区域
func (c *Cartridge) Read(addr uint16) byte {
	if addr >= 0x8000 {
		offset := int(addr - 0x8000)
		if len(c.PRG) == 0x4000 {
			offset %= 0x4000 // NROM-128 镜像
		}
		return c.PRG[offset]
	}
	return 0
}

func (c *Cartridge) Write(addr uint16, val byte) {
	// Mapper 0: PRG 是只读的
}

// 直接从数据加载 ROM
func LoadCartridgeFromData(data []byte) (*Cartridge, error) {
	fmt.Println("[gones] LoadCartridgeFromData: 开始解析...")
	if len(data) < 16 {
		fmt.Println("[gones] LoadCartridgeFromData: 数据过短")
		return nil, fmt.Errorf("invalid rom: too short")
	}

	header := ParseINESHeader(data[:16])
	fmt.Println("[gones] LoadCartridgeFromData: header 解析完成")
	if string(header.Magic[:]) != "NES\x1A" {
		fmt.Println("[gones] LoadCartridgeFromData: 非法 iNES header")
		return nil, fmt.Errorf("invalid iNES header")
	}

	prgSize := int(header.PRGROMSize) * 0x4000
	chrSize := int(header.CHRROMSize) * 0x2000

	prgStart := 16
	chrStart := prgStart + prgSize

	if len(data) < chrStart+chrSize {
		fmt.Println("[gones] LoadCartridgeFromData: 文件过小")
		return nil, fmt.Errorf("invalid rom: file too small")
	}

	fmt.Println("[gones] LoadCartridgeFromData: 解析 Cartridge 结构体 ...")
	cart := &Cartridge{
		PRG:       data[prgStart:chrStart],
		CHR:       data[chrStart : chrStart+chrSize],
		MapperID:  (header.Flags6 >> 4) | (header.Flags7 & 0xF0),
		Mirroring: header.Flags6 & 0x01,
	}
	fmt.Println("[gones] LoadCartridgeFromData: 完成")
	return cart, nil
}
