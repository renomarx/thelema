package game

import "math/rand"
import "time"

type FighterInterface interface {
	ChooseAction(ring *FightingRing) int
	Fight(ring *FightingRing)
	IsDead() bool
	TakeDamages(damages int)
	GetTile() Tile
	GetHealth() Characteristic
	GetEnergy() Characteristic
	IsHurt() int
	IsAttacking() bool
}

type FightingRing struct {
	IsOpen                    bool
	Round                     int
	SelectedPlayerAction      string
	Player                    FighterInterface
	Friends                   []FighterInterface
	Enemies                   []FighterInterface
	Stage                     FightingStage
	AttacksMenuOpen           bool
	AttackTargetSelectionOpen bool
	TargetSelected            int
	PossibleAttacks           struct {
		List     []*Attack
		Selected int
	}
	roundFighters []RoundFighter
}

type Attack struct {
	Damages    int
	Speed      int
	EnergyCost int
	Name       string
	Range      int
	Type       AttackType
	MagickType PowerType
}
type AttackType string

const AttackTypePhysical AttackType = "PHYSICAL"
const AttackTypeMagick AttackType = "MAGICK"

type RoundFighter struct {
	speed int
	f     FighterInterface
}

type FightingStage string

const FightingChoice FightingStage = "CHOICE"
const FightingAttacks FightingStage = "ATTACKS"

func (g *Game) FightMonsters(bestiary []*MonsterType) {
	var enemies []FighterInterface
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
		enemies = append(enemies, mo)
	}
	g.Fight(enemies)
}

func (g *Game) Fight(enemies []FighterInterface) {
	EM.Dispatch(&Event{
		Action:  ActionFight,
		Message: "You're being attacked!",
	})
	g.FightingRing = NewFightingRing()
	p := g.Level.Player
	for _, e := range enemies {
		g.FightingRing.AddEnemy(e)
	}
	g.FightingRing.Player = p
	if p.Friend != nil && !p.Friend.IsDead() {
		g.FightingRing.AddFriend(p.Friend)
	}
	g.FightingRing.Start()
	for g.FightingRing.IsOpen {
		g.FightingRing.PlayRound(g)
	}
	g.FightingRing = nil
	EM.Dispatch(&Event{
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

func (ring *FightingRing) AddFriend(f FighterInterface) {
	ring.Friends = append(ring.Friends, f)
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
	p := g.Level.Player
	ring.LoadPossibleAttacks(p)
	g.OpenFightingMenu()
	for g.FightingMenu.IsOpen {
		g.HandleInputFightingMenu()
	}
	speed := ring.Player.ChooseAction(ring)
	ring.prepareRoundFighter(ring.Player, speed)
	for _, e := range ring.Friends {
		if !e.IsDead() {
			speed = e.ChooseAction(ring)
			ring.prepareRoundFighter(e, speed)
		}
	}
	for _, e := range ring.Enemies {
		if !e.IsDead() {
			speed = e.ChooseAction(ring)
			ring.prepareRoundFighter(e, speed)
		}
	}

	ring.Stage = FightingAttacks
	for _, rf := range ring.roundFighters {
		if !rf.f.IsDead() {
			rf.f.Fight(ring)
			EM.Dispatch(&Event{
				Action: ActionAttack,
			})
		}
		if !ring.IsOpen {
			return
		}
	}
	ring.clearRound()
	ring.Round++
}

func (fr *FightingRing) clearRound() {
	fr.roundFighters = nil
	var enemies []FighterInterface
	for _, e := range fr.Enemies {
		if !e.IsDead() {
			enemies = append(enemies, e)
		}
	}
	var friends []FighterInterface
	for _, f := range fr.Friends {
		if !f.IsDead() {
			friends = append(friends, f)
		}
	}
	fr.Friends = friends
	fr.TargetSelected = 0
}

func (ring *FightingRing) prepareRoundFighter(f FighterInterface, speed int) {
	if f == nil {
		return
	}
	rf := RoundFighter{
		speed: speed,
		f:     f,
	}
	//ring.roundFighters = append(ring.roundFighters, f)
	pos := len(ring.roundFighters)
	for i, rff := range ring.roundFighters {
		if rf.speed > rff.speed {
			pos = i
			break
		}
	}
	ring.roundFighters = append(ring.roundFighters[:pos], append([]RoundFighter{rf}, ring.roundFighters[pos:]...)...)
}

func (fr *FightingRing) LoadPossibleAttacks(p *Player) {
	att := &Attack{
		Speed:   p.Dexterity.Current,
		Damages: 10,
		Name:    "Sword attack",
		Type:    AttackTypePhysical,
	}
	att.Range = 1
	fr.PossibleAttacks.List = append(fr.PossibleAttacks.List, att)
	for _, pow := range p.Powers {
		att := &Attack{
			Damages:    pow.Strength,
			Name:       pow.Name,
			EnergyCost: pow.Energy,
			Speed:      pow.Speed,
			Range:      pow.Range,
			Type:       AttackTypeMagick,
			MagickType: pow.Type,
		}
		fr.PossibleAttacks.List = append(fr.PossibleAttacks.List, att)
	}
	fr.PossibleAttacks.Selected = 0
}

func (fr *FightingRing) NextPossibleAttack() {
	i := fr.PossibleAttacks.Selected + 1
	if i >= len(fr.PossibleAttacks.List) {
		i = len(fr.PossibleAttacks.List) - 1
	}
	fr.PossibleAttacks.Selected = i
}

func (fr *FightingRing) LastPossibleAttack() {
	i := fr.PossibleAttacks.Selected - 1
	if i <= 0 {
		i = 0
	}
	fr.PossibleAttacks.Selected = i
}

func (fr *FightingRing) GetSelectedAttack() *Attack {
	return fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
}

func (fr *FightingRing) NextTarget() {
	i := fr.TargetSelected + 1
	if i >= len(fr.Enemies) {
		i = len(fr.Enemies) - 1
	}
	fr.TargetSelected = i
}

func (fr *FightingRing) LastTarget() {
	i := fr.TargetSelected - 1
	if i < 0 {
		i = 0
	}
	fr.TargetSelected = i
}

func (fr *FightingRing) GetFirstEnemyNotDead() FighterInterface {
	for _, e := range fr.Enemies {
		if !e.IsDead() {
			return e
		}
	}
	return nil
}

func (a *Attack) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/a.Speed) * time.Millisecond)
}
