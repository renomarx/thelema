package uipixel

import "golang.org/x/image/colornames"

func (ui *UI) DrawMenu() {
	menu := ui.Game.GetMenu()
	if menu.IsOpen {
		ui.drawMenuBox()
		var offsetH float64 = 20
		for _, choice := range menu.Choices {
			color := colornames.White
			if choice.Highlighted {
				color = colornames.Green
			}
			if choice.Disabled {
				color = colornames.Gray
			}
			_, h := ui.DrawText(choice.Cmd, TextSizeXL, color, 20, offsetH)
			offsetH += h + 10
		}
	}
}

func (ui *UI) drawMenuBox() {
	ui.win.Clear(colornames.Black)
}
