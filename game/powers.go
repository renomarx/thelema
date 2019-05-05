package game

import "sort"

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
	Lifetime    int
	Range       int
	Description string
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
			p.Powers[string(PowerEnergyBall)] = &PlayerPower{Type: PowerEnergyBall, Speed: 10, Energy: 30, Tile: Energyball,
				Description: "Boule d'énergie. Explose quand recontre un obstacle ou un ennemi."}
		case PowerInvocation:
			p.Powers[string(PowerInvocation)] = &PlayerPower{Type: PowerInvocation, Strength: 100, Speed: 10, Energy: 100, Lifetime: 15, Tile: Fox,
				Description: "Invoque un familier qui attaque vos ennemis et attire leur attention."}
		case PowerStorm:
			p.Powers[string(PowerStorm)] = &PlayerPower{Type: PowerStorm, Speed: 10, Energy: 50, Lifetime: 2, Tile: Storm, Range: 7,
				Description: "Crée un éclair devant vous qui blesse vos ennemis durant un court laps de temps."}
		case PowerFlames:
			p.Powers[string(PowerFlames)] = &PlayerPower{Type: PowerFlames, Speed: 10, Energy: 200, Lifetime: 3, Tile: Flames, Range: 5,
				Description: "Crée un incendie tout autour de vous qui blesse vos ennemis durant un court laps de temps."}
		case PowerHealing:
			p.Powers[string(PowerHealing)] = &PlayerPower{Type: PowerHealing, Speed: 10, Energy: 20, Tile: Healing,
				Description: "Régénère votre santé."}
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
