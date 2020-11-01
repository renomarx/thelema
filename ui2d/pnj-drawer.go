package ui2d

import (
	"github.com/renomarx/thelema/game"
)

func (ui *UI) drawNpc(p *game.Npc) {
	ui.drawCharacter(&p.Character, ui.npcTextures[p.Name])
}
