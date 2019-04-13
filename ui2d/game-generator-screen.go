package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
)

const GGScreenOffsetX = 6

func (ui *UI) DrawGameGeneratorScreen() {
	gg := ui.Game.GG
	if gg != nil && gg.IsOpen {
		var offsetH int32 = 0
		_, h := ui.DrawText("Choisissez votre personnage", TextSizeL, ColorActive, GGScreenOffsetX*Res, offsetH)
		offsetH += h + 10
		for i, player := range gg.Players {
			x := GGScreenOffsetX*Res + int32(i*64)
			if gg.IsHighlighted(i) {
				ui.renderer.Copy(ui.textureAtlas,
					&ui.textureIndex['ʆ'][0],
					&sdl.Rect{X: x, Y: offsetH, W: 64, H: 64})
			}
			ui.renderer.Copy(ui.playerTextures[player.Name],
				&sdl.Rect{X: 0, Y: 128, W: 64, H: 64},
				&sdl.Rect{X: x, Y: offsetH, W: 64, H: 64})
		}
		offsetH += 64 + 40

		currentPlayer := gg.GetCurrentPlayer()
		_, h = ui.DrawText(currentPlayer.Name, TextSizeL, ColorActive, GGScreenOffsetX*Res, offsetH)
		offsetH += h + 10
		ui.DrawCharacteristics(&currentPlayer.Character, GGScreenOffsetX, offsetH)
	}
}
