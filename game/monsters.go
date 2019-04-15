package game

import "time"
import "math"

type Monster struct {
	Character
	target *Invoked
}

func NewMonster(mt *MonsterType, p Pos) *Monster {
	monster := &Monster{}
	monster.Rune = rune(mt.Tile)
	monster.Name = mt.Name
	monster.Hitpoints.Init(mt.Hitpoints)
	monster.Strength.Init(mt.Strength)
	monster.Speed.Init(mt.Speed)
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
		if m.canAttackPlayer(positions[1], level) {
			g.GetEventManager().Dispatch(&Event{Action: ActionRoar, Payload: map[string]string{"monster": string(m.Rune)}})
			m.AttackPlayer(level.Player, g)
		}
		m.ActionPoints = 0.0
	}
	m.LastActionTime = time.Now()
}

func (m *Monster) getTargetPos(g *Game) Pos {
	level := g.Level
	if m.target != nil {
		_, monsterExists := level.Invocations[m.target.Pos]
		if !monsterExists {
			m.target = nil
		}
	}
	for mmpos, mm := range level.Invocations {
		if math.Abs(float64(m.Pos.X-mmpos.X)) < float64(m.VisionRange) || math.Abs(float64(m.Pos.Y-mmpos.Y)) < float64(m.VisionRange) {
			m.target = mm
			return mmpos
		}
	}
	if math.Abs(float64(m.Pos.X-level.Player.Pos.X)) < float64(m.VisionRange) || math.Abs(float64(m.Pos.Y-level.Player.Pos.Y)) < float64(m.VisionRange) {
		return level.Player.Pos
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
	p.TakeDamage(game, m.Strength.Current)
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
	p.TakeDamage(g, m.Strength.Current)
}

func (m *Monster) TakeDamage(g *Game, damage int) {
	if m.Hitpoints.Current <= 0 {
		m.Die(g.Level)
	}
	m.Hitpoints.Current -= damage
	g.MakeExplosion(m.Pos, damage, 50)
}

func (m *Monster) Die(level *Level) {
	Mux.Lock()
	delete(level.Monsters, m.Pos)
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
