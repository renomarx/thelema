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

func (ch *Characteristic) Add(value int) {
	ch.Current += value
	if ch.Current > ch.Initial {
		ch.Current = ch.Initial
	}
}

func (ch *Characteristic) RaiseXp(value int) {
	ch.Xp += value
	if ch.Xp >= ch.Initial*CharacteristicXpMultiplier {
		ch.Initial += ch.Initial / 20
		ch.Current = ch.Initial
		ch.Xp = 0
	}
}

const VoiceMaleStandard = "MALE_STANDARD"
const VoiceFemaleStandard = "FEMALE_STANDARD"

type Talker struct {
	IsTalking bool
	Dialog    *Dialog
	Voice     string
	Talkable  bool
}

type Character struct {
	MovingObject
	LookAt               InputType
	Name                 string
	Health               Characteristic
	Energy               Characteristic
	Speed                Characteristic
	RegenerationSpeed    Characteristic
	Strength             Characteristic
	Dexterity            Characteristic
	Beauty               Characteristic
	Will                 Characteristic
	Intelligence         Characteristic
	Charisma             Characteristic
	Luck                 Characteristic
	Aggressiveness       Characteristic
	Affinity             string
	ActionPoints         float64
	LastActionTime       time.Time
	Dead                 bool
	VisionRange          int
	Powers               map[string]*PlayerPower
	CurrentPower         *PlayerPower
	Attacks              []*Attack
	IsPowerUsing         bool
	isAttacking          bool
	AttackPos            int
	PowerPos             int
	LastRegenerationTime time.Time
	damagesTaken         int
}

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetHealth() Characteristic {
	return c.Health
}

func (c *Character) GetEnergy() Characteristic {
	return c.Energy
}

func (c *Character) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/c.Speed.Current) * time.Millisecond)
}

func (c *Character) DoLookAt(to Pos) {
	from := c.Pos
	if from.Y == to.Y {
		if from.X < to.X {
			c.LookAt = Right
		} else if from.X > to.X {
			c.LookAt = Left
		}
	}
	if from.X == to.X {
		if from.Y < to.Y {
			c.LookAt = Down
		} else if from.Y > to.Y {
			c.LookAt = Up
		}
	}
}

func (c *Character) moveFromTo(from Pos, to Pos) {
	c.Pos = to
	if from.Y == to.Y {
		if from.X < to.X {
			c.LookAt = Right
			c.Xb = CaseLen
			c.moveRight()
		} else if from.X > to.X {
			c.LookAt = Left
			c.Xb = -1 * CaseLen
			c.moveLeft()
		}
	}
	if from.X == to.X {
		if from.Y < to.Y {
			c.LookAt = Down
			c.Yb = CaseLen
			c.moveDown()
		} else if from.Y > to.Y {
			c.LookAt = Up
			c.Yb = -1 * CaseLen
			c.moveUp()
		}
	}
}

func (c *Character) moveLeft() {
	for c.Xb = -1 * CaseLen; c.Xb < 0; c.Xb++ {
		c.adaptSpeed()
	}
}

func (c *Character) moveRight() {
	for c.Xb = CaseLen; c.Xb > 0; c.Xb-- {
		c.adaptSpeed()
	}
}

func (c *Character) moveUp() {
	for c.Yb = -1 * CaseLen; c.Yb < 0; c.Yb++ {
		c.adaptSpeed()
	}
}

func (c *Character) moveDown() {
	for c.Yb = CaseLen; c.Yb > 0; c.Yb-- {
		c.adaptSpeed()
	}
}

func (c *Character) IsDead() bool {
	return c.Dead
}

func (c *Character) CalculateAttackScore() int {
	score := float64(c.Strength.Current) * (1.0 + float64(c.Luck.Current)/100.0)
	iscore := int(score)
	return iscore
}

func (c *Character) CalculatePowerAttackScore() int {
	score := float64(c.Will.Current) * (1.0 + float64(c.Luck.Current)/100.0)
	iscore := int(score)
	return iscore
}

func (p *Character) regenerate() {
	t := time.Now()
	deltaD := t.Sub(p.LastRegenerationTime)
	if deltaD > time.Duration(1000/p.RegenerationSpeed.Current)*time.Millisecond {
		if p.Energy.Current < p.Energy.Initial {
			p.Energy.Current += 5
		}
		if p.Health.Current < p.Health.Initial {
			p.Health.Current += 1
		}
		p.LastRegenerationTime = time.Now()
	}
}

func (c *Character) HasEnoughEnergy(cost int) bool {
	return c.Energy.Current > cost
}

func (c *Character) LooseEnergy(cost int) {
	c.Energy.Current -= cost
	if c.Energy.Current < 0 {
		c.Energy.Current = 0
	}
}

func (c *Character) TakeDamages(damage int) {
	if c.Dead {
		return
	}
	c.damagesTaken = damage
	for i := 0; i < 10; i++ {
		c.adaptSpeed()
	}
	c.damagesTaken = 0
	c.Health.Current -= damage
	if c.Health.Current <= 0 {
		c.Dead = true
		return
	}
}

func (c *Character) IsHurt() int {
	return c.damagesTaken
}

func (c *Character) IsAttacking() bool {
	return c.isAttacking
}

func (c *Character) SetAggressiveness(ag int) {
	c.Aggressiveness.Current = ag
}

func (c *Character) GetAggressiveness() int {
	return c.Aggressiveness.Current
}
