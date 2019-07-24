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
		ui.drawFightingFriends(offsetX+32, offsetY+100)
		ui.drawFightingEnemies(offsetX*2, offsetY)
		if fr.CurrentEffect != nil {
			ui.drawFightingEffect(fr.CurrentEffect, offsetX, offsetY)
		}
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
			enemy, isPnj := e.(*game.Enemy)
			if fr.AttackTargetSelectionOpen {
				att := fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
				if i >= fr.TargetSelected && i < fr.TargetSelected+att.Range {
					ui.renderer.Copy(ui.textureAtlas,
						&ui.textureIndex['ʆ'][0],
						&sdl.Rect{X: int32(offsetX), Y: int32(offsetY), W: 64, H: 64})
				}
			}
			if isPnj {
				ui.drawFightingEnemy(enemy, offsetX, offsetY)
			} else {
				ui.drawFightingMonster(e, offsetX, offsetY)
			}
			offsetX += 20
			offsetY += 70
		}
	}
}

func (ui *UI) drawFightingEnemy(e *game.Enemy, offsetX, offsetY int) {
	texture := ui.pnjTextures[e.Name]
	xb := 0
	tileY := 9 * 64
	tileX := 64 * ((-1*e.Xb + Res) / (Res / 8))
	if e.IsAttacking() {
		xb = (ui.WindowHeight / 3) * e.AttackPos / 32
		tileY = tileY + 4*64
		tileX = 64 * (6 * e.AttackPos / 32)
	}
	if e.IsDead() {
		tileY = 20 * 64
		tileX = 64 * 5
	}
	if e.IsHurt() > 0 {
		xb = -16
	}
	ui.renderer.Copy(texture,
		&sdl.Rect{X: int32(tileX), Y: int32(tileY), W: 64, H: 64},
		&sdl.Rect{X: int32(offsetX - xb), Y: int32(offsetY), W: 64, H: 64})
	ui.drawHealthBar(int32(offsetX-xb), int32(offsetY-15), e.GetHealth())
}

func (ui *UI) drawFightingMonster(e game.FighterInterface, offsetX, offsetY int) {
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
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[e.GetTile()][0],
			&sdl.Rect{X: int32(offsetX - xb), Y: int32(offsetY + yb), W: 64, H: 64})
		ui.drawHealthBar(int32(offsetX-xb), int32(offsetY-15), e.GetHealth())
	}
}

func (ui *UI) drawFightingFriends(offsetX, offsetY int) {
	fr := ui.Game.FightingRing
	if fr != nil && len(fr.Friends) > 0 {
		for _, e := range fr.Friends {
			f, isFriend := e.(*game.Friend)
			if isFriend {
				ui.drawFightingFriend(f, offsetX, offsetY)
			} else {
				ui.drawFightingInvocation(e, offsetX, offsetY)
			}
			offsetX += 20
			offsetY += 70
		}
	}
}

func (ui *UI) drawFightingFriend(f *game.Friend, offsetX, offsetY int) {
	texture := ui.pnjTextures[f.Name]
	xb := 0
	tileY := 11 * 64
	tileX := 64 * ((-1*f.Xb + Res) / (Res / 8))
	if f.IsAttacking() {
		xb = (ui.WindowHeight / 3) * f.AttackPos / 32
		tileY = tileY + 4*64
		tileX = 64 * (6 * f.AttackPos / 32)
	}
	if f.IsDead() {
		tileY = 20 * 64
		tileX = 64 * 5
	}
	if f.IsHurt() > 0 {
		xb = -16
	}
	ui.renderer.Copy(texture,
		&sdl.Rect{X: int32(tileX), Y: int32(tileY), W: 64, H: 64},
		&sdl.Rect{X: int32(offsetX + xb), Y: int32(offsetY), W: 64, H: 64})
	ui.drawHealthBar(int32(offsetX+xb), int32(offsetY-15), f.GetHealth())
}

func (ui *UI) drawFightingInvocation(e game.FighterInterface, offsetX, offsetY int) {
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
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[e.GetTile()][0],
			&sdl.Rect{X: int32(offsetX + xb), Y: int32(offsetY + yb), W: 64, H: 64})
		ui.drawHealthBar(int32(offsetX+xb), int32(offsetY-15), e.GetHealth())
	}
}

func (ui *UI) drawFightingScreen() {
	ui.renderer.Copy(ui.backgroundTextures["outdoor"],
		&sdl.Rect{X: 0, Y: 0, W: 1280, H: 832},
		&sdl.Rect{X: 0, Y: 0, W: int32(ui.WindowWidth), H: int32(ui.WindowHeight)})
}

func (ui *UI) drawFightingMenu() {
	menu := ui.Game.FightingRing.Menu
	if menu != nil && menu.IsOpen {
		ui.renderer.Copy(ui.uiTextures["downbox"],
			&sdl.Rect{X: 0, Y: 0, W: 320, H: 64},
			&sdl.Rect{X: 0, Y: int32(3 * ui.WindowHeight / 4), W: int32(ui.WindowWidth), H: int32(ui.WindowHeight / 4)})

		var offsetH int32 = int32((3 * ui.WindowHeight / 4) + 15)
		var offsetW int32 = 40
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
			ui.renderer.Copy(tex, nil, &sdl.Rect{offsetW, offsetH, w, h})
			offsetW += w + 20
		}
	}
}

func (ui *UI) drawFightingAttacks() {
	fr := ui.Game.FightingRing
	if fr != nil && fr.AttacksMenuOpen {
		offsetY := (3 * ui.WindowHeight / 4) + 50
		selectedAttack := fr.GetSelectedAttack()
		tex := ui.GetTexture(selectedAttack.Name, TextSizeXL, ColorWhite)
		_, _, w, h, _ := tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{40, int32(offsetY), w, h})

		offsetX := ui.WindowWidth / 2
		tex = ui.GetTexture(fmt.Sprintf("Power: %d", selectedAttack.Damages), TextSizeXL, ColorWhite)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{int32(offsetX), int32(offsetY), w, h})
		tex = ui.GetTexture(fmt.Sprintf("Energy cost: %d", selectedAttack.EnergyCost), TextSizeXL, ColorWhite)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{int32(offsetX), int32(offsetY + 32), w, h})

		offsetX = 3 * ui.WindowWidth / 4
		tex = ui.GetTexture(fmt.Sprintf("Speed: %d", selectedAttack.Speed), TextSizeXL, ColorWhite)
		_, _, w, h, _ = tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{int32(offsetX), int32(offsetY), w, h})
		// tex = ui.GetTexture(fmt.Sprintf("Accuracy: %d", selectedAttack.Accuracy), TextSizeXL, ColorWhite)
		// _, _, w, h, _ = tex.Query()
		// ui.renderer.Copy(tex, nil, &sdl.Rect{int32(offsetX), int32(offsetY + 32), w, h})

	}
}

func (ui *UI) drawFightingEffect(effect *game.Effect, offsetX, offsetY int) {
	fp := effect.Pos
	x := offsetX*(fp.X+1) + rand.Intn(6)
	y := offsetY + 50*fp.Y + rand.Intn(6)
	tile := game.Tile(effect.Rune)
	if len(ui.textureIndex[tile]) > 0 {
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][effect.TileIdx%len(ui.textureIndex[tile])],
			&sdl.Rect{X: int32(x), Y: int32(y), W: 64, H: 64})
	}
}
