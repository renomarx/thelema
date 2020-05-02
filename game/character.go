package game

import (
	"log"
	"math/rand"
	"time"
)

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
	FightingPos    Pos
	Attacks        []Attack
	SelectedAttack Attack
	isAttacking    bool
	AttackPos      int
	damagesTaken   int
}

type PowerUser struct {
	Powers              map[string]*Power
	CurrentPower        *Power
	IsPowerUsing        bool
	PowerUsingStage     int
	ElementalAffinities map[MagickElement]int
	MagickLevel         map[MagickCategory]int
}

type Character struct {
	MovingObject
	Fighter
	PowerUser
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
	Defense              Characteristic
	Evasion              Characteristic
	Affinity             string
	ActionPoints         float64
	LastActionTime       time.Time
	LastRegenerationTime time.Time
	Dead                 bool
	VisionRange          int
	Shadow               bool
	Meditating           bool
}

func NewCharacter() Character {
	c := Character{}
	c.LastRegenerationTime = time.Now()
	c.Powers = make(map[string]*Power)
	c.ElementalAffinities = make(map[MagickElement]int)
	c.MagickLevel = make(map[MagickCategory]int)
	return c
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

func (c *Character) TakeDamages(damage int) bool {
	if c.Dead {
		return true
	}
	if c.Evasion.Current > 0 {
		r := rand.Intn(c.Evasion.Current)
		if r > damage {
			return false
		}
	}
	def := c.Defense.Current
	if def > 0 {
		d10 := damage / 10
		d90 := 9 * d10
		d90 -= d90 * def / 100
		if d90 < 0 {
			d90 = 0
		}
		damage = d10 + d90
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
			return true
		}
	}
	return true
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
	case "Defense":
		c.Defense.Raise(value)
	case "Evasion":
		c.Evasion.Raise(value)
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
	c.Defense.Reset()
	c.Evasion.Reset()
}

func (c *Character) GetElementalAffinity(element MagickElement) int {
	_, e := c.ElementalAffinities[element]
	if !e {
		c.ElementalAffinities[element] = 0
	}
	return c.ElementalAffinities[element]
}

func (c *Character) RaiseElementalAffinity(element MagickElement, x int) {
	c.GetElementalAffinity(element)
	c.ElementalAffinities[element] += x
}

func (c *Character) GetMagickLevel(cat MagickCategory) int {
	_, e := c.MagickLevel[cat]
	if !e {
		c.MagickLevel[cat] = 0
	}
	return c.MagickLevel[cat]
}

func (c *Character) RaiseMagickLevel(cat MagickCategory, x int) {
	c.GetMagickLevel(cat)
	c.MagickLevel[cat] += x
}

func (c *Character) GetAttacks() []Attack {
	var attacks []Attack
	for _, att := range c.Attacks {
		attacks = append(attacks, att)
	}
	for _, pow := range c.Powers {
		if pow.IsAttack {
			att := Attack{
				Strength:      pow.Strength,
				Name:          pow.Name,
				EnergyCost:    pow.Energy,
				Speed:         pow.Speed,
				Range:         pow.Range,
				Accuracy:      100,
				Type:          AttackTypeMagick,
				MagickUID:     pow.UID,
				MagickElement: pow.Element,
			}
			attacks = append(attacks, att)
		}
	}
	if len(attacks) == 0 {
		allAttacks := Attacks()
		attacks = append(attacks, allAttacks[0])
	}
	return attacks
}

func (c *Character) doAttack(ring *FightingRing, to []FighterInterface) {
	att := c.SelectedAttack
	c.isAttacking = true
	for c.AttackPos = 0; c.AttackPos < CaseLen; c.AttackPos++ {
		att.adaptSpeed()
	}
	c.isAttacking = false

	damages := att.GetPower(c)
	if damages == 0 {
		EM.Dispatch(&Event{
			Message: "L'attaque a échoué!",
		})
		return
	}
	switch att.Type {
	case AttackTypePhysical:
		for _, f := range to {
			f.TakeDamages(damages)
		}
		c.Strength.RaiseXp(damages * len(to) / 10)
		c.Dexterity.RaiseXp(1)
	case AttackTypeMagick:
		switch att.MagickUID {
		case PowerBrutalStrength:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerBrutalStrength}})
			ring.MakeEffect(c.GetFightingPos(), string(Healing), 400) // FIXME
			c.RaiseCharacteristic("Strength", damages)
		case PowerQuickening:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerQuickening}})
			ring.MakeEffect(c.GetFightingPos(), string(Healing), 400) // FIXME
			c.RaiseCharacteristic("Dexterity", damages)
		case PowerRockBody:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerRockBody}})
			ring.MakeEffect(c.GetFightingPos(), string(Healing), 400) // FIXME
			c.RaiseCharacteristic("Defense", damages)
		case PowerGlaciation:
			for _, f := range to {
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerRockBody}})
				ring.MakeEffect(f.GetFightingPos(), string(Ice), 600)
				f.TakeDamages(damages)
				f.LowerCharacteristic("Dexterity", damages/10)
			}
		case PowerStorm:
			for _, f := range to {
				ring.MakeStorm(f.GetFightingPos(), damages, Right, 200)
				f.TakeDamages(damages)
			}
		case PowerLightness:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerLightness}})
			ring.MakeEffect(c.GetFightingPos(), string(Healing), 400) // FIXME
			c.RaiseCharacteristic("Evasion", damages)
		case PowerHealing:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerHealing}})
			ring.MakeEffect(c.GetFightingPos(), string(Healing), 400)
			c.Health.Restore(damages)
		case PowerCalm:
			for _, f := range to {
				ring.MakeEffect(f.GetFightingPos(), string(Calm), 400)
				f.LowerCharacteristic("Aggressiveness", damages)
			}
		case PowerInvocation:
			monster := NewInvokedSpirit(damages)
			// TODO : case enemy
			monster.FightingPos = Pos{X: c.FightingPos.X, Y: c.FightingPos.Y + 1}
			ring.AddFriend(monster)
		case PowerFlames:
			for _, f := range to {
				ring.MakeFlame(f.GetFightingPos(), damages, 400)
				f.TakeDamages(damages)
			}
		default:
			log.Println("power default : ", att.MagickUID)
			for _, f := range to {
				f.TakeDamages(damages)
			}
		}
		c.LooseEnergy(att.EnergyCost)
		c.Intelligence.RaiseXp(1)
		targetsNumber := len(to)
		if targetsNumber == 0 {
			targetsNumber = 1
		}
		c.Will.RaiseXp(damages * targetsNumber / 10)
		c.Energy.RaiseXp(att.EnergyCost)
		c.RaiseElementalAffinity(att.MagickElement, 1)
		c.RaiseMagickLevel(att.MagickCategory, 1)
	}
}
