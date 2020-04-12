package ui2d

import (
	"math"
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
		minY := int(math.Floor(math.Max(0, float64(player.Y-(ui.WindowHeight/2/Res)-2))))
		maxY := int(math.Floor(math.Min(float64(len(level.Map)), float64(player.Y+(ui.WindowHeight/2/Res)+2))))
		minX := int(math.Floor(math.Max(0, float64(player.X-(ui.WindowWidth/2/Res)-2))))
		maxX := int(math.Floor(math.Min(float64(len(level.Map[0])), float64(player.X+(ui.WindowWidth/2/Res)+2))))
		levelMap, mapExists := ui.mapTextures[level.Name]
		levelMapOver, mapOverExists := ui.mapTextures[level.Name+"_over"]

		// Under player
		if mapExists {
			_, _, mapWidth, mapHeight, _ := levelMap.Query()
			x := player.X*Res - player.Xb
			y := player.Y*Res - player.Yb
			ui.renderer.Copy(levelMap,
				&sdl.Rect{X: 0, Y: 0, W: int32(mapWidth), H: int32(mapHeight)},
				&sdl.Rect{X: int32(ui.WindowWidth/2 - x), Y: int32(ui.WindowHeight/2 - y), W: int32(mapWidth), H: int32(mapHeight)})
		} else {
			for y := minY; y < maxY; y++ {
				row := level.Map[y]
				for x := minX; x < maxX; x++ {
					c := row[x]
					tile := c.T
					if len(ui.textureIndex[tile]) > 0 {
						srcRect := ui.textureIndex[tile][(x*(y+1)+y*(x+3))%len(ui.textureIndex[tile])]
						dstRect := sdl.Rect{X: int32(x*Res) + ui.Cam.X, Y: int32(y*Res) + ui.Cam.Y, W: Res, H: Res}

						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}
				}
			}
		}

		// Under player
		for y := minY; y < maxY; y++ {
			row := level.Map[y]
			for x := minX; x < maxX; x++ {
				c := row[x]

				object := c.Object
				if object != nil {
					if !mapExists || !object.Static {
						ui.drawObject(game.Pos{X: x, Y: y}, game.Tile(object.Rune))
					}
				}

				pnj := c.Pnj
				if pnj != nil {
					ui.drawPnj(pnj)
				}
			}
		}

		ui.drawPlayer()

		// Over player
		if mapOverExists {
			_, _, mapWidth, mapHeight, _ := levelMapOver.Query()
			x := player.X*Res - player.Xb
			y := player.Y*Res - player.Yb
			ui.renderer.Copy(levelMapOver,
				&sdl.Rect{X: 0, Y: 0, W: int32(mapWidth), H: int32(mapHeight)},
				&sdl.Rect{X: int32(ui.WindowWidth/2 - x), Y: int32(ui.WindowHeight/2 - y), W: int32(mapWidth), H: int32(mapHeight)})
		}

		// Over player
		for y := minY; y < maxY; y++ {
			row := level.Map[y]
			for x := minX; x < maxX; x++ {
				c := row[x]

				effect := c.Effect
				if effect != nil {
					ui.drawEffect(game.Pos{X: x, Y: y}, effect)
				}
			}
		}

		// Dialogs
		if player.TalkingTo != nil {
			ui.DrawDialog(player.TalkingTo)
		}

		// Menus
		ui.DrawMinimap()
		ui.DrawPlayerStats()
		ui.DrawPlayerMenu()
	}
}
