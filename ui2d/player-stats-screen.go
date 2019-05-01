package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

func (ui *UI) DrawPlayerStats() {
	p := ui.Game.Level.Player
	if !p.IsTalking && !ui.Game.Paused {
		ui.drawPlayerStatsBox()
		offsetH := ui.WindowHeight - Res
		w, _ := ui.DrawText(
			"Health : "+strconv.Itoa(p.Health.Current)+"/"+strconv.Itoa(p.Health.Initial),
			TextSizeM,
			ColorActive,
			10,
			int32(offsetH))
		ui.DrawText(
			"Energy : "+strconv.Itoa(p.Energy.Current)+"/"+strconv.Itoa(p.Energy.Initial),
			TextSizeM,
			ColorActive,
			w+Res+10,
			int32(offsetH))
	}
}

func (ui *UI) drawPlayerStatsBox() {
	for x := 0; x <= ui.WindowWidth/Res; x++ {
		for y := ui.WindowHeight/Res - 1; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['Ʈ'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
