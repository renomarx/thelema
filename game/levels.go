package game

import "math/rand"
import "time"

type Level struct {
	Name       string
	Width      int
	Height     int
	Player     *Player
	Map        [][][]Case
	Paused     bool
	PRay       int
	Discovered bool
	Bestiary   []*MonsterType
}

type Case struct {
	T                   Tile
	Portal              *Portal
	Object              *Object
	Effect              *Effect
	Npc                 *Npc
	MonstersProbability int
}

func (l *Level) GetObject(p Pos) *Object {
	x := p.X
	y := p.Y
	z := p.Z
	if y >= 0 && y < len(l.Map[z]) {
		if x >= 0 && x < len(l.Map[z][y]) {
			return l.Map[z][y][x].Object
		}
	}
	return nil
}

func (l *Level) GetNpc(p Pos) *Npc {
	x := p.X
	y := p.Y
	z := p.Z
	if y >= 0 && y < len(l.Map[z]) {
		if x >= 0 && x < len(l.Map[z][y]) {
			return l.Map[z][y][x].Npc
		}
	}
	return nil
}

func (l *Level) GetRandomFreePos(z int) *Pos {
	if len(l.Map[z]) == 0 {
		return nil
	}
	y := rand.Intn(len(l.Map[z]))
	if len(l.Map[z][y]) == 0 {
		return nil
	}
	x := rand.Intn(len(l.Map[z][y]))
	pos := Pos{X: x, Y: y, Z: z}
	i := 0
	for !canGo(l, pos) && i < 333 {
		y := rand.Intn(len(l.Map[z]))
		x := rand.Intn(len(l.Map[z][y]))
		pos = Pos{X: x, Y: y, Z: z}
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

func (l *Level) handleMap() {
	minZ, maxZ := l.GetZBounds()
	for z := minZ; z < maxZ; z++ {
		minY, maxY := l.GetYBounds(z)
		for y := minY; y < maxY; y++ {
			minX, maxX := l.GetXBounds(z, y)
			for x := minX; x < maxX; x++ {
				c := l.Map[z][y][x]
				l.handleNpc(c.Npc)
			}
		}
	}
}

func (l *Level) GetZBounds() (int, int) {
	return 0, len(l.Map)
}

func (l *Level) GetYBounds(z int) (int, int) {
	minY := 0
	if l.Player.Y-l.PRay > minY {
		minY = l.Player.Y - l.PRay
	}
	maxY := len(l.Map[z])
	if l.Player.Y+l.PRay < maxY {
		maxY = l.Player.Y + l.PRay
	}
	return minY, maxY
}

func (l *Level) GetXBounds(z, y int) (int, int) {
	minX := 0
	if l.Player.X-l.PRay > minX {
		minX = l.Player.X - l.PRay
	}
	maxX := len(l.Map[z][y])
	if l.Player.X+l.PRay < maxX {
		maxX = l.Player.X + l.PRay
	}
	return minX, maxX
}

func (l *Level) handleNpc(m *Npc) {
	if m != nil && !m.IsPlaying {
		m.IsPlaying = true
		go func(m *Npc) {
			m.Update(l)
			m.IsPlaying = false
		}(m)
	}
}

func (level *Level) OpenPortal(g *Game, pos Pos) {
	port := level.Map[pos.Z][pos.Y][pos.X].Portal
	if port != nil {
		p := level.Player
		if port.Key != "" {
			if !p.Inventory.HasKey(port.Key) {
				EM.Dispatch(&Event{Message: "Clé nécessaire: " + port.Key + " pour ouvrir cette porte."})
				return
			}
		}
		p.X = port.PosTo.X
		p.Y = port.PosTo.Y
		g.Level = g.Levels[port.LevelTo]
		g.Level.Discovered = true
		g.Level.Player = p

		EM.Dispatch(&Event{
			Action:  ActionChangeLevel,
			Payload: map[string]string{"levelName": g.Level.Name},
			Message: "Entré à " + port.LevelTo})
	}
}

func (level *Level) AddPortal(posFrom Pos, portal *Portal) {
	level.Map[posFrom.Z][posFrom.Y][posFrom.X].Portal = portal
}

func (g *Game) SendToLevel(fromName, npcName, toName string) {
	from, exists := g.Levels[fromName]
	if !exists {
		panic("Level " + fromName + " does not exist")
	}
	npc := from.SearchNpc(npcName)
	if npc == nil {
		panic("Npc " + npcName + " on level " + fromName + " does not exist")
	}
	to, exists := g.Levels[toName]
	if !exists {
		panic("Level " + toName + " does not exist")
	}
	npc.ChangeLevel(from, to)
}

func (l *Level) SearchNpc(npcName string) *Npc {
	minZ, maxZ := l.GetZBounds()
	for z := minZ; z < maxZ; z++ {
		minY, maxY := l.GetYBounds(z)
		for y := minY; y < maxY; y++ {
			minX, maxX := l.GetXBounds(z, y)
			for x := minX; x < maxX; x++ {
				if l.Map[z][y][x].Npc != nil {
					npc := l.Map[z][y][x].Npc
					if npc.Name == npcName {
						return npc
					}
				}
			}
		}
	}
	return nil
}

func (l *Level) MakeExplosion(p Pos, size int, lifetime int) {
	eff := NewExplosion(p, size, lifetime)
	l.Map[p.Z][p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	l.Map[p.Z][p.Y][p.X].Effect = nil
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
	l.Map[storm.Pos.Z][storm.Pos.Y][storm.Pos.X].Effect = storm
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
	l.Map[p.Z][p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 1
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 2
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 3
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.Die(l)
}

func (l *Level) MakeEffect(p Pos, r string, lifetime int) {
	eff := NewEffect(p, r)

	l.Map[p.Z][p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	l.Map[p.Z][p.Y][p.X].Effect = nil
}

func (g *Game) DiscoverLevel(levelName string) {
	level, e := g.Levels[levelName]
	if e {
		level.Discovered = true
	}
}
