package game

import "time"
import "math/rand"

type Friend struct {
	Character
	target *Character
}

func (level *Level) MakeFriend(pnj *Pnj) {
	np := pnj.Pos
	f := &Friend{}
	f.Character = pnj.Character
	f.Speed.Init(f.Speed.Current * 2)
	f.VisionRange = 7
	level.Friends[np] = f
}

func (m *Friend) Update(g *Game) {
	if m.IsDead() {
		return
	}
	level := g.Level
	t := time.Now()
	deltaD := t.Sub(m.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	m.ActionPoints += float64(m.Speed.Current) * delta
	monsterPos := m.getTargetPos(level)
	positions := level.astar(m.Pos, monsterPos, m)
	if len(positions) > 1 && m.ActionPoints >= 100000 { // 0.1 second
		if m.canMove(positions[1], level) {
			m.Move(positions[1], g)
		}
		if m.canAttack(positions[1], level) {
			m.Attack(g, positions[1])
		}
		m.ActionPoints = 0.0
	}
	m.LastActionTime = time.Now()
}

func (m *Friend) getTargetPos(l *Level) Pos {
	if m.target != nil {
		if m.target.IsDead() {
			m.target = nil
		} else {
			return m.target.Pos
		}
	}
	for y := m.Y - m.VisionRange; y < m.Y+m.VisionRange; y++ {
		for x := m.X - m.VisionRange; x < m.X+m.VisionRange; x++ {
			mm, e := l.Monsters[Pos{X: x, Y: y}]
			if e {
				m.target = &mm.Character
				return Pos{X: x, Y: y}
			}
			n, en := l.Enemies[Pos{X: x, Y: y}]
			if en && !n.IsDead() {
				m.target = &n.Character
				return Pos{X: x, Y: y}
			}
		}
	}
	return l.Player.Pos
}

func (m *Friend) canMove(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	if isThereAnEnemyCharacter(level, to) {
		return false
	}
	return true
}

func (m *Friend) Move(to Pos, g *Game) {
	lastPos := Pos{X: m.Pos.X, Y: m.Pos.Y}
	g.Mux.Lock()
	delete(g.Level.Friends, m.Pos)
	g.Level.Friends[to] = m
	g.Mux.Unlock()
	m.moveFromTo(lastPos, to)
}

func (m *Friend) canAttack(to Pos, level *Level) bool {
	return isThereAnEnemyCharacter(level, to)
}

func (m *Friend) TakeDamage(g *Game, damage int) {
	if m.Health.Current <= 0 {
		m.Die(g)
	}
	m.Health.Current -= damage
	g.MakeExplosion(m.Pos, damage, 50)
	m.ParalyzedTime = rand.Intn(damage) * 10
}

func (m *Friend) Die(g *Game) {
	m.isDead = true
	g.Mux.Lock()
	delete(g.Level.Friends, m.Pos)
	g.Mux.Unlock()
}

func (m *Friend) CanSee(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAPnj(level, pos) {
		return false
	}
	if isThereAFriend(level, pos) {
		return false
	}

	return isInsideMap(level, pos)
}
