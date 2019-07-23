package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"sort"
	"thelema/game"
)

func (ui *UI) DrawQuests() {
	p := ui.Game.Level.Player
	if p.QuestMenuOpen {
		ui.drawQuestsBox()
		var offsetH int32 = 0
		_, h := ui.DrawText("Quêtes", TextSizeL, ColorActive, PlayerMenuOffsetX*32, offsetH)
		offsetH += h + 10

		runningIds := make([]string, 0, len(p.Quests))
		finishedIds := make([]string, 0, len(p.Quests))
		for id, q := range p.Quests {
			if q.IsRunning() {
				runningIds = append(runningIds, id)
			} else if q.IsFinished {
				finishedIds = append(finishedIds, id)
			}
		}
		sort.Strings(runningIds)
		sort.Strings(finishedIds)

		_, h = ui.DrawText("En cours:", TextSizeM, ColorActive, PlayerMenuOffsetX*32, offsetH)
		offsetH += h
		i := 0
		for _, qid := range runningIds {
			q := p.Quests[qid]
			h = ui.DrawQuest(q, offsetH, ColorActive)
			offsetH = h
			i++
		}
		offsetH += Res
		_, h = ui.DrawText("Terminées:", TextSizeM, ColorDisabled, PlayerMenuOffsetX*32, offsetH)
		offsetH += h
		i = 0
		for _, qid := range finishedIds {
			q := p.Quests[qid]
			h = ui.DrawFinishedQuest(q, offsetH, ColorDisabled)
			offsetH = h
			i++
		}
	}
}

func (ui *UI) DrawQuest(q *game.Quest, offsetH int32, color sdl.Color) int32 {
	_, h := ui.DrawText(q.Name, TextSizeM, color, PlayerMenuOffsetX*32+10, offsetH)
	offsetH += h

	steps := q.GetOrderedSteps()
	for _, st := range steps {
		stColor := color
		if st.IsFinished {
			stColor = ColorDisabled
		}
		_, h := ui.DrawText("- "+st.Description, TextSizeS, stColor, PlayerMenuOffsetX*32+10+10, offsetH)
		offsetH += h
	}
	offsetH += 10
	return offsetH
}

func (ui *UI) DrawFinishedQuest(q *game.Quest, offsetH int32, color sdl.Color) int32 {
	_, h := ui.DrawText(q.Name, TextSizeM, color, PlayerMenuOffsetX*32+10, offsetH)
	offsetH += h
	return offsetH
}

func (ui *UI) drawQuestsBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['Ʈ'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
