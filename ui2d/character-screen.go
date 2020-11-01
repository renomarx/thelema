package ui2d

import (
	"fmt"
	"strconv"
	"github.com/renomarx/thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

const CharacteristicColumnLength int32 = 200

func (ui *UI) DrawCharacteristics(p *game.Character, offsetX, offsetH int32) int32 {
	_, h := ui.DrawText(
		"Santé : "+strconv.Itoa(p.Health.Current)+"/"+strconv.Itoa(p.Health.Initial)+" (xp: "+strconv.Itoa(p.Health.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Energie : "+strconv.Itoa(p.Energy.Current)+"/"+strconv.Itoa(p.Energy.Initial)+" (xp: "+strconv.Itoa(p.Energy.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)

	offsetH += h
	_, h = ui.DrawText(
		"Force : "+strconv.Itoa(p.Strength.Current)+" (xp: "+strconv.Itoa(p.Strength.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Volonté : "+strconv.Itoa(p.Will.Current)+" (xp: "+strconv.Itoa(p.Will.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)
	offsetH += h

	_, h = ui.DrawText(
		"Dextérité : "+strconv.Itoa(p.Dexterity.Current)+" (xp: "+strconv.Itoa(p.Dexterity.Xp)+")",
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
		"Beauté : "+strconv.Itoa(p.Beauty.Current)+" (xp: "+strconv.Itoa(p.Beauty.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)
	_, h = ui.DrawText(
		"Charisme : "+strconv.Itoa(p.Charisma.Current)+" (xp: "+strconv.Itoa(p.Charisma.Xp)+")",
		TextSizeM,
		ColorActive,
		offsetX*Res+CharacteristicColumnLength,
		offsetH)
	offsetH += h

	_, h = ui.DrawText(
		"Vitesse : "+strconv.Itoa(p.Speed.Current),
		TextSizeM,
		ColorActive,
		offsetX*Res, offsetH)
	_, h = ui.DrawText(
		"Vitesse de régéneration : "+strconv.Itoa(p.RegenerationSpeed.Current),
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
	_, h := ui.DrawText("Magies", TextSizeM, ColorGreen, offsetX, offsetH)
	offsetH += h + 10

	powernames := p.GetSortedPowernames()

	categoriesOffset := map[game.MagickCategory]int32{
		game.MagickCategoryPhysical: 150,
		game.MagickCategoryAstral:   300,
		game.MagickCategoryMental:   450,
		game.MagickCategoryHigh:     600,
	}

	elementsOffset := map[game.MagickElement]int32{
		game.MagickElementEarth: 40,
		game.MagickElementWater: 80,
		game.MagickElementAir:   120,
		game.MagickElementFire:  160,
		game.MagickElementEther: 200,
	}

	ui.DrawText(fmt.Sprintf("Physique (%d)", p.GetMagickLevel(game.MagickCategoryPhysical)), TextSizeM, ColorGreen, offsetX+categoriesOffset[game.MagickCategoryPhysical], offsetH)
	ui.DrawText(fmt.Sprintf("Astrale (%d)", p.GetMagickLevel(game.MagickCategoryAstral)), TextSizeM, ColorGreen, offsetX+categoriesOffset[game.MagickCategoryAstral], offsetH)
	ui.DrawText(fmt.Sprintf("Mentale (%d)", p.GetMagickLevel(game.MagickCategoryMental)), TextSizeM, ColorGreen, offsetX+categoriesOffset[game.MagickCategoryMental], offsetH)
	ui.DrawText(fmt.Sprintf("Sacré (%d)", p.GetMagickLevel(game.MagickCategoryHigh)), TextSizeM, ColorGreen, offsetX+categoriesOffset[game.MagickCategoryHigh], offsetH)

	ui.DrawText(fmt.Sprintf("Terre (%d)", p.GetElementalAffinity(game.MagickElementEarth)), TextSizeM, ColorGreen, offsetX, offsetH+elementsOffset[game.MagickElementEarth])
	ui.DrawText(fmt.Sprintf("Eau (%d)", p.GetElementalAffinity(game.MagickElementWater)), TextSizeM, ColorGreen, offsetX, offsetH+elementsOffset[game.MagickElementWater])
	ui.DrawText(fmt.Sprintf("Air (%d)", p.GetElementalAffinity(game.MagickElementAir)), TextSizeM, ColorGreen, offsetX, offsetH+elementsOffset[game.MagickElementAir])
	ui.DrawText(fmt.Sprintf("Feu (%d)", p.GetElementalAffinity(game.MagickElementFire)), TextSizeM, ColorGreen, offsetX, offsetH+elementsOffset[game.MagickElementFire])
	ui.DrawText(fmt.Sprintf("Ether (%d)", p.GetElementalAffinity(game.MagickElementEther)), TextSizeM, ColorGreen, offsetX, offsetH+elementsOffset[game.MagickElementEther])

	magicksNumber := make(map[string]int32)
	for _, powername := range powernames {
		power := p.Powers[powername]
		n, e := magicksNumber[string(power.Category)+string(power.Element)]
		if !e {
			magicksNumber[string(power.Category)+string(power.Element)] = 0
			n = 0
		}
		x := int32((PlayerMenuOffsetX+n)*Res) + categoriesOffset[power.Category]
		y := offsetH + elementsOffset[power.Element]
		if p.CurrentPower != nil && p.CurrentPower.UID == power.UID {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["ʆ"][0],
				&sdl.Rect{X: x, Y: y, W: Res, H: Res})
		}
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[power.Tile][0],
			&sdl.Rect{X: x, Y: y, W: Res, H: Res})
		magicksNumber[string(power.Category)+string(power.Element)]++
	}
	offsetH += 250

	_, h = ui.DrawText("Description", TextSizeM, ColorWhite, offsetX, offsetH)
	offsetH += h + 10
	ui.DrawPower(p.CurrentPower, PlayerMenuOffsetX, offsetH)
	offsetH += 32 + 40

	return offsetH
}

func (ui *UI) DrawPower(power *game.Power, offsetX, offsetH int32) int32 {
	_, h := ui.DrawText(
		power.Name,
		TextSizeM,
		ColorGreen,
		offsetX*Res,
		offsetH)

	offsetH += h
	_, h = ui.DrawText(
		power.Description,
		TextSizeM,
		ColorActive,
		offsetX*Res,
		offsetH)

	offsetH += h
	_, h = ui.DrawText(
		"Coût en énergie : "+strconv.Itoa(power.Energy),
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
				&ui.textureIndex["ß"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
