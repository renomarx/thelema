package game

import "time"
import "math/rand"

type Monster struct {
	Character
	target *Character
}

func NewMonster(mt *MonsterType, p Pos) *Monster {
	monster := &Monster{}
	monster.Rune = rune(mt.Tile)
	monster.Name = mt.Name
	monster.Health.Init(mt.Health)
	monster.Energy.Init(mt.Energy)
	monster.Speed.Init(mt.Speed)
	monster.Strength.Init(mt.Stats)
	monster.Dexterity.Init(mt.Stats)
	monster.Will.Init(mt.Stats)
	monster.Intelligence.Init(mt.Stats)
	monster.Luck.Init(mt.Luck)
	monster.Beauty.Init(rand.Intn(20))
	monster.ActionPoints = 0.0
	monster.Pos = p
	monster.Xb = 0
	monster.Yb = 0
	monster.LastActionTime = time.Now()
	monster.IsMoving = false
	monster.VisionRange = mt.VisionRange
	return monster
}

func (m *Monster) Update(g *Game) {
	level := g.Level
	p := level.Player
	if m.IsMoving || p.IsDead() {
		return
	}
	t := time.Now()
	deltaD := t.Sub(m.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	m.ActionPoints += float64(m.Speed.Current) * delta
	playerPos := m.getTargetPos(g)
	positions := level.astar(m.Pos, playerPos, m)
	if len(positions) > 1 && m.ActionPoints >= 100000 { // 0.1 second
		if m.canMove(positions[1], level) {
			m.Move(positions[1], level)
		}
		if m.canAttackInvocation(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackInvocation(level.Invocations[positions[1]], g)
		}
		if m.canAttackFriend(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackFriend(level.Friends[positions[1]], g)
		}
		if m.canAttackPlayer(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackPlayer(level.Player, g)
		}
		m.ActionPoints = 0.0
	}
	m.LastActionTime = time.Now()
}

func (m *Monster) getTargetPos(g *Game) Pos {
	l := g.Level
	if m.target != nil {
		if m.target.IsDead() {
			m.target = nil
		} else {
			return m.target.Pos
		}
	}

	for y := m.Y - m.VisionRange; y < m.Y+m.VisionRange; y++ {
		for x := m.X - m.VisionRange; x < m.X+m.VisionRange; x++ {
			mm, e := l.Invocations[Pos{X: x, Y: y}]
			if e {
				m.target = &mm.Character
				return Pos{X: x, Y: y}
			}
			f, ef := l.Friends[Pos{X: x, Y: y}]
			if ef {
				m.target = &f.Character
				return Pos{X: x, Y: y}
			}
			if l.Player.Pos.X == x && l.Player.Pos.Y == y {
				return l.Player.Pos
			}
		}
	}
	return m.Pos
}

func (m *Monster) canMove(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	if isThereAnInvocation(level, to) {
		return false
	}
	if isThereAFriend(level, to) {
		return false
	}
	return true
}

func (m *Monster) Move(to Pos, level *Level) {
	m.IsMoving = true
	lastPos := Pos{X: m.Pos.X, Y: m.Pos.Y}
	Mux.Lock()
	delete(level.Monsters, m.Pos)
	level.Monsters[to] = m
	Mux.Unlock()
	m.moveFromTo(lastPos, to)
}

func (m *Monster) canAttackPlayer(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return true
	}
	return false
}

func (m *Monster) AttackPlayer(p *Player, game *Game) {
	m.IsMoving = true
	m.IsAttacking = true
	go func(m *Monster) {
		for i := 0; i < CaseLen; i++ {
			m.adaptSpeed()
		}
		m.IsMoving = false
		m.IsAttacking = false
	}(m)
	p.TakeDamage(game, m.CalculateAttackScore())
}

func (m *Monster) canAttackInvocation(to Pos, level *Level) bool {
	return isThereAnInvocation(level, to)
}

func (m *Monster) AttackInvocation(p *Invoked, g *Game) {
	m.IsMoving = true
	m.IsAttacking = true
	go func(m *Monster) {
		for i := 0; i < CaseLen; i++ {
			m.adaptSpeed()
		}
		m.IsMoving = false
		m.IsAttacking = false
	}(m)
	p.TakeDamage(g, m.CalculateAttackScore())
}

func (m *Monster) canAttackFriend(to Pos, level *Level) bool {
	return isThereAFriend(level, to)
}

func (m *Monster) AttackFriend(p *Friend, g *Game) {
	m.IsMoving = true
	m.IsAttacking = true
	go func(m *Monster) {
		for i := 0; i < CaseLen; i++ {
			m.adaptSpeed()
		}
		m.IsMoving = false
		m.IsAttacking = false
	}(m)
	p.TakeDamage(g, m.CalculateAttackScore())
}

func (m *Monster) TakeDamage(g *Game, damage int, c *Character) {
	if m.Health.Current <= 0 {
		m.Die(g.Level)
	}
	m.Health.Current -= damage
	g.MakeExplosion(m.Pos, damage, 50)
	m.ParalyzedTime = rand.Intn(damage) * 10
	m.target = c
}

func (m *Monster) Die(level *Level) {
	Mux.Lock()
	delete(level.Monsters, m.Pos)
	b := &Object{Rune: rune(Steak), Blocking: true}
	b.Pos = m.Pos
	level.Objects[m.Pos] = b
	Mux.Unlock()
}

func (m *Monster) CanSee(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAMonster(level, pos) {
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
