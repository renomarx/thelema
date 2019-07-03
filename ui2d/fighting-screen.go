package ui2d

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	// "github.com/veandco/go-sdl2/ttf"
	// "log"
	// "path/filepath"
	// "thelema/game"
)

func (ui *UI) DrawFightingRing() {
	fr := ui.Game.FightingRing
	if fr != nil && fr.IsOpen {
		ui.drawFightingScreen()
		ui.drawFightingMenu()
		ui.drawFightingAttacks()
		ui.drawFightingPlayer()
		ui.drawFightingEnemies()
	}
}

func (ui *UI) drawFightingPlayer() {
	p := ui.Game.Level.Player
	texture := ui.playerTextures[p.Name]
	if p.Weapon != nil {
		texture = ui.playerTextures[p.Name+"_with_"+p.Weapon.Typ]
	}
	tileY := 11 * 64
	tileX := 64 * ((-1*p.Xb + Res) / (Res / 8))
	ui.renderer.Copy(texture,
		&sdl.Rect{X: int32(tileX), Y: int32(tileY), W: 64, H: 64},
		&sdl.Rect{X: 100, Y: 100, W: 64, H: 64})
	ui.drawHealthBar(100, 65, p.GetHealth())
	ui.drawEnergyBar(100, 85, p.GetEnergy())
}

func (ui *UI) drawFightingEnemies() {
	fr := ui.Game.FightingRing
	if fr != nil && len(fr.Enemies) > 0 {
		offsetX := int32(600)
		offsetY := int32(100)
		for _, e := range fr.Enemies {
			if !e.IsDead() {
				ui.renderer.Copy(ui.textureAtlas,
					&ui.textureIndex[e.GetTile()][0],
					&sdl.Rect{X: offsetX, Y: offsetY, W: 32, H: 32})
				ui.drawHealthBar(offsetX, offsetY-15, e.GetHealth())
				offsetX += int32(16)
				offsetY += int32(50)
			}
		}
	}
}

func (ui *UI) drawFightingScreen() {
	for x := 0; x < ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['ß'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}

func (ui *UI) drawFightingMenu() {
	menu := ui.Game.FightingMenu
	if menu != nil && menu.IsOpen {
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

func (ui *UI) drawFightingAttacks() {
	fr := ui.Game.FightingRing
	if fr != nil && fr.AttacksMenuOpen {
		selectedAttack := fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
		tex := ui.GetTexture(fmt.Sprintf(selectedAttack.Name+" (%d)", selectedAttack.Damages), TextSizeXL, ColorActive)
		_, _, w, h, _ := tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{10, 500, w, h})
	}
}
