package ui2d

import (
	"thelema/game"
)

func (ui *UI) drawNpc(p *game.Npc) {
	ui.drawCharacter(&p.Character, ui.npcTextures[p.Name])
}
