package game

import "math/rand"

const LevelTypeOutdoor = "OUTDOOR"
const LevelTypeGrotto = "GROTTO"
const LevelTypeCity = "CITY"
const LevelTypeHouse = "HOUSE"

type Level struct {
	Type        string
	Portals     map[Pos]*Portal
	Player      *Player
	Map         [][]Tile
	Monsters    [][]*Monster
	Objects     [][]*Object
	Effects     [][]*Effect
	Projectiles [][]*Projectile
	Pnjs        [][]*Pnj
	Invocations [][]*Invoked
	Friends     [][]*Friend
	Enemies     [][]*Enemy
	Paused      bool
	PRay        int
}

func (l *Level) GetMonster(x, y int) *Monster {
	if y >= 0 && y < len(l.Monsters) {
		if x >= 0 && x < len(l.Monsters[y]) {
			return l.Monsters[y][x]
		}
	}
	return nil
}

func (l *Level) GetObject(x, y int) *Object {
	if y >= 0 && y < len(l.Objects) {
		if x >= 0 && x < len(l.Objects[y]) {
			return l.Objects[y][x]
		}
	}
	return nil
}

func (l *Level) GetEffect(x, y int) *Effect {
	if y >= 0 && y < len(l.Effects) {
		if x >= 0 && x < len(l.Effects[y]) {
			return l.Effects[y][x]
		}
	}
	return nil
}

func (l *Level) GetProjectile(x, y int) *Projectile {
	if y >= 0 && y < len(l.Projectiles) {
		if x >= 0 && x < len(l.Projectiles[y]) {
			return l.Projectiles[y][x]
		}
	}
	return nil
}

func (l *Level) GetPnj(x, y int) *Pnj {
	if y >= 0 && y < len(l.Pnjs) {
		if x >= 0 && x < len(l.Pnjs[y]) {
			return l.Pnjs[y][x]
		}
	}
	return nil
}

func (l *Level) GetInvocation(x, y int) *Invoked {
	if y >= 0 && y < len(l.Invocations) {
		if x >= 0 && x < len(l.Invocations[y]) {
			return l.Invocations[y][x]
		}
	}
	return nil
}

func (l *Level) GetFriend(x, y int) *Friend {
	if y >= 0 && y < len(l.Friends) {
		if x >= 0 && x < len(l.Friends[y]) {
			return l.Friends[y][x]
		}
	}
	return nil
}

func (l *Level) GetEnemy(x, y int) *Enemy {
	if y >= 0 && y < len(l.Enemies) {
		if x >= 0 && x < len(l.Enemies[y]) {
			return l.Enemies[y][x]
		}
	}
	return nil
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
	level.PRay = 100
	return level
}

func (level *Level) InitMaps(height, width int) {
	level.Map = make([][]Tile, height)
	level.Monsters = make([][]*Monster, height)
	level.Objects = make([][]*Object, height)
	level.Effects = make([][]*Effect, height)
	level.Projectiles = make([][]*Projectile, height)
	level.Pnjs = make([][]*Pnj, height)
	level.Invocations = make([][]*Invoked, height)
	level.Friends = make([][]*Friend, height)
	level.Enemies = make([][]*Enemy, height)
	for i := range level.Map {
		level.Map[i] = make([]Tile, width)
		level.Monsters[i] = make([]*Monster, width)
		level.Objects[i] = make([]*Object, width)
		level.Effects[i] = make([]*Effect, width)
		level.Projectiles[i] = make([]*Projectile, width)
		level.Pnjs[i] = make([]*Pnj, width)
		level.Invocations[i] = make([]*Invoked, width)
		level.Friends[i] = make([]*Friend, width)
		level.Enemies[i] = make([]*Enemy, width)
	}
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
		p.IsPlaying = true
		go func(p *Player) {
			p.Update(g)
			p.IsPlaying = false
		}(p)
	}
}

func (g *Game) handleMonsters() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			if y >= 0 && y < len(l.Monsters) {
				if x >= 0 && x < len(l.Monsters[y]) {
					m := l.GetMonster(x, y)
					if m != nil && !m.IsPlaying {
						m.IsPlaying = true
						go func(m *Monster) {
							m.Update(g)
							m.IsPlaying = false
						}(m)
					}
				}
			}
		}
	}
}

func (g *Game) handleInvocations() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetInvocation(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Invoked) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handlePnjs() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetPnj(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Pnj) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleFriends() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetFriend(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Friend) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleEnemies() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetEnemy(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Enemy) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleProjectiles() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetProjectile(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Projectile) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (g *Game) handleEffects() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			m := l.GetEffect(x, y)
			if m != nil && !m.IsPlaying {
				m.IsPlaying = true
				go func(m *Effect) {
					m.Update(g)
					m.IsPlaying = false
				}(m)
			}
		}
	}
}

func (level *Level) OpenPortal(g *Game, pos Pos) {
	port := level.Portals[pos]
	if port != nil {
		p := level.Player
		p.X = port.PosTo.X
		p.Y = port.PosTo.Y
		//levelFrom := *g.Level
		g.Level = g.Levels[port.LevelTo]
		g.Level.Player = p
		// TODO
		// for oldP, f := range levelFrom.Friends {
		// 	f.Pos = port.PosTo
		// 	g.Level.Friends[port.PosTo] = f
		// 	g.Mux.Lock()
		// 	delete(levelFrom.Friends, oldP)
		// 	g.Mux.Unlock()
		// }

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
