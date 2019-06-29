package game

import "math/rand"
import "fmt"

type FighterInterface interface {
	Fight(ring *FightingRing) AttackInterface
	IsDead() bool
	TakeDamages(damages int)
	GetTile() Tile
	GetHealth() Characteristic
	GetEnergy() Characteristic
}

type FightingRing struct {
	IsOpen               bool
	Round                int
	SelectedPlayerAction string
	Player               FighterInterface
	Friends              []FighterInterface
	Enemies              []FighterInterface
	Stage                FightingStage
	Attacks              []AttackInterface
	AttacksMenuOpen      bool
	PossibleAttacks      struct {
		List     []AttackInterface
		Selected int
	}
}

type FightingStage string

const FightingChoice FightingStage = "CHOICE"
const FightingAttacks FightingStage = "ATTACKS"

func (g *Game) FightMonsters(bestiary []*MonsterType) {
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionFight,
		Message: "You're being attacked!",
	})
	g.FightingRing = NewFightingRing()
	p := g.Level.Player
	g.FightingRing.LoadPossibleAttacks(p)
	nb := rand.Intn(2) + 1
	for i := 0; i < nb; i++ {
		m := rand.Intn(len(bestiary))
		proba := rand.Intn(100)
		mt := bestiary[m]
		for proba > mt.Probability {
			m := rand.Intn(len(bestiary))
			proba = rand.Intn(100)
			mt = bestiary[m]
		}
		mo := NewMonster(mt)
		g.FightingRing.AddEnemy(mo)
		g.FightingRing.Player = p
	}
	g.FightingRing.Start()
	for g.FightingRing.IsOpen {
		g.FightingRing.PlayRound(g)
	}
	g.FightingRing = nil
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionStopFight,
		Payload: map[string]string{"levelType": g.Level.Type},
	})
}

func NewFightingRing() *FightingRing {
	fr := &FightingRing{
		IsOpen: false,
		Round:  0,
		Stage:  FightingChoice,
	}
	return fr
}

func (ring *FightingRing) AddEnemy(e FighterInterface) {
	ring.Enemies = append(ring.Enemies, e)
}

func (ring *FightingRing) Start() {
	ring.IsOpen = true
}

func (ring *FightingRing) End() {
	ring.IsOpen = false
}

func (ring *FightingRing) IsFinished() bool {
	if ring.Player.IsDead() {
		return true
	}
	for _, e := range ring.Enemies {
		if !e.IsDead() {
			return false
		}
	}

	return true
}

func (ring *FightingRing) PlayRound(g *Game) {
	if !ring.IsOpen {
		return
	}
	if ring.IsFinished() {
		ring.End()
		return
	}

	ring.Stage = FightingChoice
	g.OpenFightingMenu()
	for g.FightingMenu.IsOpen {
		g.HandleInputFightingMenu()
	}
	a := ring.Player.Fight(ring)
	if a == nil {
		ring.End()
		return
	}
	ring.prepareAttack(a)
	for _, e := range ring.Enemies {
		if !e.IsDead() {
			a = e.Fight(ring)
			ring.prepareAttack(a)
		}
	}
	for _, f := range ring.Friends {
		if !f.IsDead() {
			a = f.Fight(ring)
			ring.prepareAttack(a)
		}
	}

	ring.Stage = FightingAttacks
	fmt.Println(ring.Round, ring.Attacks)
	for _, at := range ring.Attacks {
		at.Play(ring)
		g.GetEventManager().Dispatch(&Event{
			Action: ActionAttack,
		})
	}
	ring.Attacks = nil
	ring.Round++
}

func (ring *FightingRing) prepareAttack(a AttackInterface) {
	if a == nil {
		return
	}
	//ring.Attacks = append(ring.Attacks, a)
	pos := len(ring.Attacks)
	for i, at := range ring.Attacks {
		if a.GetSpeed() < at.GetSpeed() {
			pos = i
			break
		}
	}
	ring.Attacks = append(ring.Attacks[:pos], append([]AttackInterface{a}, ring.Attacks[pos:]...)...)
}

func (fr *FightingRing) LoadPossibleAttacks(p *Player) {
	fr.PossibleAttacks.List = append(fr.PossibleAttacks.List, &SwordAttack{
		Speed:   p.Weapon.Speed,
		Damages: p.CalculateAttackScore(),
	})
	fr.PossibleAttacks.List = append(fr.PossibleAttacks.List, &BiteAttack{
		Damages: 10,
	})
	fr.PossibleAttacks.Selected = 0
}

func (fr *FightingRing) NextPossibleAttack() {
	i := fr.PossibleAttacks.Selected + 1
	if i >= len(fr.PossibleAttacks.List) {
		i = 0
	}
	fr.PossibleAttacks.Selected = i
}

func (fr *FightingRing) LastPossibleAttack() {
	i := fr.PossibleAttacks.Selected - 1
	if i <= 0 {
		i = len(fr.PossibleAttacks.List) - 1
	}
	fr.PossibleAttacks.Selected = i
}
