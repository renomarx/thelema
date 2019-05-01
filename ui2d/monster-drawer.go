package ui2d

import (
	"math/rand"
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawMonster(pos game.Pos, m *game.Monster) {
	tile := game.Tile(m.Rune)
	if len(ui.textureIndex[tile]) > 0 {

		xb := m.Xb
		yb := m.Yb
		fieldLen := 4
		if m.IsAttacking {
			xb = rand.Intn(fieldLen*2) - fieldLen
			yb = rand.Intn(fieldLen*2) - fieldLen
		}

		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][m.Beauty.Initial%len(ui.textureIndex[tile])],
			&sdl.Rect{X: int32(pos.X*Res-xb) + ui.Cam.X, Y: int32(pos.Y*Res-yb) + ui.Cam.Y, W: Res, H: Res})
	}
}
