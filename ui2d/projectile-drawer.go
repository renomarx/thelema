package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawProjectile(p *game.Projectile, tile game.Tile) {
	tile = game.Tile(p.Rune)
	if len(ui.textureIndex[tile]) > 0 {
		pos := p.Pos
		idx := 0
		if p.Direction == game.Left {
			idx = 6
		}
		if p.Direction == game.Right {
			idx = 2
		}
		if p.Direction == game.Down {
			idx = 4
		}
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][idx],
			&sdl.Rect{X: int32(pos.X*Res-p.Xb) + ui.Cam.X, Y: int32(pos.Y*Res-p.Yb) + ui.Cam.Y, W: Res, H: Res})
	}
}
