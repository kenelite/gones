package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
)

var (
	LogoImage image.Image
	XboxImage image.Image
)

//go:embed embed/images/*.png
var imageFiles embed.FS

func LoadImages() {
	LogoImage = loadImage("logo.png")
	XboxImage = loadImage("xbox.png")
}

func loadImage(name string) image.Image {
	data, err := imageFiles.ReadFile("embed/images/" + name)
	if err != nil {
		log.Fatalf("Failed to read image file %s: %v", name, err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to decode image %s: %v", name, err)
	}
	return img
}
