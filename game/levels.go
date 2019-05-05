package game

import "math/rand"

const LevelTypeOutdoor = "OUTDOOR"
const LevelTypeGrotto = "GROTTO"
const LevelTypeCity = "CITY"
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
	PRay        int
}

func (l *Level) GetRandomFreePos() *Pos {
	x := rand.Intn(len(l.Map[0]))
	y := rand.Intn(len(l.Map))
	pos := Pos{X: x, Y: y}
	i := 0
	for !canGo(l, pos) && i < 333 {
		x := rand.Intn(len(l.Map[0]))
		y := rand.Intn(len(l.Map))
		pos = Pos{X: x, Y: y}
		i++
	}
	if i >= 333 {
		return nil
	}
	return &pos
}

func NewLevel(levelType string) *Level {
	level := &Level{}
	level.Type = levelType
	level.Monsters = make(map[Pos]*Monster)
	level.Objects = make(map[Pos]*Object)
	level.Effects = make(map[Pos]*Effect)
	level.Projectiles = make(map[Pos]*Projectile)
	level.Pnjs = make(map[Pos]*Pnj)
	level.Invocations = make(map[Pos]*Invoked)
	level.PRay = 100
	return level
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
		g.handleEffects()
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
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m, e := l.Monsters[Pos{X: x, Y: y}]
			if e {
				m.Update(g)
			}
		}
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

func (g *Game) handleEffects() {
	for _, eff := range g.Level.Effects {
		eff.Update(g)
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
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelType": g.Level.Type},
			Message: "Going to " + port.LevelTo})
	}
}

func (level *Level) AddPortal(posFrom Pos, portal *Portal) {
	if len(level.Portals) == 0 {
		level.Portals = make(map[Pos]*Portal)
	}
	level.Portals[posFrom] = portal
}
