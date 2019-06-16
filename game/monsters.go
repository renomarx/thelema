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
	monster.VisionRange = mt.VisionRange
	return monster
}

func (m *Monster) Update(g *Game) {
	level := g.Level
	p := level.Player
	if p.IsDead() {
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
			m.Move(positions[1], g)
		}
		if m.canAttackInvocation(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackInvocation(level.Invocations[positions[1].Y][positions[1].X], g)
		}
		if m.canAttackFriend(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackFriend(level.Friends[positions[1].Y][positions[1].X], g)
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
			mm := l.GetInvocation(x, y)
			if mm != nil {
				m.target = &mm.Character
				return Pos{X: x, Y: y}
			}
			f := l.GetFriend(x, y)
			if f != nil && !f.IsDead() {
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
	if isThereAPlayerCharacter(level, to) {
		return false
	}
	return true
}

func (m *Monster) Move(to Pos, g *Game) {
	lastPos := Pos{X: m.Pos.X, Y: m.Pos.Y}
	g.Level.Monsters[m.Y][m.X] = nil
	g.Level.Monsters[to.Y][to.X] = m
	m.moveFromTo(lastPos, to)
}

func (m *Monster) canAttackPlayer(to Pos, level *Level) bool {
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return true
	}
	return false
}

func (m *Monster) AttackPlayer(p *Player, game *Game) {
	m.IsAttacking = true
	for i := 0; i < CaseLen; i++ {
		m.adaptSpeed()
	}
	p.TakeDamage(game, m.CalculateAttackScore())
	m.IsAttacking = false
}

func (m *Monster) canAttackInvocation(to Pos, level *Level) bool {
	return isThereAnInvocation(level, to)
}

func (m *Monster) AttackInvocation(p *Invoked, g *Game) {
	m.IsAttacking = true
	for i := 0; i < CaseLen; i++ {
		m.adaptSpeed()
	}
	p.TakeDamage(g, m.CalculateAttackScore())
	m.IsAttacking = false
}

func (m *Monster) canAttackFriend(to Pos, level *Level) bool {
	return isThereAFriend(level, to)
}

func (m *Monster) AttackFriend(p *Friend, g *Game) {
	m.IsAttacking = true
	for i := 0; i < CaseLen; i++ {
		m.adaptSpeed()
	}
	p.TakeDamage(g, m.CalculateAttackScore())
	m.IsAttacking = false
}

func (m *Monster) TakeDamage(g *Game, damage int, c *Character) {
	if m == nil {
		return
	}
	if m.Health.Current <= 0 {
		m.Die(g)
	}
	m.Health.Current -= damage
	g.MakeExplosion(m.Pos, damage, 50)
	m.ParalyzedTime = rand.Intn(damage) * 10
	m.target = c
}

func (m *Monster) Die(g *Game) {
	m.isDead = true
	g.Level.Monsters[m.Y][m.X] = nil
	b := &Object{Rune: rune(Steak), Blocking: true}
	b.Pos = m.Pos
	g.Level.Objects[m.Y][m.X] = b
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
			return true
		}
	}
	return false
}
