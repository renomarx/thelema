package game

import (
	"math/rand"
	"time"
)

type Invoked struct {
	Character
	CreatedAt time.Time
	Lifetime  int
	target    *Character
}

func (g *Game) MakeInvocation(p Pos, dir InputType, pp *PlayerPower) bool {
	level := g.Level
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
		m := NewFox(np, pp.Lifetime)
		level.Map[np.Y][np.X].Invoked = m
		return true
	}
	return false
}

func NewFox(p Pos, lifetime int) *Invoked {
	monster := &Invoked{}
	monster.Rune = rune(Fox)
	monster.Name = "Invoked Fox"
	monster.Health.Init(200)
	monster.Energy.Init(200)
	monster.Strength.Init(30)
	monster.Dexterity.Init(30)
	monster.Will.Init(20)
	monster.Intelligence.Init(20)
	monster.Luck.Init(20)
	monster.Beauty.Init(0)
	monster.Speed.Init(10)
	monster.VisionRange = 5
	monster.ActionPoints = 0.0
	monster.Pos = p
	monster.Xb = 0
	monster.Yb = 0
	monster.CreatedAt = time.Now()
	monster.LastActionTime = time.Now()
	monster.Lifetime = lifetime

	return monster
}

func (m *Invoked) Update(g *Game) {
	level := g.Level
	if level.Player.IsDead() {
		m.Die(g)
		return
	}
	t := time.Now()
	deltaC := t.Sub(m.CreatedAt)
	if deltaC >= time.Duration(m.Lifetime)*time.Second {
		m.Die(g)
		return
	}
	deltaD := t.Sub(m.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	m.ActionPoints += float64(m.Speed.Current) * delta
	monsterPos := m.getTargetPos(level)
	positions := level.astar(m.Pos, monsterPos, m)
	if len(positions) > 1 && m.ActionPoints >= 100000 { // 0.1 second
		if m.canMove(positions[1], level) {
			m.Move(positions[1], g)
		}
		if m.canAttackEnemy(positions[1], level) {
			m.AttackEnemy(level.Map[positions[1].Y][positions[1].X].Enemy, g)
		}
		if m.canAttackMonster(positions[1], level) {
			m.AttackMonster(level.Map[positions[1].Y][positions[1].X].Monster, g)
		}
		m.ActionPoints = 0.0
	}
	m.LastActionTime = time.Now()
}

func (m *Invoked) getTargetPos(l *Level) Pos {
	if m.target != nil {
		if m.target.IsDead() {
			m.target = nil
		} else {
			return m.target.Pos
		}
	}
	for y := m.Y - m.VisionRange; y < m.Y+m.VisionRange; y++ {
		for x := m.X - m.VisionRange; x < m.X+m.VisionRange; x++ {
			mm := l.GetMonster(x, y)
			if mm != nil {
				m.target = &mm.Character
				return Pos{X: x, Y: y}
			}
			n := l.GetEnemy(x, y)
			if n != nil && !n.IsDead() {
				m.target = &n.Character
				return Pos{X: x, Y: y}
			}
		}
	}
	return l.Player.Pos
}

func (m *Invoked) canMove(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	if isThereAnEnemyCharacter(level, to) {
		return false
	}
	return true
}

func (m *Invoked) Move(to Pos, g *Game) {
	lastPos := Pos{X: m.Pos.X, Y: m.Pos.Y}
	g.Level.Map[m.Y][m.X].Invoked = nil
	g.Level.Map[to.Y][to.X].Invoked = m
	m.moveFromTo(lastPos, to)
}

func (m *Invoked) canAttackMonster(to Pos, level *Level) bool {
	return isThereAMonster(level, to)
}

func (m *Invoked) canAttackEnemy(to Pos, level *Level) bool {
	return isThereAnEnemy(level, to)
}

func (m *Invoked) AttackMonster(mm *Monster, g *Game) {
	m.IsAttacking = true
	for i := 0; i < CaseLen; i++ {
		m.adaptSpeed()
	}
	m.IsAttacking = false
	mm.TakeDamage(g, m.CalculateAttackScore(), &m.Character)
}

func (m *Invoked) AttackEnemy(mm *Enemy, g *Game) {
	m.IsAttacking = true
	for i := 0; i < CaseLen; i++ {
		m.adaptSpeed()
	}
	m.IsAttacking = false
	mm.TakeDamage(g, m.CalculateAttackScore())
}

func (m *Invoked) TakeDamage(g *Game, damage int) {
	if m.Health.Current <= 0 {
		m.Die(g)
	}
	m.Health.Current -= damage
	g.MakeExplosion(m.Pos, damage, 50)
	m.ParalyzedTime = rand.Intn(damage) * 10
}

func (m *Invoked) Die(g *Game) {
	m.isDead = true
	g.Level.Map[m.Y][m.X].Invoked = nil
}

func (m *Invoked) CanSee(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAnInvocation(level, pos) {
		return false
	}
	if isThereAFriend(level, pos) {
		return false
	}
	if isThereAPnj(level, pos) {
		return false
	}

	return isInsideMap(level, pos)
}
