package game

type Enemy struct {
	Character
}

func (level *Level) MakeEnemy(pnj *Pnj) *Enemy {
	e := &Enemy{}
	e.Character = pnj.Character
	e.Aggressiveness.Init(20)
	e.Speed.Init(10)
	return e
}

func (m *Enemy) ChooseAction(ring *FightingRing) int {
	// TODO : enemy IA
	return m.Dexterity.Current
}

func (m *Enemy) Fight(ring *FightingRing) {
	if m.GetAggressiveness() <= 0 {
		EM.Dispatch(&Event{
			Message: m.Name + " is calmed, not attacking.",
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
