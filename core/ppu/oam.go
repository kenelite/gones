package ppu

type Sprite struct {
	Y    byte
	Tile byte
	Attr byte
	X    byte
}

type OAM struct {
	Sprites [64]Sprite
}

func NewOAM() *OAM {
	return &OAM{}
}

func (o *OAM) ReadOAMByte(addr byte) byte {
	offset := int(addr)
	if offset < 0 || offset >= 256 {
		return 0
	}
	sprIdx := offset / 4
	field := offset % 4
	spr := o.Sprites[sprIdx]
	switch field {
	case 0:
		return spr.Y
	case 1:
		return spr.Tile
	case 2:
		return spr.Attr
	case 3:
		return spr.X
	}
	return 0
}

func (o *OAM) WriteOAMByte(addr byte, val byte) {
	offset := int(addr)
	if offset < 0 || offset >= 256 {
		return
	}
	sprIdx := offset / 4
	field := offset % 4
	spr := &o.Sprites[sprIdx]
	switch field {
	case 0:
		spr.Y = val
	case 1:
		spr.Tile = val
	case 2:
		spr.Attr = val
	case 3:
		spr.X = val
	}
	println("[gones] OAM Write: addr=", addr, "val=", val)
}
