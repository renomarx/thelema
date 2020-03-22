package uipixel

import "github.com/faiface/pixel"

const GGScreenOffsetX float64 = 200

func (ui *UI) DrawGameGeneratorScreen() {
	gg := ui.Game.GG
	if gg != nil && gg.IsOpen {
		var offsetH float64 = 0
		_, h := ui.DrawText("Choisissez votre personnage", TextSizeL, ColorActive, GGScreenOffsetX, offsetH)
		offsetH += h + 10
		for i, player := range gg.Players {
			x := GGScreenOffsetX + float64(i*64)
			if gg.IsHighlighted(i) {
				bkSprite := pixel.NewSprite(ui.textureAtlas, ui.textureIndex["ʆ"][0])
				ui.drawSprite(bkSprite, float64(x), offsetH)
			}
			sprite := pixel.NewSprite(ui.playerTextures[player.Name], pixel.R(0, 128, 64, 64))
			ui.drawSprite(sprite, float64(x), offsetH)
		}
		offsetH += 64 + 40

		currentPlayer := gg.GetCurrentPlayer()
		_, h = ui.DrawText("Affinity:  "+currentPlayer.Affinity, TextSizeL, ColorActive, GGScreenOffsetX, offsetH)
		// offsetH += h + 10
		// ui.DrawCharacteristics(&currentPlayer.Character, GGScreenOffsetX, offsetH)
	}
}
