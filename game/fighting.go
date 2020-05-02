package game

import (
	"math/rand"
	"time"
)

type FighterInterface interface {
	ChooseAction(ring *FightingRing) int
	Fight(ring *FightingRing)
	IsDead() bool
	IsCalmed() bool
	TakeDamages(damages int) bool
	GetTile() Tile
	GetHealth() Characteristic
	GetEnergy() Characteristic
	IsHurt() int
	IsAttacking() bool
	GetFightingPos() Pos
	SetFightingPos(p Pos)
	LowerCharacteristic(name string, value int)
	RaiseCharacteristic(name string, value int)
	ResetFightingSkills()
}

type FightingRing struct {
	IsOpen                    bool
	Running                   bool
	Round                     int
	SelectedPlayerAction      string
	Player                    FighterInterface
	Friends                   []FighterInterface
	Enemies                   []FighterInterface
	Stage                     FightingStage
	Menu                      *Menu
	AttacksMenuOpen           bool
	AttackTargetSelectionOpen bool
	TargetSelected            int
	PossibleAttacks           struct {
		List     []Attack
		Selected int
	}
	roundFighters []RoundFighter
	CurrentEffect *Effect
}

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
		Message: "Vous êtes attaqué!",
	})
	p := g.Level.Player
	g.Level.MakeExplosion(p.Pos, 88, 1000)
	g.FightingRing = NewFightingRing()
	y := 0
	for _, e := range enemies {
		e.SetFightingPos(Pos{X: 1, Y: y})
		g.FightingRing.AddEnemy(e)
		y++
	}
	g.FightingRing.Player = p
	g.FightingRing.Player.SetFightingPos(Pos{X: 0, Y: 0})
	if p.Friend != nil && !p.Friend.IsDead() {
		p.Friend.SetFightingPos(Pos{X: 0, Y: 1})
		g.FightingRing.AddFriend(p.Friend)
	}
	g.FightingRing.Start()
	for g.FightingRing.Running {
		g.FightingRing.PlayRound(g)
	}
	EM.Dispatch(&Event{
		Message: "Vous avez vaincu tous les ennemis!",
	})
	time.Sleep(2 * time.Second)
	g.FightingRing.Close()
	g.FightingRing = nil
	EM.Dispatch(&Event{
		Action:  ActionStopFight,
		Payload: map[string]string{"levelName": g.Level.Name},
	})
}

func (ring *FightingRing) Close() {
	ring.IsOpen = false
	ring.Player.ResetFightingSkills()
	for _, e := range ring.Enemies {
		e.ResetFightingSkills()
	}
	for _, f := range ring.Friends {
		f.ResetFightingSkills()
	}
}

func NewFightingRing() *FightingRing {
	fr := &FightingRing{
		IsOpen: false,
		Round:  0,
		Stage:  FightingChoice,
	}
	fr.LoadFightingMenu()
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
	ring.Running = true
}

func (ring *FightingRing) End() {
	ring.Running = false
}

func (ring *FightingRing) IsFinished() bool {
	if ring.Player.IsDead() {
		return true
	}
	for _, e := range ring.Enemies {
		if !e.IsDead() && !e.IsCalmed() {
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
	ring.OpenFightingMenu()
	for ring.Menu.IsOpen {
		ring.HandleInputFightingMenu(g.input)
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
	i := 0
	for i < len(fr.Enemies) && fr.Enemies[i].IsDead() {
		i++
	}
	if i == len(fr.Enemies) {
		fr.End()
		return
	}
	fr.TargetSelected = i
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
	attacks := p.GetAttacks()
	fr.PossibleAttacks.List = attacks
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

func (fr *FightingRing) GetSelectedAttack() Attack {
	return fr.PossibleAttacks.List[fr.PossibleAttacks.Selected]
}

func (fr *FightingRing) NextTarget() {
	i := fr.TargetSelected + 1
	for i < len(fr.Enemies) && fr.Enemies[i].IsDead() {
		i++
	}
	if i >= len(fr.Enemies) {
		return
	}
	fr.TargetSelected = i
}

func (fr *FightingRing) LastTarget() {
	i := fr.TargetSelected - 1
	for i >= 0 && fr.Enemies[i].IsDead() {
		i--
	}
	if i < 0 {
		return
	}
	fr.TargetSelected = i
}

func (a *Attack) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/a.Speed) * time.Millisecond)
}
