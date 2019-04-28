package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"thelema/game"
)

func (ui *UI) drawMapBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['Ʈ'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}

func (ui *UI) DrawMap() {
	level := ui.Game.Level
	player := level.Player
	if player.MapMenuOpen {
		ui.drawMapBox()

		// Working only because game world width < screen width && game world height < screen height
		CamX := int32((ui.WindowWidth - len(level.Map[0])) / 2)
		CamY := int32((ui.WindowHeight - len(level.Map)) / 2)

		for y := 0; y < len(level.Map); y++ {
			row := level.Map[y]
			for x := 0; x < len(row); x++ {
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

		for pos, object := range level.Objects {
			ui.drawMapObject(game.Pos{X: pos.X + int(CamX), Y: pos.Y + int(CamY)}, game.Tile(object.Rune))
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
	ui.drawMapPoints(pos, 1)
}

func (ui *UI) drawMapPlayer(pos game.Pos) {
	ui.renderer.SetDrawColor(255, 0, 0, 255)
	ui.drawMapPoints(pos, 3)
}

func (ui *UI) drawMapPoints(pos game.Pos, ray int) {
	for y := pos.Y - ray; y < pos.Y+ray; y++ {
		for x := pos.X - ray; x < pos.X+ray; x++ {
			if x > 0 && y > 0 {
				ui.renderer.DrawPoint(int32(x), int32(y))
			}
		}
	}
}
