package ui2d

import (
	"thelema/game"
)

func (ui *UI) drawFriend(p *game.Friend) {
	ui.drawCharacter(&p.Character, ui.pnjTextures[p.Name])
}
