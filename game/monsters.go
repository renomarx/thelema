package game

import "math/rand"

type Monster struct {
	Character
}

func NewMonster(mt *MonsterType) *Monster {
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
	return monster
}

func (m *Monster) Fight(ring *FightingRing) AttackInterface {
	// TODO : monster IA
	bt := &BiteAttack{
		Damages: m.CalculateAttackScore(),
	}
	bt.From = m
	bt.To = append(bt.To, ring.Player)
	bt.Speed = 10
	return bt
}

func (m *Monster) TakeDamages(damage int) {
	if m.isDead {
		return
	}
	m.Health.Current -= damage
	m.Health.RaiseXp(damage)
	if m.Health.Current <= 0 {
		m.isDead = true
		return
	}
}
