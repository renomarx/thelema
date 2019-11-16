package ui2d

import (
	"strconv"
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

const CharacteristicColumnLength int32 = 200

func (ui *UI) DrawCharacteristics(p *game.Character, offsetX, offsetH int32) int32 {
	_, h := ui.DrawText(
		"Health : "+strconv.Itoa(p.Health.Current)+"/"+strconv.Itoa(p.Health.Initial)+" (xp: "+strconv.Itoa(p.Health.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Energy : "+strconv.Itoa(p.Energy.Current)+"/"+strconv.Itoa(p.Energy.Initial)+" (xp: "+strconv.Itoa(p.Energy.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)

	offsetH += h
	_, h = ui.DrawText(
		"Strength : "+strconv.Itoa(p.Strength.Current)+" (xp: "+strconv.Itoa(p.Strength.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Will : "+strconv.Itoa(p.Will.Current)+" (xp: "+strconv.Itoa(p.Will.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)
	offsetH += h

	_, h = ui.DrawText(
		"Dexterity : "+strconv.Itoa(p.Dexterity.Current)+" (xp: "+strconv.Itoa(p.Dexterity.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Intelligence : "+strconv.Itoa(p.Intelligence.Current)+" (xp: "+strconv.Itoa(p.Intelligence.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)
	offsetH += h

	_, h = ui.DrawText(
		"Beauty : "+strconv.Itoa(p.Beauty.Current)+" (xp: "+strconv.Itoa(p.Beauty.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Charisma : "+strconv.Itoa(p.Charisma.Current)+" (xp: "+strconv.Itoa(p.Charisma.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)
	offsetH += h

	_, h = ui.DrawText(
		"Speed : "+strconv.Itoa(p.Speed.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res, offsetH)
	_, h = ui.DrawText(
		"Regeneration speed : "+strconv.Itoa(p.RegenerationSpeed.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
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
		offsetH = ui.DrawPowers(offsetH)
	}
}

func (ui *UI) DrawPowers(offsetH int32) int32 {
	p := ui.Game.Level.Player
	var offsetX = int32(PlayerMenuOffsetX * Res)
	_, h := ui.DrawText("Magies (gauche ou droite pour changer)", TextSizeM, ColorGreen, offsetX, offsetH)
	offsetH += h + 10

	powernames := p.GetSortedPowernames()

	for i, powername := range powernames {
		power := p.Powers[powername]
		x := int32((PlayerMenuOffsetX + i) * Res)
		if p.CurrentPower != nil && p.CurrentPower.Type == power.Type {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["ʆ"][0],
				&sdl.Rect{X: x, Y: offsetH, W: Res, H: Res})
		}
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[power.Tile][0],
			&sdl.Rect{X: int32((PlayerMenuOffsetX + i) * Res), Y: offsetH, W: Res, H: Res})
	}
	offsetH += 42
	ui.DrawPower(p.CurrentPower, PlayerMenuOffsetX, offsetH)
	offsetH += 32 + 40

	return offsetH
}

func (ui *UI) DrawPower(power *game.PlayerPower, offsetX, offsetH int32) int32 {
	_, h := ui.DrawText(
		power.Description,
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)

	offsetH += h
	_, h = ui.DrawText(
		"Energy : "+strconv.Itoa(power.Energy),
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	offsetH += h

	return offsetH
}

func (ui *UI) drawCharacterBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["Ʈ"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
