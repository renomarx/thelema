package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawEffect(pos game.Pos, effect *game.Effect) {
	tile := game.Tile(effect.Rune)
	idx := 0
	if effect.Size%2 == 1 {
		idx = 1
	}
	if effect.Size >= 10 {
		idx += 2
	}
	if effect.Size >= 50 {
		idx += 2
	}
	if len(ui.textureIndex[tile]) > 0 {
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][idx],
			&sdl.Rect{X: int32(pos.X*Res) + ui.Cam.X, Y: int32(pos.Y*Res) + ui.Cam.Y, W: Res, H: Res})
	}
}
