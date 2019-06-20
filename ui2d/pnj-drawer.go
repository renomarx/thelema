package ui2d

import (
	"thelema/game"
)

func (ui *UI) drawPnj(p *game.Pnj) {
	ui.drawCharacter(&p.Character, ui.pnjTextures[p.Name])
}
