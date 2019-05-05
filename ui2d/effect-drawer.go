package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawEffect(pos game.Pos, effect *game.Effect) {
	tile := game.Tile(effect.Rune)
	if len(ui.textureIndex[tile]) > 0 {
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][effect.TileIdx%len(ui.textureIndex[tile])],
			&sdl.Rect{X: int32(pos.X*Res) + ui.Cam.X, Y: int32(pos.Y*Res) + ui.Cam.Y, W: Res, H: Res})
	}
}
