package ui2d

func (ui *UI) drawPlayer() {
	p := ui.Game.Level.Player
	texture := ui.playerTextures[p.Name]
	ui.drawCharacter(&p.Character, texture)
}
