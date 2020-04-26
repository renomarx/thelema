package game

import "math/rand"

type Monster struct {
	Character
}

func NewMonster(mt *MonsterType) *Monster {
	monster := &Monster{}
	monster.Speed.Init(10)
	monster.Rune = string(mt.Tile)
	monster.Name = mt.Name
	monster.Health.Init(mt.Health)
	monster.Energy.Init(mt.Energy)
	monster.Strength.Init(mt.Strength)
	monster.Dexterity.Init(mt.Speed)
	monster.Will.Init(mt.Strength)
	monster.Intelligence.Init(mt.Speed)
	monster.Luck.Init(mt.Luck)
	monster.Beauty.Init(rand.Intn(20))
	monster.Aggressiveness.Init(mt.Aggressiveness)
	return monster
}

func (m *Monster) ChooseAction(ring *FightingRing) int {
	// TODO : monster IA
	return m.Dexterity.Current * 2
}

func (m *Monster) Fight(ring *FightingRing) {
	if m.IsCalmed() {
		EM.Dispatch(&Event{
			Message: m.Name + " est calmé, n'attaquera pas.",
		})
		return
	}
	m.isAttacking = true
	for m.AttackPos = 0; m.AttackPos < CaseLen; m.AttackPos++ {
		m.adaptSpeed()
	}
	m.isAttacking = false
	if len(ring.Friends) > 0 {
		for _, f := range ring.Friends {
			if !f.IsDead() {
				f.TakeDamages(m.CalculateAttackScore())
				return
			}
		}
	}
	ring.Player.TakeDamages(m.CalculateAttackScore())
}
