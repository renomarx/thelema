package game

const LevelTypeOutdoor = "OUTDOOR"
const LevelTypeGrotto = "GROTTO"
const LevelTypeHouse = "HOUSE"

type Level struct {
	Type        string
	Map         [][]Tile
	Player      *Player
	Monsters    map[Pos]*Monster
	Objects     map[Pos]*Object
	Portals     map[Pos]*Portal
	Effects     map[Pos]*Effect
	Projectiles map[Pos]*Projectile
	Pnjs        map[Pos]*Pnj
	Invocations map[Pos]*Invoked
	Paused      bool
}

func (g *Game) UpdateLevel() {
	input := g.input
	if g.Level.Paused {
		g.HandleInputPlayerMenu()
	} else {
		g.handleInput()
		g.handleMonsters()
		g.handlePnjs()
		g.handleInvocations()
		g.handleProjectiles()
		if input.Typ == Select {
			g.OpenPlayerMenu()
		}
	}
	if input.Typ == Escape {
		g.OpenMenu()
	}
}

func (g *Game) handleInput() {
	level := g.Level
	p := level.Player
	p.Update(g)
}

func (g *Game) handleMonsters() {
	for _, monster := range g.Level.Monsters {
		monster.Update(g)
	}
}

func (g *Game) handleInvocations() {
	for _, m := range g.Level.Invocations {
		m.Update(g)
	}
}

func (g *Game) handlePnjs() {
	for _, pnj := range g.Level.Pnjs {
		pnj.Update(g)
	}
}

func (g *Game) handleProjectiles() {
	for _, projectile := range g.Level.Projectiles {
		projectile.Update(g)
	}
}

func (level *Level) OpenPortal(g *Game, pos Pos) {
	port := level.Portals[pos]
	if port != nil {
		p := level.Player
		p.X = port.PosTo.X
		p.Y = port.PosTo.Y
		g.Level = g.Levels[port.LevelTo]
		g.Level.Player = p

		g.GetEventManager().Dispatch(&Event{
			Type:    PlayerEventsType,
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelType": g.Level.Type},
			Message: "Going to level " + port.LevelTo})
	}
}

func (level *Level) AddPortal(posFrom Pos, portal *Portal) {
	if len(level.Portals) == 0 {
		level.Portals = make(map[Pos]*Portal)
	}
	level.Portals[posFrom] = portal
}
