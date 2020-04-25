package game

type Enemy struct {
	Character
	SelectedAttack *Attack
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
	m.SelectedAttack = &Attack{
		Speed:   20,
		Damages: 20,
		Name:    "Sword attack",
		Type:    AttackTypePhysical,
		Range:   1,
	}
	for _, pow := range m.Powers {
		att := &Attack{
			Damages:    pow.Strength,
			Name:       pow.Name,
			EnergyCost: pow.Energy,
			Speed:      pow.Speed,
			Range:      pow.Range,
			Type:       AttackTypeMagick,
			MagickType: pow.UID,
		}
		m.SelectedAttack = att
		return att.Speed
	}
	return m.SelectedAttack.GetSpeed(&m.Character)
}

func (m *Enemy) Fight(ring *FightingRing) {
	if m.GetAggressiveness() <= 0 {
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

	damages := m.CalculateAttackScore()
	if m.SelectedAttack != nil {
		att := m.SelectedAttack
		damages = att.GetPower(&m.Character)
		switch att.Type {
		case AttackTypePhysical:
		case AttackTypeMagick:
			switch att.MagickType {
			case PowerHealing:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerHealing}})
				ring.MakeEffect(Pos{X: 1, Y: 0}, string(Healing), 400)
				m.Health.Add(damages)
			// case PowerInvocation:
			// 	monster := NewInvokedSpirit()
			// 	ring.AddFriend(monster)
			// 	p.Will.RaiseXp(monster.Strength.Initial / 10)
			case PowerFlames:
				ring.MakeFlame(Pos{X: 0, Y: 0}, damages, 400)
			case PowerStorm:
				ring.MakeStorm(Pos{X: 0, Y: 0}, damages, Right, 200)
			case PowerCalm:
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Calm), 400)
			}
			m.LooseEnergy(att.EnergyCost)
		}

	}
	if len(ring.Friends) > 0 {
		for _, f := range ring.Friends {
			if !f.IsDead() {
				f.TakeDamages(damages)
				return
			}
		}
	}
	ring.Player.TakeDamages(damages)
}
