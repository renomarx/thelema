package game

import "time"

const CharacterDeltaTime = 100

type Characteristic struct {
	Initial int
	Current int
}

func (ch *Characteristic) Init(value int) {
	ch.Initial = value
	ch.Current = value
}

type Fighter struct {
	IsAttacking      bool
	IsPowerAttacking bool
}

type Talker struct {
	IsTalking bool
	Dialog    *Dialog
}

type Character struct {
	MovingObject
	Fighter
	LookAt            InputType
	Name              string
	Hitpoints         Characteristic
	Energy            Characteristic
	Strength          Characteristic
	Speed             Characteristic
	RegenerationSpeed Characteristic
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
