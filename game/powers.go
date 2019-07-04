package game

import "sort"
import "time"

type PowerType string

const PowerEnergyBall = "energy_ball"
const PowerInvocation = "invocation"
const PowerFlames = "flames"
const PowerStorm = "storm"
const PowerHealing = "healing"

type PlayerPower struct {
	Type        PowerType
	Tile        Tile
	Strength    int
	Speed       int
	Energy      int
	Range       int
	Description string
	Name        string
}

func (pow *PlayerPower) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/pow.Speed) * time.Millisecond)
}

func (p *Player) NewPower(powername string, g *Game) {
	p.newPowerRaw(powername)
	pp := p.Powers[powername]
	g.GetEventManager().Dispatch(&Event{
		Message: "You learned power: '" + string(pp.Type) + "' with this book!",
		Action:  ActionPower,
		Payload: map[string]string{"type": string(pp.Type)}})
}

func (p *Player) newPowerRaw(powername string) {
	_, exists := p.Powers[powername]
	if !exists {
		switch powername {
		case PowerEnergyBall:
			p.Powers[string(PowerEnergyBall)] = &PlayerPower{Type: PowerEnergyBall, Speed: 10, Strength: 20, Energy: 10, Tile: Energyball, Range: 1,
				Description: "Boule d'énergie. Explose quand recontre un obstacle ou un ennemi.", Name: "Energy ball"}
		case PowerInvocation:
			p.Powers[string(PowerInvocation)] = &PlayerPower{Type: PowerInvocation, Speed: 5, Energy: 50, Tile: Fox,
				Description: "Invoque un familier qui attaque vos ennemis et attire leur attention.", Name: "Invocation"}
		case PowerStorm:
			p.Powers[string(PowerStorm)] = &PlayerPower{Type: PowerStorm, Speed: 15, Energy: 30, Strength: 30, Tile: Storm, Range: 1,
				Description: "Crée un éclair devant vous qui blesse vos ennemis durant un court laps de temps.", Name: "Storm"}
		case PowerFlames:
			p.Powers[string(PowerFlames)] = &PlayerPower{Type: PowerFlames, Speed: 3, Energy: 100, Strength: 60, Tile: Flames, Range: 5,
				Description: "Crée un incendie tout autour de vous qui blesse vos ennemis durant un court laps de temps.", Name: "Flames"}
		case PowerHealing:
			p.Powers[string(PowerHealing)] = &PlayerPower{Type: PowerHealing, Speed: 4, Energy: 20, Strength: 20, Tile: Healing,
				Description: "Régénère votre santé.", Name: "Healing"}
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
