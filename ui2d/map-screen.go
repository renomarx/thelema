package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) drawMapBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["Ʈ"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}

func (ui *UI) DrawMap() {
	g := ui.Game
	level := g.Level
	player := level.Player
	if player.MapMenuOpen {
		ui.drawMapBox()

		mapHeight := ui.WindowHeight / 2
		mapWidth := ui.WindowWidth / 2
		offsetX := ui.WindowWidth / 4
		offsetY := ui.WindowHeight / 4

		for y := 0; y < mapHeight; y++ {
			for x := 0; x < mapWidth; x++ {
				r := 255
				g := 219
				b := 182

				ui.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
				ui.renderer.DrawPoint(int32(x+offsetX), int32(y+offsetY))
				ui.renderer.SetDrawColor(0, 0, 0, 0)
			}
		}

		portals := make(map[game.Pos]*game.Portal)
		for y := 0; y < len(level.Map[player.Z]); y++ {
			row := level.Map[player.Z][y]
			for x := 0; x < len(row); x++ {
				mapX := x * mapWidth / len(row)
				mapY := y * mapHeight / len(level.Map)

				object := row[x].Object
				if object != nil {
					ui.drawMapObject(game.Pos{X: mapX + offsetX, Y: mapY + offsetY}, game.Tile(object.Rune))
				}

				portal := row[x].Portal
				if portal != nil {
					portals[game.Pos{X: mapX, Y: mapY}] = portal
				}
			}
		}

		// Cities names
		for pos, portal := range portals {
			if portal.Discovered(ui.Game) {
				tex := ui.GetTexture(portal.LevelTo, TextSizeXS, ColorWhite)
				_, _, w, h, _ := tex.Query()
				for j := -2; j < int(h)+2; j++ {
					for i := -2; i < int(w)+2; i++ {
						ui.renderer.SetDrawColor(0, 0, 0, 255)
						ui.renderer.DrawPoint(int32(pos.X+offsetX+i), int32(pos.Y+offsetY+j))
					}
				}
				ui.renderer.Copy(tex, nil, &sdl.Rect{X: int32(pos.X + offsetX), Y: int32(pos.Y + offsetY), W: w, H: h})
			}
		}

		// Player
		mapX := player.X * mapWidth / len(level.Map[0])
		mapY := player.Y * mapHeight / len(level.Map)
		ui.drawMapPlayer(game.Pos{X: mapX + offsetX, Y: mapY + offsetY}, 3)
	}
}

func (ui *UI) drawMapObject(pos game.Pos, tile game.Tile) {
	r := 0
	g := 0
	b := 0
	rr := 0
	switch tile {
	case game.Ocean:
		r = 13
		g = 61
		b = 122
	case game.Door:
		r = 1
		g = 126
		b = 255
		rr = 1
	case game.Wall:
		r = 136
		g = 134
		b = 131
	}
	ui.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
	ui.drawMapPoints(pos, rr)
	ui.renderer.SetDrawColor(0, 0, 0, 0)
}

func (ui *UI) drawMapPlayer(pos game.Pos, ray int) {
	ui.renderer.SetDrawColor(255, 0, 0, 255)
	ui.drawMapPoints(pos, ray)
	ui.renderer.SetDrawColor(0, 0, 0, 0)
}

func (ui *UI) drawMapPoints(pos game.Pos, ray int) {
	if ray == 0 {
		ui.renderer.DrawPoint(int32(pos.X), int32(pos.Y))
		return
	}
	for y := pos.Y - ray; y < pos.Y+ray; y++ {
		for x := pos.X - ray; x < pos.X+ray; x++ {
			if x > 0 && y > 0 {
				ui.renderer.DrawPoint(int32(x), int32(y))
			}
		}
	}
}
