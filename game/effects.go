package game

import "time"

type Effect struct {
	Object
	TileIdx int
	Damages int
}

const ExplosionSizeSmall = "SMALL"
const ExplosionSizeMedium = "MEDIUM"
const ExplosionSizeLarge = "LARGE"

func (g *Game) MakeExplosion(p Pos, size int, lifetime int) {
	level := g.Level
	// _, alreadyExists := level.Effects[p]
	// if alreadyExists {
	// 	return
	// }
	esize := ExplosionSizeSmall
	if size >= 50 {
		esize = ExplosionSizeMedium
		if size >= 100 {
			esize = ExplosionSizeLarge
		}
	}
	EM.Dispatch(&Event{Action: ActionExplode, Payload: map[string]string{"size": esize}})
	eff := &Effect{}
	eff.Rune = rune(Explosion)
	eff.Blocking = false
	idx := 0
	if size%2 == 1 {
		idx = 1
	}
	if size >= 10 {
		idx += 2
	}
	if size >= 50 {
		idx += 2
	}
	eff.TileIdx = idx

	level.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	level.Map[p.Y][p.X].Effect = nil
}

func (g *Game) MakeRangeStorm(p Pos, damages int, dir InputType, lifetime int, rg int) {
	poss := []Pos{}
	switch dir {
	case Left:
		for i := 1; i <= rg; i++ {
			poss = append(poss, Pos{X: p.X - i, Y: p.Y})
		}
	case Right:
		for i := 1; i <= rg; i++ {
			poss = append(poss, Pos{X: p.X + i, Y: p.Y})
		}
	case Up:
		for i := 1; i <= rg; i++ {
			poss = append(poss, Pos{X: p.X, Y: p.Y - i})
		}
	case Down:
		for i := 1; i <= rg; i++ {
			poss = append(poss, Pos{X: p.X, Y: p.Y + i})
		}
	}
	for _, pp := range poss {
		go g.MakeStorm(pp, damages, dir, lifetime)
	}
}

func (g *Game) MakeStorm(p Pos, damages int, dir InputType, lifetime int) {
	level := g.Level
	eff := &Effect{}
	if !eff.canBe(level, p) {
		return
	}
	eff.Pos = p
	eff.Rune = rune(Storm)
	eff.Blocking = false
	eff.Damages = damages
	switch dir {
	case Left:
		eff.TileIdx = 2
	case Right:
		eff.TileIdx = 2
	case Up:
		eff.TileIdx = 0
	case Down:
		eff.TileIdx = 0
	}
	level.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Second)
	eff.Die(g)
}

func (g *Game) MakeFlames(p Pos, damages int, lifetime int, rg int) {
	for y := p.Y - rg; y <= p.Y+rg; y++ {
		for x := p.X - rg; x <= p.X+rg; x++ {
			if x != p.X || y != p.Y {
				go g.MakeFlame(Pos{X: x, Y: y}, damages, lifetime)
			}
		}
	}
}

func (g *Game) MakeFlame(p Pos, damages int, lifetime int) {
	level := g.Level
	eff := &Effect{}
	if !eff.canBe(level, p) {
		return
	}
	eff.Pos = p
	eff.Rune = rune(Flames)
	eff.Blocking = false
	eff.Damages = damages
	eff.TileIdx = 0
	level.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 1
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 2
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.TileIdx = 3
	time.Sleep(time.Duration(lifetime) * time.Millisecond * 250)
	eff.Die(g)
}

func (g *Game) MakeEffect(p Pos, r rune, lifetime int) {
	level := g.Level
	eff := &Effect{}
	eff.Rune = r
	eff.Blocking = false
	eff.TileIdx = 0

	level.Map[p.Y][p.X].Effect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	level.Map[p.Y][p.X].Effect = nil
}

func (e *Effect) Die(g *Game) {
	g.Level.Map[e.Y][e.X].Effect = nil
}

func (e *Effect) canBe(level *Level, pos Pos) bool {
	if pos.Y < 0 || pos.Y >= len(level.Map) || pos.X < 0 || pos.X >= len(level.Map[pos.Y]) {
		return false
	}
	if isThereABlockingObject(level, pos) {
		return false
	}
	return true
}
