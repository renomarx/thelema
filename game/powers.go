package game

import "sort"
import "time"

type PowerType string

const PowerInvocation = "invocation"
const PowerFlames = "flames"
const PowerStorm = "storm"
const PowerHealing = "healing"
const PowerDeadSpeaking = "dead_speaking"

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

func (p *Player) NewPower(powername string, g *Game) {
	p.newPowerRaw(powername)
	pp := p.Powers[powername]
	EM.Dispatch(&Event{
		Message: "You learned power: '" + string(pp.Type) + "' with this book!",
		Action:  ActionPower,
		Payload: map[string]string{"type": string(pp.Type)}})
}

func (p *Player) newPowerRaw(powername string) {
	_, exists := p.Powers[powername]
	if !exists {
		switch powername {
		case PowerInvocation:
			p.Powers[string(PowerInvocation)] = &PlayerPower{Type: PowerInvocation, Speed: 5, Energy: 50, Tile: Spirit,
				Description: "Invoque un esprit.", Name: "Invocation", IsAttack: true}
		case PowerStorm:
			p.Powers[string(PowerStorm)] = &PlayerPower{Type: PowerStorm, Speed: 15, Energy: 30, Strength: 30, Tile: Storm, Range: 1,
				Description: "Lance un éclair devant vous.", Name: "Storm", IsAttack: true}
		case PowerFlames:
			p.Powers[string(PowerFlames)] = &PlayerPower{Type: PowerFlames, Speed: 3, Energy: 100, Strength: 60, Tile: Flames, Range: 5,
				Description: "Crée un incendie tout autour de vous.", Name: "Flames", IsAttack: true}
		case PowerHealing:
			p.Powers[string(PowerHealing)] = &PlayerPower{Type: PowerHealing, Speed: 4, Energy: 20, Strength: 20, Tile: Healing,
				Description: "Régénère votre santé.", Name: "Healing", IsAttack: true}
		case PowerDeadSpeaking:
			p.Powers[string(PowerDeadSpeaking)] = &PlayerPower{Type: PowerDeadSpeaking, Tile: Skull, Speed: 5, Energy: 20,
				Description: "Parle avec l'esprit d'un mort.", Name: "Dead speaking"}
		}
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
