package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
)

const PlayerMenuOffsetX = 6

func (ui *UI) DrawPlayerMenu() {
	menu := ui.Game.Level.Player.Menu
	if menu.IsOpen {
		ui.drawPlayerMenuBox()
		var offsetH int32 = 0
		for _, choice := range menu.Choices {
			tex := ui.GetTexture(choice.Cmd, TextSizeL, ColorActive)
			if choice.Highlighted {
				tex.SetColorMod(0, 255, 0)
			} else {
				tex.SetColorMod(255, 255, 255)
			}
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{10, offsetH, w, h})
			offsetH += h
		}
		ui.DrawInventory()
		ui.DrawLibrary()
		ui.DrawQuests()
		ui.DrawPlayerCharacter()
		ui.DrawMap()
	}
}

func (ui *UI) drawPlayerMenuBox() {
	for x := 0; x < PlayerMenuOffsetX; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["Ʈ"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
