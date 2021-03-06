package ui2d

import (
	"github.com/renomarx/thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawCharacter(p *game.Character, texture *sdl.Texture) {
	pos := p.Pos

	tileX := 0
	tileY := 0
	dir := p.LookAt
	if dir == game.Left {
		tileY = 9 * 64
		tileX = 64 * ((p.Xb + Res) / (Res / 8))
	}
	if dir == game.Right {
		tileY = 11 * 64
		tileX = 64 * ((-1*p.Xb + Res) / (Res / 8))
	}
	if dir == game.Up {
		tileY = 8 * 64
		tileX = 64 * ((p.Yb + Res) / (Res / 8))
	}
	if dir == game.Down {
		tileY = 10 * 64
		tileX = 64 * ((-1*p.Yb + Res) / (Res / 8))
	}
	if p.IsPowerUsing {
		tileY = tileY - 8*64
		tileX = 64 * (p.PowerUsingStage / 6)
	}
	if p.IsDead() {
		tileY = 20 * 64
		tileX = 64 * 5
	}
	if p.Meditating {
		tileY = 20 * 64
		tileX = 64 * 2
	}
	if p.Shadow {
		texture.SetColorMod(0, 0, 0)
	} else {
		texture.SetColorMod(255, 255, 255)
	}
	ui.renderer.Copy(texture,
		&sdl.Rect{X: int32(tileX), Y: int32(tileY), W: 64, H: 64},
		&sdl.Rect{X: int32(pos.X*Res-p.Xb-(Res/8)) + ui.Cam.X, Y: int32(pos.Y*Res-p.Yb-(Res/4)) + ui.Cam.Y, W: Res + (Res / 4), H: Res + (Res / 4)})
}
