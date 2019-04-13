package game

import "time"

const ProjectileDeltaTime = 20

type Projectile struct {
	MovingObject
	Size      int
	Speed     int
	Direction InputType
}

func (level *Level) MakeEnergyball(p Pos, dir InputType, size int, speed int) {
	eb := &Projectile{}
	eb.Rune = rune(Energyball)
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

func (p *Projectile) Update(game *Game) {
	if p.IsMoving {
		return
	}
	level := game.Level

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

	p.Move(to, level)
}

func (p *Projectile) canMove(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	return level.Map[pos.Y][pos.X] != StoneWall && level.Map[pos.Y][pos.X] != DoorClosed
}

func (p *Projectile) Move(to Pos, level *Level) {
	if p.IsMoving {
		return
	}
	p.IsMoving = true
	Mux.Lock()
	delete(level.Projectiles, p.Pos)
	level.Projectiles[to] = p
	Mux.Unlock()
	p.Pos = to

	if !p.canMove(level, to) {
		p.Die(level)
		return
	}

	p.MakeDamage(level)

	go func(p *Projectile) {

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
		p.IsMoving = false
	}(p)
}

func (p *Projectile) adaptSpeed() {
	time.Sleep(time.Duration(ProjectileDeltaTime/p.Speed) * time.Millisecond)
}

func (p *Projectile) MakeDamage(level *Level) {
	Mux.Lock()
	m, ok := level.Monsters[p.Pos]
	Mux.Unlock()
	if ok {
		// There is a monster !
		m.TakeDamage(level, p.Size)
	}
}

func (p *Projectile) Die(level *Level) {
	Mux.Lock()
	delete(level.Projectiles, p.Pos)
	Mux.Unlock()
	level.MakeExplosion(p.Pos, 100, 100)
}
