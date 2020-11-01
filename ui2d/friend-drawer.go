package ui2d

import (
	"github.com/renomarx/thelema/game"
)

func (ui *UI) drawFriend(p *game.Friend) {
	ui.drawCharacter(&p.Character, ui.npcTextures[p.Name])
}
