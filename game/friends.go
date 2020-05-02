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
	attacks := m.GetAttacks()
	m.SelectedAttack = attacks[0]
	return m.SelectedAttack.GetSpeed(&m.Character)
}

func (m *Friend) Fight(ring *FightingRing) {
	if m.IsCalmed() {
		EM.Dispatch(&Event{
			Message: m.Name + " est calmé, n'attaquera pas.",
		})
		return
	}

	var to []FighterInterface
	idx := 0
	att := m.SelectedAttack
	for i := 0; idx < att.Range && i < len(ring.Enemies); i++ {
		f := ring.Enemies[i]
		if !f.IsDead() {
			to = append(to, f)
			idx++
		}
	}

	m.doAttack(ring, to)
}
