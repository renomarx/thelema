package game

type Invocation struct {
	Character
}

func NewInvokedSpirit(strength int) *Invocation {
	monster := &Invocation{}
	monster.Character = NewCharacter()
	monster.Rune = string(Spirit)
	monster.Name = "Invoked Spirit"
	monster.Health.Init(strength * 10)
	monster.Energy.Init(strength * 10)
	monster.Strength.Init(strength)
	monster.Dexterity.Init(strength)
	monster.Will.Init(strength)
	monster.Intelligence.Init(strength)
	monster.Luck.Init(20)
	monster.Beauty.Init(0)
	monster.Speed.Init(10)
	monster.Aggressiveness.Init(strength * 10)
	return monster
}

func (m *Invocation) ChooseAction(ring *FightingRing) int {
	// TODO : invocation IA
	attacks := m.GetAttacks()
	m.SelectedAttack = attacks[0]
	return m.SelectedAttack.GetSpeed(&m.Character)
}

func (m *Invocation) Fight(ring *FightingRing) {
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
