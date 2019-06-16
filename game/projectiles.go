package game

import "time"

const ProjectileDeltaTime = 20

type Projectile struct {
	MovingObject
	Size      int
	Speed     int
	Direction InputType
	From      *Character
}

func (level *Level) MakeEnergyball(p Pos, dir InputType, size int, from *Character) {
	eb := &Projectile{}
	eb.From = from
	eb.Rune = rune(Energyball)
	eb.Blocking = false
	eb.Size = size
	eb.Speed = 10
	eb.Pos = p
	eb.Direction = dir
	if eb.Direction == Left {
		eb.Xb = 32
	}
	if eb.Direction == Right {
		eb.Xb = -32
	}
	if eb.Direction == Up {
		eb.Yb = 32
	}
	if eb.Direction == Down {
		eb.Yb = -32
	}

	level.Projectiles[p] = eb
}

func (level *Level) MakeArrow(p Pos, dir InputType, size int, speed int, from *Character) {
	eb := &Projectile{}
	eb.From = from
	eb.Rune = rune(Arrow)
	eb.Blocking = false
	eb.Size = size
	eb.Speed = speed
	eb.Pos = p
	eb.Direction = dir
	if eb.Direction == Left {
		eb.Xb = 32
	}
	if eb.Direction == Right {
		eb.Xb = -32
	}
	if eb.Direction == Up {
		eb.Yb = 32
	}
	if eb.Direction == Down {
		eb.Yb = -32
	}

	level.Projectiles[p] = eb
}

func (p *Projectile) Update(g *Game) {
	to := p.Pos
	if p.Direction == Left {
		to.X--
	}
	if p.Direction == Right {
		to.X++
	}
	if p.Direction == Up {
		to.Y--
	}
	if p.Direction == Down {
		to.Y++
	}

	p.Move(to, g)
}

func (p *Projectile) canMove(level *Level, pos Pos) bool {
	if !isInsideMap(level, pos) {
		return false
	}
	if isThereABlockingObject(level, pos) {
		return false
	}
	return true
}

func (p *Projectile) Move(to Pos, g *Game) {
	level := g.Level
	g.Mux.Lock()
	delete(level.Projectiles, p.Pos)
	level.Projectiles[to] = p
	g.Mux.Unlock()
	p.Pos = to

	if !p.canMove(level, to) {
		p.Die(g)
		return
	}

	p.MakeDamage(g)

	if p.Direction == Right {
		for p.Xb = CaseLen; p.Xb > 0; p.Xb-- {
			p.adaptSpeed()
		}
	}
	if p.Direction == Left {
		for p.Xb = -1 * CaseLen; p.Xb < 0; p.Xb++ {
			p.adaptSpeed()
		}
	}
	if p.Direction == Down {
		for p.Yb = CaseLen; p.Yb > 0; p.Yb-- {
			p.adaptSpeed()
		}
	}
	if p.Direction == Up {
		for p.Yb = -1 * CaseLen; p.Yb < 0; p.Yb++ {
			p.adaptSpeed()
		}
	}
}

func (p *Projectile) adaptSpeed() {
	time.Sleep(time.Duration(ProjectileDeltaTime/p.Speed) * time.Millisecond)
}

func (p *Projectile) MakeDamage(g *Game) {
	level := g.Level
	g.Mux.Lock()
	m, ok := level.Monsters[p.Pos]
	g.Mux.Unlock()
	if ok {
		// There is a monster !
		m.TakeDamage(g, p.Size, p.From)
		p.Die(g)
	}
	g.Mux.Lock()
	e, ok := level.Enemies[p.Pos]
	g.Mux.Unlock()
	if ok {
		// There is an annemy !
		e.TakeDamage(g, p.Size)
		p.Die(g)
	}
}

func (p *Projectile) Die(g *Game) {
	g.Mux.Lock()
	delete(g.Level.Projectiles, p.Pos)
	g.Mux.Unlock()
	g.MakeExplosion(p.Pos, 100, 100)
}
