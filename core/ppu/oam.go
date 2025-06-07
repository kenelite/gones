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
