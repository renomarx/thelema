package ui2d

import (
	"sort"
	"strconv"
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawInventory() {
	inventory := ui.Game.Level.Player.Inventory
	if inventory.IsOpen {
		ui.drawInventoryBox()
		var offsetX = int32(PlayerMenuOffsetX*Res + 10)
		var offsetH int32 = 0
		_, h := ui.DrawText("Inventaire", TextSizeL, ColorActive, offsetX, offsetH)
		offsetH += h + 10
		gold := strconv.Itoa(inventory.Gold)
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[game.Gold][0],
			&sdl.Rect{X: int32(offsetX), Y: offsetH, W: Res, H: Res})
		_, h = ui.DrawText(gold, TextSizeM, ColorDisabled, (PlayerMenuOffsetX+1)*Res, offsetH)
		offsetH += h + Res
		_, h = ui.DrawText("Objets spéciaux:", TextSizeM, ColorDisabled, offsetX, offsetH)
		offsetH += h

		runes := make([]string, 0, len(inventory.QuestObjects))
		for r, _ := range inventory.QuestObjects {
			runes = append(runes, string(r))
		}
		sort.Strings(runes) //sort by key
		i := 0
		for _, r := range runes {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex[game.Tile(r)][0],
				&sdl.Rect{X: int32((PlayerMenuOffsetX + i) * Res), Y: offsetH, W: Res, H: Res})
			i++
		}
		offsetH += 64
		_, h = ui.DrawText("Objets communs:", TextSizeM, ColorActive, offsetX, offsetH)
		offsetH += h
		i = 0
		for _, usable := range inventory.Usables {
			x := int32((PlayerMenuOffsetX + i) * Res)
			if usable.Highlighted {
				ui.renderer.Copy(ui.textureAtlas,
					&ui.textureIndex["ʆ"][0],
					&sdl.Rect{X: x, Y: offsetH, W: Res, H: Res})
			}
			o := usable.GetObject()
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex[game.Tile(o.Rune)][0],
				&sdl.Rect{X: int32((PlayerMenuOffsetX + i) * Res), Y: offsetH, W: Res, H: Res})
			i++
		}
	}
}

func (ui *UI) drawInventoryBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["ß"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
