package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"thelema/game"
)

func (ui *UI) DrawLevel() {
	level := ui.Game.Level
	if level != nil && level.Player != nil {
		player := level.Player

		ui.Cam.X = int32((ui.WindowWidth / 2) - player.X*Res + player.Xb)
		ui.Cam.Y = int32((ui.WindowHeight / 2) - player.Y*Res + player.Yb)

		ui.renderer.Clear()
		minY := int(math.Floor(math.Max(0, float64(player.Y-(ui.WindowHeight/2/Res)-2))))
		maxY := int(math.Floor(math.Min(float64(len(level.Map)), float64(player.Y+(ui.WindowHeight/2/Res)+2))))
		minX := int(math.Floor(math.Max(0, float64(player.X-(ui.WindowWidth/2/Res)-2))))
		maxX := int(math.Floor(math.Min(float64(len(level.Map[0])), float64(player.X+(ui.WindowWidth/2/Res)+2))))
		for y := minY; y < maxY; y++ {
			row := level.Map[y]
			for x := minX; x < maxX; x++ {
				tile := row[x]
				if len(ui.textureIndex[tile]) > 0 {
					srcRect := ui.textureIndex[tile][(x+y)%len(ui.textureIndex[tile])]
					dstRect := sdl.Rect{X: int32(x*Res) + ui.Cam.X, Y: int32(y*Res) + ui.Cam.Y, W: Res, H: Res}

					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
				}
			}
		}

		for y := minY; y < maxY; y++ {
			for x := minX; x < maxX; x++ {
				pos := game.Pos{X: x, Y: y}
				game.Mux.Lock()
				object, exists := level.Objects[pos]
				game.Mux.Unlock()
				if exists {
					ui.drawObject(pos, game.Tile(object.Rune))
				}
			}
		}

		for y := minY; y < maxY; y++ {
			for x := minX; x < maxX; x++ {
				pos := game.Pos{X: x, Y: y}
				game.Mux.Lock()
				pnj, exists := level.Pnjs[pos]
				game.Mux.Unlock()
				if exists {
					ui.drawPnj(pnj)
				}
			}
		}

		for y := minY; y < maxY; y++ {
			for x := minX; x < maxX; x++ {
				pos := game.Pos{X: x, Y: y}
				game.Mux.Lock()
				monster, exists := level.Monsters[pos]
				game.Mux.Unlock()
				if exists {
					ui.drawMonster(pos, monster)
				}
			}
		}

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

		ui.DrawMap()
		ui.DrawPlayerStats()
		ui.DrawPlayerMenu()
	}
}
