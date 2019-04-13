package game

import "time"
import "math"

type Invoked struct {
	Character
	CreatedAt time.Time
	Lifetime  int
	target    *Monster
}

func (level *Level) MakeInvocation(p Pos, dir InputType, pp *PlayerPower) bool {
	np := Pos{X: p.X, Y: p.Y}
	if dir == Left {
		np.X--
	}
	if dir == Right {
		np.X++
	}
	if dir == Up {
		np.Y--
	}
	if dir == Down {
		np.Y++
	}
	if canGo(level, np) {
		m := NewFox(np, pp)
		level.Invocations[np] = m
		return true
	}
	return false
}

func NewFox(p Pos, pp *PlayerPower) *Invoked {
	monster := &Invoked{}
	monster.Rune = rune(Fox)
	monster.Name = "Invoked Fox"
	monster.Hitpoints.Init(pp.Strength)
	monster.Strength.Init(pp.Strength / 10)
	monster.Speed.Init(pp.Speed)
	monster.VisionRange = 5
	monster.ActionPoints = 0.0
	monster.Pos = p
	monster.Xb = 0
	monster.Yb = 0
	monster.CreatedAt = time.Now()
	monster.LastActionTime = time.Now()
	monster.IsMoving = false
	monster.Lifetime = pp.Lifetime

	return monster
}

func (m *Invoked) Update(g *Game) {
	level := g.Level
	if m.IsMoving {
		return
	}
	if level.Player.IsDead() {
		m.Die(g.Level)
		return
	}
	t := time.Now()
	deltaC := t.Sub(m.CreatedAt)
	if deltaC >= time.Duration(m.Lifetime)*time.Second {
		m.Die(g.Level)
		return
	}
	deltaD := t.Sub(m.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	m.ActionPoints += float64(m.Speed.Current) * delta
	monsterPos := m.getTargetPos(level)
	positions := level.astar(m.Pos, monsterPos, m)
	if len(positions) > 1 && m.ActionPoints >= 100000 { // 0.1 second
		if m.canMove(positions[1], level) {
			m.Move(positions[1], level)
		}
		if m.canAttack(positions[1], level) {
			m.Attack(level.Monsters[positions[1]], g)
		}
		m.ActionPoints = 0.0
	}
	m.LastActionTime = time.Now()
}

func (m *Invoked) getTargetPos(level *Level) Pos {
	if m.target != nil {
		_, monsterExists := level.Monsters[m.target.Pos]
		if !monsterExists {
			m.target = nil
		}
	}
	for mmpos, mm := range level.Monsters {
		if math.Abs(float64(m.Pos.X-mmpos.X)) < float64(m.VisionRange) || math.Abs(float64(m.Pos.Y-mmpos.Y)) < float64(m.VisionRange) {
			m.target = mm
			return mmpos
		}
	}
	return level.Player.Pos
}

func (m *Invoked) canMove(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	if isThereAMonster(level, to) {
		return false
	}
	return true
}

func (m *Invoked) Move(to Pos, level *Level) {
	m.IsMoving = true
	lastPos := Pos{X: m.Pos.X, Y: m.Pos.Y}
	Mux.Lock()
	delete(level.Invocations, m.Pos)
	level.Invocations[to] = m
	Mux.Unlock()
	m.moveFromTo(lastPos, to)
}

func (m *Invoked) canAttack(to Pos, level *Level) bool {
	return isThereAMonster(level, to)
}

func (m *Invoked) Attack(mm *Monster, g *Game) {
	m.IsMoving = true
	m.IsAttacking = true
	go func(m *Invoked) {
		for i := 0; i < CaseLen; i++ {
			m.adaptSpeed()
		}
		m.IsMoving = false
		m.IsAttacking = false
	}(m)
	mm.TakeDamage(g.Level, m.Strength.Current)
}

func (m *Invoked) TakeDamage(level *Level, damage int) {
	if m.Hitpoints.Current <= 0 {
		m.Die(level)
	}
	m.Hitpoints.Current -= damage
	level.MakeExplosion(m.Pos, damage, 50)
}

func (m *Invoked) Die(level *Level) {
	Mux.Lock()
	delete(level.Invocations, m.Pos)
	Mux.Unlock()
}

func (m *Invoked) CanSee(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAnInvocation(level, pos) {
		return false
	}
	if isThereAPnj(level, pos) {
		return false
	}
	if pos.Y >= 0 && pos.Y < len(level.Map) {
		if pos.X >= 0 && pos.X < len(level.Map[pos.Y]) {
			return level.Map[pos.Y][pos.X] != StoneWall && level.Map[pos.Y][pos.X] != DoorClosed
		}
	}
	return false
}
