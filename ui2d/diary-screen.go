package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawDiary() {
	p := ui.Game.Level.Player
	if p.QuestMenuOpen {
		ui.drawDiaryBox()
		var offsetH int32 = 0
		_, h := ui.DrawText("Journal", TextSizeL, ColorActive, PlayerMenuOffsetX*32, offsetH)
		offsetH += h + 10

		todo := ui.Game.GetSteps(game.StepStateTODO)
		done := ui.Game.GetSteps(game.StepStateDONE)

		_, h = ui.DrawText("A faire:", TextSizeM, ColorActive, PlayerMenuOffsetX*32, offsetH)
		offsetH += h
		for _, st := range todo {
			h = ui.DrawStep(st, offsetH, ColorActive)
			offsetH = h
		}
		offsetH += Res
		_, h = ui.DrawText("Terminées:", TextSizeM, ColorDisabled, PlayerMenuOffsetX*32, offsetH)
		offsetH += h
		for _, st := range done {
			h = ui.DrawStep(st, offsetH, ColorDisabled)
			offsetH = h
		}
	}
}

func (ui *UI) DrawStep(st *game.Step, offsetH int32, color sdl.Color) int32 {
	_, h := ui.DrawText("- "+st.Name, TextSizeS, color, PlayerMenuOffsetX*32+10+10, offsetH)
	offsetH += h
	offsetH += 10
	return offsetH
}

func (ui *UI) drawDiaryBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["ß"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
