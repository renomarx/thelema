package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
)

const TextSizeXXL = 24
const TextSizeXL = 20
const TextSizeL = 16
const TextSizeM = 14
const TextSizeS = 12
const TextSizeXS = 10

func (ui *UI) DrawText(text string, size int, color sdl.Color, offsetX, offsetY int32) (int32, int32) {
	tex := ui.GetTexture(text, size, color)
	_, _, w, h, _ := tex.Query()
	ui.renderer.Copy(tex, nil, &sdl.Rect{offsetX, offsetY, w, h})

	return w, h
}
