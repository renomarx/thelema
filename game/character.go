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

func (ch *Characteristic) Restore(value int) {
	ch.Current += value
	if ch.Current > ch.Initial {
		ch.Current = ch.Initial
	}
}

func (ch *Characteristic) RaisePermanently(value int) {
	ch.Initial += value
	ch.Current = ch.Initial
}

func (ch *Characteristic) RaiseXp(value int) {
	ch.Xp += value
	if ch.Xp >= ch.Initial*CharacteristicXpMultiplier {
		ch.Initial += ch.Initial / 20
		ch.Current = ch.Initial
		ch.Xp = 0
	}
}

func (ch *Characteristic) Raise(value int) {
	ch.Current += value
}

func (ch *Characteristic) Lower(value int) {
	tmp := ch.Current - value
	if tmp < 0 {
		tmp = 0
	}
	ch.Current = tmp
}

func (ch *Characteristic) Reset() {
	ch.Current = ch.Initial
}

const VoiceMaleStandard = "MALE_STANDARD"
const VoiceFemaleStandard = "FEMALE_STANDARD"

type Talker struct {
	Dialog   *Dialog
	Voice    string
	Talkable bool
}

type Fighter struct {
	FightingPos  Pos
	Attacks      []*Attack
	isAttacking  bool
	AttackPos    int
	PowerPos     int
	damagesTaken int
}

type Character struct {
	MovingObject
	Fighter
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
	LastRegenerationTime time.Time
	IsPowerUsing         bool
	Shadow               bool
	Meditating           bool
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
	defer func() {
		c.damagesTaken = 0
	}()
	for i := 0; i < damage; i++ {
		c.adaptSpeed()
		c.Health.Current--
		if c.Health.Current <= 0 {
			c.Dead = true
			return
		}
	}
}

func (c *Character) IsHurt() int {
	return c.damagesTaken
}

func (c *Character) IsAttacking() bool {
	return c.isAttacking
}

func (c *Character) IsCalmed() bool {
	return c.Aggressiveness.Current <= 0
}

func (c *Character) GetFightingPos() Pos {
	return c.FightingPos
}

func (c *Character) SetFightingPos(p Pos) {
	c.FightingPos = p
}

func (c *Character) LowerCharacteristic(name string, value int) {
	switch name {
	case "Health":
		c.Health.Lower(value)
	case "Energy":
		c.Energy.Lower(value)
	case "Strength":
		c.Strength.Lower(value)
	case "Dexterity":
		c.Dexterity.Lower(value)
	case "Beauty":
		c.Beauty.Lower(value)
	case "Will":
		c.Will.Lower(value)
	case "Intelligence":
		c.Intelligence.Lower(value)
	case "Charisma":
		c.Charisma.Lower(value)
	case "Luck":
		c.Luck.Lower(value)
	case "Aggressiveness":
		c.Aggressiveness.Lower(value)
	}
}

func (c *Character) RaiseCharacteristic(name string, value int) {
	switch name {
	case "Strength":
		c.Strength.Raise(value)
	case "Dexterity":
		c.Dexterity.Raise(value)
	case "Beauty":
		c.Beauty.Raise(value)
	case "Will":
		c.Will.Raise(value)
	case "Intelligence":
		c.Intelligence.Raise(value)
	case "Charisma":
		c.Charisma.Raise(value)
	case "Luck":
		c.Luck.Raise(value)
	case "Aggressiveness":
		c.Aggressiveness.Raise(value)
	}
}

func (c *Character) ResetFightingSkills() {
	c.Strength.Reset()
	c.Dexterity.Reset()
	c.Beauty.Reset()
	c.Will.Reset()
	c.Intelligence.Reset()
	c.Charisma.Reset()
	c.Luck.Reset()
	c.Aggressiveness.Reset()
}
