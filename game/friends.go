package game

type Friend struct {
	Character
}

func (level *Level) MakeFriend(pnj *Pnj) *Friend {
	f := &Friend{}
	f.Character = pnj.Character
	level.Map[pnj.Z][pnj.Y][pnj.X].Pnj = nil
	return f
}

func (m *Friend) ChooseAction(ring *FightingRing) int {
	// TODO : friend IA
	return m.Dexterity.Current
}

func (m *Friend) Fight(ring *FightingRing) {
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
