package ui2d

func (ui *UI) drawPlayer() {
	p := ui.Game.Level.Player
	ui.drawCharacter(&p.Character, ui.playerTextures[p.Name])
}
