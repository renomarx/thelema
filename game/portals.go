package game

type Portal struct {
	LevelTo string
	PosTo   Pos
	Key     string
}

func (port *Portal) Discovered(g *Game) bool {
	level, e := g.Levels[port.LevelTo]
	return e && level.Discovered
}
