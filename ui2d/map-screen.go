package ui2d

import (
	"math"
	"thelema/game"
)

func (ui *UI) DrawMap() {
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
				tile := row[x]
				r := 0
				g := 0
				b := 0
				switch tile {
				case game.DirtFloor:
					r = 255
					g = 219
					b = 182
				}

				ui.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
				ui.renderer.DrawPoint(int32(x)+CamX, int32(y)+CamY)
			}
		}

		for y := minY; y < maxY; y++ {
			for x := minX; x < maxX; x++ {
				pos := game.Pos{X: x, Y: y}
				game.Mux.Lock()
				object, exists := level.Objects[pos]
				game.Mux.Unlock()
				if exists {
					ui.drawMapObject(game.Pos{X: pos.X + int(CamX), Y: pos.Y + int(CamY)}, game.Tile(object.Rune))
				}
			}
		}

		// Player
		ui.drawMapPlayer(game.Pos{X: player.X + int(CamX), Y: player.Y + int(CamY)})
		ui.renderer.SetDrawColor(0, 0, 0, 0)
	}
}

func (ui *UI) drawMapObject(pos game.Pos, tile game.Tile) {
	r := 0
	g := 0
	b := 0
	switch tile {
	case game.Ocean:
		r = 0
		g = 0
		b = 255
	case game.Tree:
		r = 0
		g = 255
		b = 0
	}
	ui.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
	ui.renderer.DrawPoint(int32(pos.X), int32(pos.Y))
}

func (ui *UI) drawMapPlayer(pos game.Pos) {
	ui.renderer.SetDrawColor(255, 0, 0, 255)
	for y := pos.Y - 2; y < pos.Y+2; y++ {
		for x := pos.X - 2; x < pos.X+2; x++ {
			if x > 0 && y > 0 {
				ui.renderer.DrawPoint(int32(x), int32(y))
			}
		}
	}
}
