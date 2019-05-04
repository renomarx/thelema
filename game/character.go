package game

import "time"

const CharacterDeltaTime = 100

const CharacteristicXpMultiplier = 2

type Characteristic struct {
	Initial int
	Current int
	Xp      int
}

func (ch *Characteristic) Init(value int) {
	ch.Initial = value
	ch.Current = value
	ch.Xp = 0
}

func (ch *Characteristic) RaiseXp(value int, g *Game) {
	ch.Xp += value
	if ch.Xp >= ch.Initial*CharacteristicXpMultiplier {
		g.GetEventManager().Dispatch(&Event{Action: ActionCharacteristicUp})
		ch.Initial += ch.Initial / 20
		ch.Current = ch.Initial
		ch.Xp = 0
	}
}

type Fighter struct {
	IsAttacking      bool
	IsPowerAttacking bool
	AttackPos        int
}

const VoiceMaleStandard = "MALE_STANDARD"
const VoiceFemaleStandard = "FEMALE_STANDARD"

type Talker struct {
	IsTalking bool
	Dialog    *Dialog
	Voice     string
}

type Character struct {
	MovingObject
	Fighter
	LookAt            InputType
	Name              string
	Health            Characteristic
	Energy            Characteristic
	Speed             Characteristic
	RegenerationSpeed Characteristic
	Strength          Characteristic
	Dexterity         Characteristic
	Beauty            Characteristic
	Will              Characteristic
	Intelligence      Characteristic
	Charisma          Characteristic
	Luck              Characteristic
	Affinity          string
	ActionPoints      float64
	LastActionTime    time.Time
	isDead            bool
	VisionRange       int
	Weapon            *Weapon
	Powers            map[string]*PlayerPower
	CurrentPower      *PlayerPower
}

func (c *Character) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/c.Speed.Current) * time.Millisecond)
}

func (c *Character) moveLeft() {
	for c.Xb = -1 * CaseLen; c.Xb < 0; c.Xb++ {
		c.adaptSpeed()
	}
	c.IsMoving = false
}

func (c *Character) moveRight() {
	for c.Xb = CaseLen; c.Xb > 0; c.Xb-- {
		c.adaptSpeed()
	}
	c.IsMoving = false
}

func (c *Character) moveUp() {
	for c.Yb = -1 * CaseLen; c.Yb < 0; c.Yb++ {
		c.adaptSpeed()
	}
	c.IsMoving = false
}

func (c *Character) moveDown() {
	for c.Yb = CaseLen; c.Yb > 0; c.Yb-- {
		c.adaptSpeed()
	}
	c.IsMoving = false
}

func (c *Character) IsDead() bool {
	return c.isDead
}

func (c *Character) Attack(g *Game, posToAttack Pos) bool {
	level := g.Level
	c.IsMoving = true
	c.IsAttacking = true
	go func(c *Character) {
		for c.AttackPos = 0; c.AttackPos < CaseLen; c.AttackPos++ {
			if c.Weapon != nil {
				c.Weapon.adaptSpeed()
			} else {
				c.adaptSpeed()
			}
		}
		c.AttackPos = 0
		c.IsMoving = false
		c.IsAttacking = false
	}(c)
	if isThereAMonster(level, posToAttack) {
		m := level.Monsters[posToAttack]
		m.TakeDamage(g, c.CalculateAttackScore())
		c.Strength.RaiseXp(2, g)
		return true
	}
	return false
}

func (c *Character) CalculateAttackScore() int {
	score := float64((c.Strength.Current+c.Dexterity.Current)/2) * (1.0 + float64(c.Luck.Current)/100.0)
	iscore := int(score)
	if c.Weapon != nil {
		iscore += c.Weapon.Damages
	}
	return iscore
}

func (c *Character) PowerAttack(g *Game) {
	if c.Energy.Current > 0 {
		c.IsMoving = true
		c.IsPowerAttacking = true
		go func(c *Character, g *Game) {
			for c.AttackPos = 0; c.AttackPos < CaseLen; c.AttackPos++ {
				c.adaptSpeed()
			}
			c.IsMoving = false
			c.IsPowerAttacking = false
			c.Energy.RaiseXp(10, g)
			switch c.CurrentPower.Type {
			case PowerEnergyBall:
				g.GetEventManager().Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerEnergyBall}})
				g.Level.MakeEnergyball(c.Pos, c.LookAt, c.CurrentPower.Strength, c.CurrentPower.Speed)
				c.Energy.Current -= c.CurrentPower.Energy
				c.Will.RaiseXp(1, g)
			case PowerInvocation:
				g.GetEventManager().Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerInvocation}})
				if g.Level.MakeInvocation(c.Pos, c.LookAt, c.CurrentPower) {
					c.Energy.Current -= c.CurrentPower.Energy
					c.Will.RaiseXp(1, g)
					c.Charisma.RaiseXp(1, g)
				}
			default:
			}
		}(c, g)
	}
}
