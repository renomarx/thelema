package game

import "time"
import "math/rand"

type Enemy struct {
	Character
	target *Character
}

func NewEnemy(mt *MonsterType, p Pos) *Enemy {
	monster := &Enemy{}
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
	monster.Weapon = mt.Weapon
	monster.CurrentPower = mt.Power
	return monster
}

func (level *Level) MakeEnemy(pnj *Pnj) {
	np := pnj.Pos
	y := &Enemy{}
	y.Character = pnj.Character
	y.Speed.Init(y.Speed.Current * 2)
	y.VisionRange = 12
	level.Enemies[np.Y][np.X] = y
}

func (e *Enemy) Update(g *Game) {
	if e.IsDead() {
		return
	}
	level := g.Level
	t := time.Now()
	deltaD := t.Sub(e.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	e.ActionPoints += float64(e.Speed.Current) * delta
	monsterPos := e.getTargetPos(g)
	positions := level.astar(e.Pos, monsterPos, e)
	if len(positions) > 1 && e.ActionPoints >= 100000 { // 0.1 second
		if e.canMove(positions[1], level) {
			e.Move(positions[1], g)
		}
		if e.canAttack(positions[1], level) {
			e.Attack(g, positions[1])
		}
		e.ActionPoints = 0.0
	}
	e.LastActionTime = time.Now()
}

func (e *Enemy) getTargetPos(g *Game) Pos {
	l := g.Level
	if e.target != nil {
		if e.target.IsDead() {
			e.target = nil
		} else {
			return e.target.Pos
		}
	}

	for y := e.Y - e.VisionRange; y < e.Y+e.VisionRange; y++ {
		for x := e.X - e.VisionRange; x < e.X+e.VisionRange; x++ {
			mm := l.GetInvocation(x, y)
			if mm != nil {
				e.target = &mm.Character
				return Pos{X: x, Y: y}
			}
			f := l.GetFriend(x, y)
			if f != nil && !f.IsDead() {
				e.target = &f.Character
				return Pos{X: x, Y: y}
			}
			if l.Player.Pos.X == x && l.Player.Pos.Y == y && !l.Player.IsDead() {
				return l.Player.Pos
			}
		}
	}
	return e.Pos
}

func (e *Enemy) canMove(to Pos, level *Level) bool {
	if isThereAPlayerCharacter(level, to) {
		return false
	}
	return true
}

func (e *Enemy) Move(to Pos, g *Game) {
	level := g.Level
	lastPos := Pos{X: e.Pos.X, Y: e.Pos.Y}
	level.Enemies[e.Y][e.X] = nil
	level.Enemies[to.Y][to.X] = e
	e.moveFromTo(lastPos, to)
}

func (e *Enemy) canAttack(to Pos, level *Level) bool {
	return isThereAPlayerCharacter(level, to)
}

func (e *Enemy) TakeDamage(g *Game, damage int) {
	if e.Health.Current <= 0 {
		e.Die(g)
	}
	e.Health.Current -= damage
	g.MakeExplosion(e.Pos, damage, 50)
	e.ParalyzedTime = rand.Intn(damage) * 10
}

func (e *Enemy) Die(g *Game) {
	e.isDead = true
	g.Level.Enemies[e.Y][e.X] = nil
}

func (e *Enemy) CanSee(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAMonster(level, pos) {
		return false
	}
	if isThereAnEnemy(level, pos) {
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
