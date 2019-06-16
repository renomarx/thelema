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
	Friends     map[Pos]*Friend
	Enemies     map[Pos]*Enemy
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
	level.Friends = make(map[Pos]*Friend)
	level.Enemies = make(map[Pos]*Enemy)
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
		g.handleFriends()
		g.handleEnemies()
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
	if !p.IsPlaying {
		go func(p *Player) {
			p.IsPlaying = true
			p.Update(g)
			p.IsPlaying = false
		}(p)
	}
}

func (g *Game) handleMonsters() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m, e := l.Monsters[Pos{X: x, Y: y}]
			if e && !m.IsPlaying {
				go func(m *Monster) {
					m.IsPlaying = true
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleInvocations() {
	for _, m := range g.Level.Invocations {
		if !m.IsPlaying {
			go func(m *Invoked) {
				m.IsPlaying = true
				m.Update(g)
				m.IsPlaying = false
			}(m)
		}
	}
}

func (g *Game) handlePnjs() {
	for _, pnj := range g.Level.Pnjs {
		if !pnj.IsPlaying {
			go func(pnj *Pnj) {
				pnj.IsPlaying = true
				pnj.Update(g)
				pnj.IsPlaying = false
			}(pnj)
		}
	}
}

func (g *Game) handleFriends() {
	for _, m := range g.Level.Friends {
		if !m.IsPlaying {
			go func(m *Friend) {
				m.IsPlaying = true
				m.Update(g)
				m.IsPlaying = false
			}(m)
		}
	}
}

func (g *Game) handleEnemies() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m, e := l.Enemies[Pos{X: x, Y: y}]
			if e && !m.IsPlaying {
				go func(m *Enemy) {
					m.IsPlaying = true
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleProjectiles() {
	for _, projectile := range g.Level.Projectiles {
		if !projectile.IsPlaying {
			go func(projectile *Projectile) {
				projectile.IsPlaying = true
				projectile.Update(g)
				projectile.IsPlaying = false
			}(projectile)
		}
	}
}

func (g *Game) handleEffects() {
	for _, eff := range g.Level.Effects {
		if !eff.IsPlaying {
			go func(eff *Effect) {
				eff.IsPlaying = true
				eff.Update(g)
				eff.IsPlaying = false
			}(eff)
		}
	}
}

func (level *Level) OpenPortal(g *Game, pos Pos) {
	port := level.Portals[pos]
	if port != nil {
		p := level.Player
		p.X = port.PosTo.X
		p.Y = port.PosTo.Y
		levelFrom := *g.Level
		g.Level = g.Levels[port.LevelTo]
		g.Level.Player = p
		for oldP, f := range levelFrom.Friends {
			f.Pos = port.PosTo
			g.Level.Friends[port.PosTo] = f
			g.Mux.Lock()
			delete(levelFrom.Friends, oldP)
			g.Mux.Unlock()
		}

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
