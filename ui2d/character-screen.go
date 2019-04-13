package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
	"thelema/game"
)

func (ui *UI) DrawCharacteristics(p *game.Character, offsetX, offsetH int32) int32 {
	_, h := ui.DrawText(
		"Health : "+strconv.Itoa(p.Hitpoints.Current)+"/"+strconv.Itoa(p.Hitpoints.Initial),
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	offsetH += h
	_, h = ui.DrawText(
		"Energy : "+strconv.Itoa(p.Energy.Current)+"/"+strconv.Itoa(p.Energy.Initial),
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	offsetH += h
	_, h = ui.DrawText(
		"Speed : "+strconv.Itoa(p.Speed.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res, offsetH)
	offsetH += h
	_, h = ui.DrawText(
		"Strength : "+strconv.Itoa(p.Strength.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	offsetH += h
	_, h = ui.DrawText(
		"Regeneration speed : "+strconv.Itoa(p.RegenerationSpeed.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	offsetH += h
	return offsetH
}

func (ui *UI) DrawPlayerCharacter() {
	p := ui.Game.Level.Player
	if p.CharacterMenuOpen {
		ui.drawCharacterBox()
		var offsetH int32 = 0
		_, h := ui.DrawText("Personnage", TextSizeL, ColorActive, PlayerMenuOffsetX*Res, offsetH)
		offsetH += h + 10
		offsetH = ui.DrawCharacteristics(&p.Character, PlayerMenuOffsetX, offsetH)
		offsetH += 40
		ui.DrawPowers(offsetH)
	}
}

func (ui *UI) DrawPowers(offsetH int32) {
	p := ui.Game.Level.Player
	var offsetX = int32(PlayerMenuOffsetX * Res)
	_, h := ui.DrawText("Magies", TextSizeL, ColorActive, offsetX, offsetH)
	offsetH += h + 10

	powernames := p.GetSortedPowernames()

	for i, powername := range powernames {
		power := p.Powers[powername]
		x := int32((PlayerMenuOffsetX + i) * Res)
		if p.CurrentPower != nil && p.CurrentPower.Type == power.Type {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['ʆ'][0],
				&sdl.Rect{X: x, Y: offsetH, W: Res, H: Res})
		}
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[power.Tile][0],
			&sdl.Rect{X: int32((PlayerMenuOffsetX + i) * Res), Y: offsetH, W: Res, H: Res})
	}
	offsetH += 32 + 40

	// currentPower := p.CurrentPower
	// _, h = ui.DrawText(currentPower.Type, TextSizeL, ColorActive, PlayerMenuOffsetX*Res, offsetH)
	// offsetH += h + 10
	//ui.DrawPower(currentPower, offsetH)
}

func (ui *UI) drawCharacterBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['Ʈ'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
