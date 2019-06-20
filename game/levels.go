package game

import "math/rand"

const LevelTypeOutdoor = "OUTDOOR"
const LevelTypeGrotto = "GROTTO"
const LevelTypeCity = "CITY"
const LevelTypeHouse = "HOUSE"

type Level struct {
	Name   string
	Width  int
	Height int
	Type   string
	Player *Player
	Map    [][]Case
	Paused bool
	PRay   int
}

type Case struct {
	T          Tile
	Portal     *Portal
	Object     *Object
	Monster    *Monster
	Effect     *Effect
	Projectile *Projectile
	Pnj        *Pnj
	Invoked    *Invoked
	Friend     *Friend
	Enemy      *Enemy
}

func (l *Level) GetMonster(x, y int) *Monster {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Monster
		}
	}
	return nil
}

func (l *Level) GetObject(x, y int) *Object {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Object
		}
	}
	return nil
}

func (l *Level) GetEffect(x, y int) *Effect {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Effect
		}
	}
	return nil
}

func (l *Level) GetProjectile(x, y int) *Projectile {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Projectile
		}
	}
	return nil
}

func (l *Level) GetPnj(x, y int) *Pnj {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Pnj
		}
	}
	return nil
}

func (l *Level) GetInvocation(x, y int) *Invoked {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Invoked
		}
	}
	return nil
}

func (l *Level) GetFriend(x, y int) *Friend {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Friend
		}
	}
	return nil
}

func (l *Level) GetEnemy(x, y int) *Enemy {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Enemy
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
	level.Width = width
	level.Height = height
	level.Map = make([][]Case, height)
	for i := range level.Map {
		level.Map[i] = make([]Case, width)
	}
}

func (g *Game) UpdateLevel() {
	input := g.input
	if g.Level.Paused {
		g.HandleInputPlayerMenu()
	} else {
		g.handleInput()
		g.handleMap()
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

func (g *Game) handleMap() {
	l := g.Level
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			if y >= 0 && y < len(l.Map) {
				if x >= 0 && x < len(l.Map[y]) {
					c := l.Map[y][x]
					g.handleMonster(c.Monster)
					g.handleInvocation(c.Invoked)
					g.handlePnj(c.Pnj)
					g.handleFriend(c.Friend)
					g.handleEnemy(c.Enemy)
					g.handleProjectile(c.Projectile)
					g.handleEffect(c.Effect)
				}
			}
		}
	}
}

func (g *Game) handleMonster(m *Monster) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Monster) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handleInvocation(m *Invoked) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Invoked) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handlePnj(m *Pnj) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Pnj) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handleFriend(m *Friend) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Friend) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handleEnemy(m *Enemy) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Enemy) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handleProjectile(m *Projectile) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Projectile) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (g *Game) handleEffect(m *Effect) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Effect) {
			m.Update(g)
			m.IsPlaying = false
		}(m)
	}
}

func (level *Level) OpenPortal(g *Game, pos Pos) {
	port := level.Map[pos.Y][pos.X].Portal
	if port != nil {
		p := level.Player
		p.X = port.PosTo.X
		p.Y = port.PosTo.Y
		levelFrom := *g.Level
		g.Level = g.Levels[port.LevelTo]
		g.Level.Player = p
		f := p.Friend
		if f != nil {
			oldP := f.Pos
			f.Pos = port.PosTo
			g.Level.Map[port.PosTo.Y][port.PosTo.X].Friend = f
			levelFrom.Map[oldP.Y][oldP.X].Friend = nil
		}

		g.GetEventManager().Dispatch(&Event{
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelType": g.Level.Type},
			Message: "Going to " + port.LevelTo})
	}
}

func (level *Level) AddPortal(posFrom Pos, portal *Portal) {
	level.Map[posFrom.Y][posFrom.X].Portal = portal
}
