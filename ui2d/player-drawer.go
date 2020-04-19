package ui2d

func (ui *UI) drawPlayer() {
	l := ui.Game.Level
	p := l.Player
	texture := ui.playerTextures[p.Name]
	ui.drawCharacter(&p.Character, texture)
	c := l.Map[p.Z][p.Y][p.X]
	effect := c.Effect
	if effect != nil {
		ui.drawEffect(p.Pos, effect)
	}
}
