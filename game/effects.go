package game

type Effect struct {
	Object
	TileIdx int
	Damages int
}

const ExplosionSizeSmall = "SMALL"
const ExplosionSizeMedium = "MEDIUM"
const ExplosionSizeLarge = "LARGE"

func NewExplosion(p Pos, size int, lifetime int) *Effect {
	esize := ExplosionSizeSmall
	if size >= 50 {
		esize = ExplosionSizeMedium
		if size >= 100 {
			esize = ExplosionSizeLarge
		}
	}
	EM.Dispatch(&Event{Action: ActionExplode, Payload: map[string]string{"size": esize}})
	eff := &Effect{}
	eff.Pos = p
	eff.Rune = string(Explosion)
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

	return eff
}

func NewRangeStorm(p Pos, damages int, dir InputType, rg int) []*Effect {
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
	var res []*Effect
	for _, pp := range poss {
		eff := NewStorm(pp, damages, dir)
		res = append(res, eff)
	}
	return res
}

func NewStorm(p Pos, damages int, dir InputType) *Effect {
	eff := &Effect{}
	eff.Pos = p
	eff.Rune = string(Storm)
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
	return eff
}

func NewFlame(p Pos, damages int) *Effect {
	eff := &Effect{}
	eff.Pos = p
	eff.Rune = string(Flames)
	eff.Blocking = false
	eff.Damages = damages
	eff.TileIdx = 0
	return eff
}

func NewEffect(p Pos, r string) *Effect {
	eff := &Effect{}
	eff.Pos = p
	eff.Rune = r
	eff.Blocking = false
	eff.TileIdx = 0
	return eff
}

func (e *Effect) Die(l *Level) {
	l.Map[e.Y][e.X].Effect = nil
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
