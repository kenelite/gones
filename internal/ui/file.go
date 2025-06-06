package ui

import (
	"github.com/sqweek/dialog"
	"os"
	"path/filepath"
)

func OpenROM() ([]byte, string, error) {
	filename, err := dialog.File().Filter("NES Files", "nes").Title("Select a ROM").Load()
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(filename)
	return data, filepath.Base(filename), err
}
