package game

import "math/rand"

type Monster struct {
	Character
	IsFriend bool
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

func (m *Monster) ChooseAction(ring *FightingRing) int {
	// TODO : monster IA
	return m.Speed.Current
}

func (m *Monster) Fight(ring *FightingRing) {
	m.isAttacking = true
	for m.AttackPos = 0; m.AttackPos < CaseLen; m.AttackPos++ {
		m.adaptSpeed()
	}
	m.isAttacking = false
	if m.IsFriend {
		e := ring.GetFirstEnemyNotDead()
		if e != nil {
			e.TakeDamages(m.CalculateAttackScore())
		}
	} else {
		ring.Player.TakeDamages(m.CalculateAttackScore())
	}
}
