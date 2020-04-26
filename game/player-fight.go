package game

import (
	"log"
	"math/rand"
)

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

func (p *Player) MeetMonsters(g *Game) {
	if p.Shadow {
		return
	}
	l := g.Level
	r := rand.Intn(100000) % 100
	cc := l.Map[p.Z][p.Y][p.X]
	if r >= cc.MonstersProbability {
		return
	}
	switch cc.T {
	// TODO : other floor types
	case MonsterFloor:
		bestiary := Bestiary()
		g.FightMonsters(bestiary)
	}
}

func (p *Player) ChooseAction(ring *FightingRing) int {
	switch ring.SelectedPlayerAction {
	case "run":
		return p.Dexterity.Current
	case "attack":
		att := ring.PossibleAttacks.List[ring.PossibleAttacks.Selected]
		p.currentAttack = att
		return att.GetSpeed(&p.Character)
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

		damages := att.GetPower(&p.Character)
		switch att.Type {
		case AttackTypePhysical:
			for _, f := range to {
				f.TakeDamages(damages)
			}
			p.Strength.RaiseXp(damages * len(to) / 10)
			p.Dexterity.RaiseXp(1)
		case AttackTypeMagick:
			switch att.MagickType {
			case PowerBrutalStrength:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerBrutalStrength}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Healing), 400) // FIXME
				p.RaiseCharacteristic("Strength", damages)
			case PowerQuickening:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerQuickening}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Healing), 400) // FIXME
				p.RaiseCharacteristic("Dexterity", damages)
			case PowerRockBody:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerRockBody}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Healing), 400) // FIXME
				p.RaiseCharacteristic("Defense", damages)
			case PowerHealing:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerHealing}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Healing), 400)
				p.Health.Restore(damages)
			case PowerInvocation:
				monster := NewInvokedSpirit()
				ring.AddFriend(monster)
				p.Will.RaiseXp(monster.Strength.Initial / 10)
			case PowerFlames:
				for i, f := range to {
					y := ring.TargetSelected + i
					ring.MakeFlame(Pos{X: 1, Y: y}, damages, 400)
					f.TakeDamages(damages)
				}
			case PowerStorm:
				for i, f := range to {
					y := ring.TargetSelected + i
					ring.MakeStorm(Pos{X: 1, Y: y}, damages, Right, 200)
					f.TakeDamages(damages)
				}
			case PowerLightness:
				EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerLightness}})
				ring.MakeEffect(Pos{X: 0, Y: 0}, string(Healing), 400) // FIXME
				p.RaiseCharacteristic("Evasion", damages)
			case PowerCalm:
				for i, f := range to {
					y := ring.TargetSelected + i
					ring.MakeEffect(Pos{X: 1, Y: y}, string(Calm), 400)
					f.LowerCharacteristic("Aggressiveness", damages)
				}
			default:
				log.Println("power default : ", att.MagickType)
				for _, f := range to {
					f.TakeDamages(damages)
				}
			}
			p.LooseEnergy(att.EnergyCost)
			p.Intelligence.RaiseXp(1)
			targetsNumber := len(to)
			if targetsNumber == 0 {
				targetsNumber = 1
			}
			p.Will.RaiseXp(damages * targetsNumber / 10)
			p.Energy.RaiseXp(att.EnergyCost)
		}
	}
}
