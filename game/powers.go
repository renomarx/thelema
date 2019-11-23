package game

import "sort"
import "time"

type PowerType string

const PowerInvocation = "invocation"
const PowerFlames = "flames"
const PowerStorm = "storm"
const PowerHealing = "healing"
const PowerDeadSpeaking = "dead_speaking"
const PowerCalm = "calm"

type PlayerPower struct {
	Type        PowerType
	Tile        Tile
	Strength    int
	Speed       int
	Energy      int
	Range       int
	Description string
	Name        string
	IsAttack    bool
}

func (pow *PlayerPower) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/pow.Speed) * time.Millisecond)
}

func Powers() map[string]*PlayerPower {
	return map[string]*PlayerPower{
		string(PowerInvocation): &PlayerPower{Type: PowerInvocation, Speed: 5, Energy: 50, Tile: Spirit,
			Description: "Summon an astral creature.", Name: "Invocation", IsAttack: true},
		string(PowerStorm): &PlayerPower{Type: PowerStorm, Speed: 15, Energy: 30, Strength: 30, Tile: Storm, Range: 1,
			Description: "Create a lighting storm in front of you.", Name: "Storm", IsAttack: true},
		string(PowerFlames): &PlayerPower{Type: PowerFlames, Speed: 3, Energy: 100, Strength: 60, Tile: Flames, Range: 5,
			Description: "Create a big fire storm all around you.", Name: "Flames", IsAttack: true},
		string(PowerHealing): &PlayerPower{Type: PowerHealing, Speed: 4, Energy: 20, Strength: 20, Tile: Healing,
			Description: "Heal yourself.", Name: "Healing", IsAttack: true},
		string(PowerDeadSpeaking): &PlayerPower{Type: PowerDeadSpeaking, Tile: Skull, Speed: 5, Energy: 20,
			Description: "Speak with the spirit of a dead.", Name: "Dead speaking"},
		string(PowerCalm): &PlayerPower{Type: PowerCalm, Tile: Calm, Speed: 10, Energy: 20, Strength: 20, Range: 1,
			Description: "Calm down a living being.", Name: "Calm", IsAttack: true},
	}
}

func (p *Player) NewPower(powername string, g *Game) *PlayerPower {
	p.newPowerRaw(powername)
	return p.Powers[powername]
}

func (p *Player) newPowerRaw(powername string) {
	_, exists := p.Powers[powername]
	if !exists {
		powers := Powers()
		p.Powers[powername] = powers[powername]
	}
}

func (p *Player) NextPower() {
	currentPowerIdx := 0
	powernames := p.GetSortedPowernames()
	for i, powername := range powernames {
		if powername == string(p.CurrentPower.Type) {
			currentPowerIdx = i
			break
		}
	}
	currentPowerIdx++
	if currentPowerIdx >= len(p.Powers) {
		currentPowerIdx = len(p.Powers) - 1
	}
	p.CurrentPower = p.Powers[powernames[currentPowerIdx]]
}

func (p *Player) LastPower() {
	currentPowerIdx := 0
	powernames := p.GetSortedPowernames()
	for i, powername := range powernames {
		if powername == string(p.CurrentPower.Type) {
			currentPowerIdx = i
			break
		}
	}
	currentPowerIdx--
	if currentPowerIdx <= 0 {
		currentPowerIdx = 0
	}
	p.CurrentPower = p.Powers[powernames[currentPowerIdx]]
}

func (p *Player) GetSortedPowernames() []string {
	powernames := make([]string, 0, len(p.Powers))
	for r, _ := range p.Powers {
		powernames = append(powernames, string(r))
	}
	sort.Strings(powernames) //sort by key
	return powernames
}

func (p *Player) getCurrentPowerIdx() int {
	currentPowerIdx := 0
	powernames := p.GetSortedPowernames()
	for i, powername := range powernames {
		if powername == string(p.CurrentPower.Type) {
			currentPowerIdx = i
			break
		}
	}
	return currentPowerIdx
}
