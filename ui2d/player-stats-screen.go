package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawPlayerStats() {
	p := ui.Game.Level.Player
	if p.TalkingTo == nil && !ui.Game.Paused {
		ui.drawPlayerStatsBox()
		offsetH := ui.WindowHeight - Res
		ui.drawHealthBar(10, int32(offsetH), p.Health)
		ui.drawEnergyBar(100, int32(offsetH), p.Energy)
		w := 200

		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[p.CurrentPower.Tile][0],
			&sdl.Rect{X: int32(w), Y: int32(offsetH - 10), W: Res, H: Res})
	}
}

func (ui *UI) drawPlayerStatsBox() {
	for x := 0; x <= ui.WindowWidth/Res; x++ {
		for y := ui.WindowHeight/Res - 1; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["Ʈ"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}

func (ui *UI) drawHealthBar(x, y int32, health game.Characteristic) {
	sizeY := 10
	sizeX := 64
	p := health.Current * sizeX / health.Initial

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			r := 77
			if i < p {
				r = 255
			}
			ui.renderer.SetDrawColor(uint8(r), 0, 0, 255)
			ui.renderer.DrawPoint(x+int32(i), y+int32(j))
			ui.renderer.SetDrawColor(0, 0, 0, 0)
		}
	}
}

func (ui *UI) drawEnergyBar(x, y int32, energy game.Characteristic) {
	sizeY := 10
	sizeX := 64
	p := energy.Current * sizeX / energy.Initial

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			b := 77
			if i < p {
				b = 255
			}
			ui.renderer.SetDrawColor(0, 0, uint8(b), 255)
			ui.renderer.DrawPoint(x+int32(i), y+int32(j))
			ui.renderer.SetDrawColor(0, 0, 0, 0)
		}
	}
}
