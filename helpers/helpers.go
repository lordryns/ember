package helpers

import (
	"ember/engine"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateProject(path string, name string, config *engine.GameConfig) error {
	var fullPath = filepath.Join(path, name)

	var _, statErr = os.Stat(fullPath)
	if !os.IsNotExist(statErr) {
		return errors.New("A folder with this name already exists in this location")
	}

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return err
	}

	config.Title = name
	var configBytes, confErr = json.Marshal(config)
	if confErr != nil {
		return confErr
	}

	var configFile, err = os.Create(filepath.Join(fullPath, "ember.json"))
	if err != nil {
		return fmt.Errorf("%s\nError occured during the creation of the %v.ember file", err, name)
	}

	configFile.Write(configBytes)

	defer configFile.Close()

	return nil
}

func IsValidProject(path string) error {
	var configPath = filepath.Join(path, "ember.json")
	var _, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return errors.New("Not a valid ember project!")
	}

	return nil
}

// got this func from chatgpt, works okay so...yeah
func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	return fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)
}

func CovertToInt(s string) int {
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return int(num)
	}

	return 0
}

// this is being used in main
type RightClickLabel struct {
	widget.Label
	OnRightClick func()
}

func NewRightClickLabel(text string) *RightClickLabel {
	r := &RightClickLabel{}
	r.ExtendBaseWidget(r)
	r.SetText(text)
	return r
}

// Left click (List selection still works with this in place)
func (r *RightClickLabel) Tapped(*fyne.PointEvent) {}

// Right click
func (r *RightClickLabel) TappedSecondary(*fyne.PointEvent) {
	if r.OnRightClick != nil {
		r.OnRightClick()
	}
}
