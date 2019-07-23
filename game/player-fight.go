package game

import "math/rand"

func (p *Player) TakeDamages(damage int) {
	p.Character.TakeDamages(damage)
	p.Health.RaiseXp(damage)
}

func (p *Player) MeetMonsters(g *Game) {
	l := g.Level
	r := rand.Intn(100000) % 100
	cc := l.Map[p.Y][p.X]
	if r >= cc.MonstersProbability {
		return
	}
	switch cc.T {
	// TODO : other floor types
	case HerbFloor:
		bestiary := Bestiary()
		g.FightMonsters(bestiary)
	case DirtFloor:
		bestiary := BestiaryUnderworld()
		g.FightMonsters(bestiary)
	}
}

func (p *Player) ChooseAction(ring *FightingRing) int {
	switch ring.SelectedPlayerAction {
	case "run":
		return p.Speed.Current
	case "attack":
		att := ring.PossibleAttacks.List[ring.PossibleAttacks.Selected]
		p.currentAttack = att
		return att.Speed
	}
	return 0
}

func (p *Player) Fight(ring *FightingRing) {
	switch ring.SelectedPlayerAction {
	case "run":
		ring.End()
	case "attack":
		att := p.currentAttack
		var to []FighterInterface
		idx := 0
		for i := ring.TargetSelected; idx < att.Range && i < len(ring.Enemies); i++ {
			f := ring.Enemies[i]
			if !f.IsDead() {
				to = append(to, f)
				idx++
			}
		}

		p.isAttacking = true
		for p.AttackPos = 0; p.AttackPos < CaseLen; p.AttackPos++ {
			att.adaptSpeed()
		}
		p.isAttacking = false

		switch att.Type {
		case AttackTypePhysical:
			damages := att.Damages * p.CalculateAttackScore() / 10
			for _, f := range to {
				f.TakeDamages(damages)
			}
			p.Strength.RaiseXp(damages * len(to) / 10)
			p.Dexterity.RaiseXp(1)
		case AttackTypeMagick:
			switch att.MagickType {
			case PowerHealing:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerHealing}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, rune(Healing), 400)
				p.Health.Add(att.Damages * p.CalculatePowerAttackScore() / 10)
			case PowerInvocation:
				monster := NewInvokedSpirit()
				ring.AddFriend(monster)
				p.Energy.RaiseXp(att.EnergyCost)
				p.Will.RaiseXp(p.Will.Initial / 10)
			case PowerFlames:
				damages := att.Damages * p.CalculatePowerAttackScore() / 10
				for i, f := range to {
					y := ring.TargetSelected + i
					ring.MakeFlame(Pos{X: 1, Y: y}, damages, 400)
					f.TakeDamages(damages)
				}
				p.Energy.RaiseXp(damages)
				p.Will.RaiseXp(damages * len(to) / 10)
			case PowerStorm:
				damages := att.Damages * p.CalculatePowerAttackScore() / 10
				for i, f := range to {
					y := ring.TargetSelected + i
					ring.MakeStorm(Pos{X: 1, Y: y}, damages, Right, 200)
					f.TakeDamages(damages)
				}
				p.Energy.RaiseXp(damages)
				p.Will.RaiseXp(damages * len(to) / 10)
			default:
				damages := att.Damages * p.CalculatePowerAttackScore() / 10
				for _, f := range to {
					f.TakeDamages(damages)
				}
				p.Energy.RaiseXp(damages)
				p.Will.RaiseXp(damages * len(to) / 10)
			}
			p.LooseEnergy(att.EnergyCost)
			p.Intelligence.RaiseXp(1)
		}
	}
}
