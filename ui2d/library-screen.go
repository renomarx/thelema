package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) DrawLibrary() {
	lib := ui.Game.Level.Player.Library
	if lib.IsOpen {
		ui.drawLibraryBox()
		var offsetX = int32(PlayerMenuOffsetX*Res + 10)
		var offsetH int32 = 0
		_, h := ui.DrawText("Books", TextSizeL, ColorActive, offsetX, offsetH)
		offsetH += h + 10
		for i, book := range lib.Books {
			x := int32((PlayerMenuOffsetX + i) * Res)
			if lib.IsHighlighted(i) {
				ui.renderer.Copy(ui.textureAtlas,
					&ui.textureIndex["ʆ"][0],
					&sdl.Rect{X: x, Y: offsetH, W: Res, H: Res})
			}
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex[game.Tile(book.Rune)][0],
				&sdl.Rect{X: int32((PlayerMenuOffsetX + i) * Res), Y: offsetH, W: Res, H: Res})
		}
		offsetH += 32 + 40

		currentBook := lib.GetCurrentBook()
		if currentBook != nil {
			_, h = ui.DrawText(currentBook.Title, TextSizeL, ColorActive, PlayerMenuOffsetX*Res, offsetH)
			offsetH += h + 10
			ui.DrawBook(currentBook, offsetH)
		}
	}
}

func (ui *UI) DrawBook(book *game.OBook, offsetH int32) {
	var offsetX = int32(PlayerMenuOffsetX*Res + 30)
	for _, line := range book.Text {
		_, h := ui.DrawText(line, TextSizeM, ColorActive, offsetX, offsetH)
		offsetH += h
	}
}

func (ui *UI) drawLibraryBox() {
	for x := PlayerMenuOffsetX; x <= ui.WindowWidth/Res; x++ {
		for y := 0; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex["Ʈ"][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
