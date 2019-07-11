package game

type Friend struct {
	Character
}

func (level *Level) MakeFriend(pnj *Pnj) *Friend {
	f := &Friend{}
	f.Character = pnj.Character
	f.Speed.Init(f.Speed.Current * 2)
	level.Map[pnj.Y][pnj.X].Pnj = nil
	return f
}

func (m *Friend) ChooseAction(ring *FightingRing) int {
	// TODO : friend IA
	return m.Speed.Current
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
