package ui2d

func (ui *UI) drawPlayer() {
	p := ui.Game.Level.Player
	texture := ui.playerTextures[p.Name]
	if p.Weapon != nil {
		texture = ui.playerTextures[p.Name+"_with_"+p.Weapon.Typ]
	}
	ui.drawCharacter(&p.Character, texture)
}
