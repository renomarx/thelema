package game

type Enemy struct {
	Character
}

func (level *Level) MakeEnemy(npc *Npc) *Enemy {
	e := &Enemy{}
	e.Character = npc.Character
	e.Aggressiveness.Init(20)
	e.Speed.Init(10)
	return e
}

func (m *Enemy) ChooseAction(ring *FightingRing) int {
	// TODO : enemy IA
	attacks := m.GetAttacks()
	m.SelectedAttack = attacks[0]
	return m.SelectedAttack.GetSpeed(&m.Character)
}

func (m *Enemy) Fight(ring *FightingRing) {
	if m.IsCalmed() {
		EM.Dispatch(&Event{
			Message: m.Name + " est calmé, n'attaquera pas.",
		})
		return
	}

	var to []FighterInterface
	idx := 0
	att := m.SelectedAttack
	// They first try to attack player friends
	for i := 0; idx < att.Range && i < len(ring.Friends); i++ {
		f := ring.Friends[i]
		if !f.IsDead() {
			to = append(to, f)
			idx++
		}
	}
	if idx < att.Range {
		// If there is still a place in attack range, attack player
		to = append(to, ring.Player)
	}

	m.doAttack(ring, to)
}
