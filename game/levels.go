package game

import "math/rand"
import "time"

type Level struct {
	Name   string
	Width  int
	Height int
	Player *Player
	Map    [][]Case
	Paused bool
	PRay   int
}

type Case struct {
	T                   Tile
	Portal              *Portal
	Object              *Object
	Effect              *Effect
	Pnj                 *Pnj
	MonstersProbability int
}

func (l *Level) GetObject(x, y int) *Object {
	if y >= 0 && y < len(l.Map) {
		if x >= 0 && x < len(l.Map[y]) {
			return l.Map[y][x].Object
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

func NewLevel() *Level {
	level := &Level{}
	level.PRay = 20
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

func (l *Level) handleMap() {
	for y := l.Player.Y - l.PRay; y < l.Player.Y+l.PRay; y++ {
		for x := l.Player.X - l.PRay; x < l.Player.X+l.PRay; x++ {
			if y >= 0 && y < len(l.Map) {
				if x >= 0 && x < len(l.Map[y]) {
					c := l.Map[y][x]
					l.handlePnj(c.Pnj)
				}
			}
		}
	}
}

func (l *Level) handlePnj(m *Pnj) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Pnj) {
			m.Update(l)
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

		EM.Dispatch(&Event{
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelName": g.Level.Name},
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

func (l *Level) MakeExplosion(p Pos, size int, lifetime int) {
	eff := NewExplosion(p, size, lifetime)
	l.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	l.Map[p.Y][p.X].Effect = nil
}

func (l *Level) MakeRangeStorm(p Pos, damages int, dir InputType, lifetime int, rg int) {
	storms := NewRangeStorm(p, damages, dir, rg)
	for _, storm := range storms {
		go l.MakeStorm(storm, lifetime)
	}
}

func (l *Level) MakeStorm(storm *Effect, lifetime int) {
	if !storm.canBe(l, storm.Pos) {
		return
	}
	l.Map[storm.Pos.Y][storm.Pos.X].Effect = storm
	time.Sleep(time.Duration(lifetime) * time.Second)
	storm.Die(l)
}

func (l *Level) MakeFlames(p Pos, damages int, lifetime int, rg int) {
	for y := p.Y - rg; y <= p.Y+rg; y++ {
		for x := p.X - rg; x <= p.X+rg; x++ {
			if x != p.X || y != p.Y {
				go l.MakeFlame(Pos{X: x, Y: y}, damages, lifetime)
			}
		}
	}
}

func (l *Level) MakeFlame(p Pos, damages int, lifetime int) {
	eff := NewFlame(p, damages)
	if !eff.canBe(l, p) {
		return
	}
	eff.TileIdx = 0
	l.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 1
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 2
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 3
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.Die(l)
}

func (l *Level) MakeEffect(p Pos, r rune, lifetime int) {
	eff := NewEffect(p, r)

	l.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	l.Map[p.Y][p.X].Effect = nil
}
