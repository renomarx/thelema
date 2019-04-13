package ui2d

import (
	"thelema/game"
	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawLevel() {
	level := ui.Game.Level
	if level != nil && level.Player != nil {
		player := level.Player

		ui.Cam.X = int32((ui.WindowWidth / 2) - player.X*Res + player.Xb)
		ui.Cam.Y = int32((ui.WindowHeight / 2) - player.Y*Res + player.Yb)

		ui.renderer.Clear()
		for y, row := range level.Map {
			for x, tile := range row {
				if len(ui.textureIndex[tile]) > 0 {
					srcRect := ui.textureIndex[tile][(x+y)%len(ui.textureIndex[tile])]
					dstRect := sdl.Rect{X: int32(x*Res) + ui.Cam.X, Y: int32(y*Res) + ui.Cam.Y, W: Res, H: Res}

					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
				}
			}
		}

		game.Mux.Lock()
		for pos, object := range level.Objects {
			ui.drawObject(pos, game.Tile(object.Rune))
		}
		game.Mux.Unlock()
		game.Mux.Lock()
		for _, pnj := range level.Pnjs {
			ui.drawPnj(pnj)
		}
		game.Mux.Unlock()
		game.Mux.Lock()
		for pos, monster := range level.Monsters {
			ui.drawMonster(pos, monster)
		}
		game.Mux.Unlock()
		game.Mux.Lock()
		for pos, invoked := range level.Invocations {
			ui.drawInvoked(pos, invoked)
		}
		game.Mux.Unlock()

		ui.drawPlayer()

		game.Mux.Lock()
		for _, projectile := range level.Projectiles {
			ui.drawProjectile(projectile, game.Tile(projectile.Rune))
		}
		game.Mux.Unlock()
		game.Mux.Lock()
		for pos, effect := range level.Effects {
			ui.drawEffect(pos, effect)
		}
		game.Mux.Unlock()

		ui.DrawPlayerMenu()
		ui.DrawPlayerStats()
	}
}
