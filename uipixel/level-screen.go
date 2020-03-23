package uipixel

import (
	"math"
)

func (ui *UI) DrawLevel() {
	level := ui.Game.Level
	if level != nil && level.Player != nil {
		player := level.Player

		ui.Cam.X = (ui.WindowWidth / 2) - float64(player.X)*Res + float64(player.Xb)
		ui.Cam.Y = (ui.WindowHeight / 2) - float64(player.Y)*Res + float64(player.Yb)

		minY := int(math.Floor(math.Max(0, float64(player.Y)-(ui.WindowHeight/2/Res)-2)))
		maxY := int(math.Floor(math.Min(float64(len(level.Map)), float64(player.Y)+(ui.WindowHeight/2/Res)+2)))
		minX := int(math.Floor(math.Max(0, float64(player.X)-(ui.WindowWidth/2/Res)-2)))
		maxX := int(math.Floor(math.Min(float64(len(level.Map[0])), float64(player.X)+(ui.WindowWidth/2/Res)+2)))
		levelMap, mapExists := ui.mapTextures[level.Name]
		if mapExists {
			mapSprite := ui.NewSprite(levelMap, levelMap.Bounds())
			x := float64(player.X)*Res - float64(player.Xb)
			y := float64(player.Y)*Res - float64(player.Yb)
			ui.DrawSprite(mapSprite, float64(ui.WindowWidth/2-x), float64(ui.WindowHeight/2-y))
		} else {
			for y := minY; y < maxY; y++ {
				row := level.Map[y]
				for x := minX; x < maxX; x++ {
					c := row[x]
					tile := c.T
					if len(ui.textureIndex[tile]) > 0 {
						spr := ui.NewSprite(ui.textureAtlas, ui.textureIndex[tile][(x*(y+1)+y*(x+3))%len(ui.textureIndex[tile])])
						ui.DrawSprite(spr, float64(x)*Res+ui.Cam.X, float64(y)*Res+ui.Cam.Y)
					}
				}
			}
		}
		//
		// for y := minY; y < maxY; y++ {
		// 	row := level.Map[y]
		// 	for x := minX; x < maxX; x++ {
		// 		c := row[x]
		// 		object := c.Object
		// 		if object != nil {
		// 			if !mapExists || !object.Static {
		// 				ui.drawObject(game.Pos{X: x, Y: y}, game.Tile(object.Rune))
		// 			}
		// 		}
		// 		pnj := c.Pnj
		// 		if pnj != nil {
		// 			ui.drawPnj(pnj)
		// 		}
		// 		effect := c.Effect
		// 		if effect != nil {
		// 			ui.drawEffect(game.Pos{X: x, Y: y}, effect)
		// 		}
		// 	}
		// }
		//
		// ui.drawPlayer()
		// if player.TalkingTo != nil {
		// 	ui.DrawDialog(player.TalkingTo)
		// }
		//
		// ui.DrawMinimap()
		// ui.DrawPlayerStats()
		// ui.DrawPlayerMenu()
	}
}
