package rom

type INESHeader struct {
	Magic      [4]byte // "NES\x1A"
	PRGROMSize byte    // PRG ROM 单位：16KB
	CHRROMSize byte    // CHR ROM 单位：8KB
	Flags6     byte
	Flags7     byte
	Flags8     byte
	Flags9     byte
	Flags10    byte
	Reserved   [5]byte
}

func ParseINESHeader(data []byte) *INESHeader {
	return &INESHeader{
		Magic:      [4]byte{data[0], data[1], data[2], data[3]},
		PRGROMSize: data[4],
		CHRROMSize: data[5],
		Flags6:     data[6],
		Flags7:     data[7],
		Flags8:     data[8],
		Flags9:     data[9],
		Flags10:    data[10],
		Reserved:   [5]byte{data[11], data[12], data[13], data[14], data[15]},
	}
}
