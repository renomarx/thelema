package game

import "math/rand"

const LevelTypeOutdoor = "OUTDOOR"
const LevelTypeGrotto = "GROTTO"
const LevelTypeCity = "CITY"
const LevelTypeHouse = "HOUSE"

type Level struct {
	Name                string
	Width               int
	Height              int
	Type                string
	Player              *Player
	Map                 [][]Case
	Paused              bool
	PRay                int
	MonstersProbability int
}

type Case struct {
	T          Tile
	Portal     *Portal
	Object     *Object
	Effect     *Effect
	Projectile *Projectile
	Pnj        *Pnj
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
	level.PRay = 20
	level.MonstersProbability = 0 // Default
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
					g.handlePnj(c.Pnj)
					g.handleProjectile(c.Projectile)
				}
			}
		}
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

func (g *Game) handleProjectile(m *Projectile) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Projectile) {
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
		g.Level = g.Levels[port.LevelTo]
		g.Level.Player = p

		g.GetEventManager().Dispatch(&Event{
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelType": g.Level.Type},
			Message: "Going to " + port.LevelTo})
	}
}

func (level *Level) AddPortal(posFrom Pos, portal *Portal) {
	level.Map[posFrom.Y][posFrom.X].Portal = portal
}

func (g *Game) SendToLevel(fromName, pnjName, toName string) {
	from, exists := g.Levels[fromName]
	if !exists {
		panic("Level " + fromName + " does not exist")
	}
	pnj := from.SearchPnj(pnjName)
	if pnj == nil {
		panic("Pnj " + pnjName + " on level " + fromName + " does not exist")
	}
	to, exists := g.Levels[toName]
	if !exists {
		panic("Level " + toName + " does not exist")
	}
	pnj.ChangeLevel(from, to)
}

func (l *Level) SearchPnj(pnjName string) *Pnj {
	for y := 0; y < len(l.Map); y++ {
		for x := 0; x < len(l.Map[y]); x++ {
			if l.Map[y][x].Pnj != nil {
				pnj := l.Map[y][x].Pnj
				if pnj.Name == pnjName {
					return pnj
				}
			}
		}
	}
	return nil
}
