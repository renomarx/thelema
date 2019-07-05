package ui2d

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	// "github.com/veandco/go-sdl2/ttf"
	// "log"
	// "path/filepath"
	"thelema/game"
)

func (ui *UI) DrawFightingRing() {
	fr := ui.Game.FightingRing
	offsetX := ui.WindowWidth / 3
	offsetY := ui.WindowHeight / 3
	if fr != nil && fr.IsOpen {
		ui.drawFightingScreen()
		ui.drawFightingMenu()
		ui.drawFightingAttacks()
		ui.drawFightingPlayer(offsetX, offsetY)
		ui.drawFightingEnemies(offsetX*2, offsetY)
	}
}

func (ui *UI) drawFightingPlayer(offsetX, offsetY int) {
	fr := ui.Game.FightingRing
	p := ui.Game.Level.Player
	texture := ui.playerTextures[p.Name]
	xb := 0
	tileY := 11 * 64
	tileX := 64 * ((-1*p.Xb + Res) / (Res / 8))
	if p.IsAttacking() {
		att := fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
		switch att.Type {
		case game.AttackTypePhysical:
			xb = (ui.WindowHeight / 3) * p.AttackPos / 32
			tileY = tileY + 4*64
			tileX = 64 * (6 * p.AttackPos / 32)
		case game.AttackTypeMagick:
			tileY = tileY - 8*64
			tileX = 64 * (p.AttackPos / 6)
		}
	}
	if p.IsDead() {
		tileY = 20 * 64
		tileX = 64 * 5
	}
	if p.IsHurt() > 0 {
		xb = -16
	}
	ui.renderer.Copy(texture,
		&sdl.Rect{X: int32(tileX), Y: int32(tileY), W: 64, H: 64},
		&sdl.Rect{X: int32(offsetX + xb), Y: int32(offsetY), W: 64, H: 64})
	ui.drawHealthBar(int32(offsetX+xb), int32(offsetY-35), p.GetHealth())
	ui.drawEnergyBar(int32(offsetX+xb), int32(offsetY-15), p.GetEnergy())
}

func (ui *UI) drawFightingEnemies(offsetX, offsetY int) {
	fr := ui.Game.FightingRing
	if fr != nil && len(fr.Enemies) > 0 {
		for i, e := range fr.Enemies {
			if !e.IsDead() {
				xb := 0
				yb := 0
				fieldLen := 4
				if e.IsAttacking() {
					xb = ui.WindowHeight / 3
					yb = rand.Intn(fieldLen*2) - fieldLen
				}
				if e.IsHurt() > 0 {
					xb = -16
				}
				if fr.AttackTargetSelectionOpen {
					att := fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
					if i >= fr.TargetSelected && i < fr.TargetSelected+att.Range {
						ui.renderer.Copy(ui.textureAtlas,
							&ui.textureIndex['ʆ'][0],
							&sdl.Rect{X: int32(offsetX), Y: int32(offsetY), W: 32, H: 32})
					}
				}
				ui.renderer.Copy(ui.textureAtlas,
					&ui.textureIndex[e.GetTile()][0],
					&sdl.Rect{X: int32(offsetX - xb), Y: int32(offsetY + yb), W: 32, H: 32})
				ui.drawHealthBar(int32(offsetX-xb), int32(offsetY-15), e.GetHealth())
				offsetX += 16
				offsetY += 50
			}
		}
	}
}

func (ui *UI) drawFightingScreen() {
	ui.renderer.Copy(ui.backgroundTextures["outdoor"],
		&sdl.Rect{X: 0, Y: 0, W: 1280, H: 832},
		&sdl.Rect{X: 0, Y: 0, W: int32(ui.WindowWidth), H: int32(ui.WindowHeight)})
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
