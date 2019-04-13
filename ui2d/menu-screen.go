package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawMenu() {
	menu := ui.Game.GetMenu()
	if menu.IsOpen {
		ui.drawMenuBox()
		var offsetH int32 = 0
		for _, choice := range menu.Choices {
			tex := ui.GetTexture(choice.Cmd, TextSizeXL, ColorActive)
			if choice.Highlighted {
				tex.SetColorMod(0, 255, 0)
			} else if choice.Disabled {
				tex.SetColorMod(150, 150, 150)
			} else {
				tex.SetColorMod(255, 255, 255)
			}
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{10, offsetH, w, h})
			offsetH += h
		}
	}
}

func (ui *UI) drawMenuBox() {
	for x := 0; x < ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['ß'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
