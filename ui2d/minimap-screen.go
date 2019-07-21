package ui2d

import (
	"math"
	"thelema/game"
)

func (ui *UI) DrawMinimap() {
	level := ui.Game.Level
	player := level.Player
	if !player.Menu.IsOpen {
		mapHeight := ui.WindowHeight / 10
		mapWidth := ui.WindowWidth / 10

		CamX := int32((mapWidth / 2) - player.X)
		CamY := int32((mapHeight / 2) - player.Y)

		minY := int(math.Floor(math.Max(0, float64(player.Y-(mapHeight/2)-2))))
		maxY := int(math.Floor(math.Min(float64(len(level.Map)), float64(player.Y+(mapHeight/2)+2))))
		minX := int(math.Floor(math.Max(0, float64(player.X-(mapWidth/2)-2))))
		maxX := int(math.Floor(math.Min(float64(len(level.Map[0])), float64(player.X+(mapWidth/2)+2))))
		for y := minY; y < maxY; y++ {
			row := level.Map[y]
			for x := minX; x < maxX; x++ {
				tile := row[x].T
				r := 0
				g := 0
				b := 0
				switch tile {
				case game.DirtFloor, game.HerbFloor:
					r = 255
					g = 219
					b = 182
				}

				ui.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
				ui.renderer.DrawPoint(int32(x)+CamX, int32(y)+CamY)
				object := row[x].Object
				if object != nil {
					ui.drawMapObject(game.Pos{X: x + int(CamX), Y: y + int(CamY)}, game.Tile(object.Rune))
				}
			}
		}

		// Player
		ui.drawMapPlayer(game.Pos{X: player.X + int(CamX), Y: player.Y + int(CamY)}, 2)
		ui.renderer.SetDrawColor(0, 0, 0, 0)
	}
}
