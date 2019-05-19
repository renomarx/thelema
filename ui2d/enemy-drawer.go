package ui2d

import (
	"thelema/game"
)

func (ui *UI) drawEnemy(p *game.Enemy) {
	ui.drawCharacter(&p.Character, ui.pnjTextures[p.Name])
}
