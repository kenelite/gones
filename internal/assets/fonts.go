package assets

import (
	"embed"
	"io/fs"
)

//go:embed embed/fonts/*.ttf
var fontFiles embed.FS

func GetFontFS() fs.FS {
	return fontFiles
}

// LoadFonts can be expanded for actual font registration in the UI
func LoadFonts() {
	// Example usage: ui.RegisterFont(fontFiles, "embed/fonts/pixel.ttf")
}
