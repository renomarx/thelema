package game

func (p *Player) TakeDamages(damage int) bool {
	done := p.Character.TakeDamages(damage)
	if !done {
		EM.Dispatch(&Event{
			Message: "Vous avez esquivé l'attaque!",
		})
		return false
	}
	p.Health.RaiseXp(damage)
	return true
}

func (p *Player) ChooseAction(ring *FightingRing) int {
	switch ring.SelectedPlayerAction {
	case FightingActionRun:
		return p.Dexterity.Current
	case FightingActionAttack:
		att := ring.PossibleAttacks.List[ring.PossibleAttacks.Selected]
		p.SelectedAttack = att
		return att.GetSpeed(&p.Character)
	}
	return 0
}

func (p *Player) Fight(ring *FightingRing) {
	switch ring.SelectedPlayerAction {
	case FightingActionRun:
		ring.End()
	case FightingActionAttack:
		att := p.SelectedAttack
		var to []FighterInterface
		idx := 0
		for i := ring.TargetSelected; idx < att.Range && i < len(ring.Enemies); i++ {
			f := ring.Enemies[i]
			if !f.IsDead() {
				to = append(to, f)
				idx++
			}
		}
		p.doAttack(ring, to)
	}
}
