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
