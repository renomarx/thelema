package game

type Invocation struct {
	Character
}

func NewInvokedSpirit() *Invocation {
	monster := &Invocation{}
	monster.Rune = string(Spirit)
	monster.Name = "Invoked Spirit"
	monster.Health.Init(200)
	monster.Energy.Init(200)
	monster.Strength.Init(30)
	monster.Dexterity.Init(30)
	monster.Will.Init(20)
	monster.Intelligence.Init(20)
	monster.Luck.Init(20)
	monster.Beauty.Init(0)
	return monster
}

func (m *Invocation) ChooseAction(ring *FightingRing) int {
	// TODO : invocation IA
	return m.Dexterity.Current
}

func (m *Invocation) Fight(ring *FightingRing) {
	m.isAttacking = true
	for m.AttackPos = 0; m.AttackPos < CaseLen; m.AttackPos++ {
		m.adaptSpeed()
	}
	m.isAttacking = false
	e := ring.GetFirstEnemyNotDead()
	if e != nil {
		e.TakeDamages(m.CalculateAttackScore())
	}
}
