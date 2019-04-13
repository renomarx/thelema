package game

import "sort"

type PowerType string

const PowerEnergyBall = "energy_ball"
const PowerInvocation = "invocation"

type PlayerPower struct {
	Type     PowerType
	Tile     Tile
	Strength int
	Speed    int
	Energy   int
	Lifetime int
}

func (p *Player) NewPower(powername string) {
	_, exists := p.Powers[powername]
	if exists {
		return
	}
	switch powername {
	case PowerEnergyBall:
		p.Powers[string(PowerEnergyBall)] = &PlayerPower{Type: PowerEnergyBall, Strength: 50, Speed: 10, Energy: 10, Tile: Energyball}
	case PowerInvocation:
		p.Powers[string(PowerInvocation)] = &PlayerPower{Type: PowerInvocation, Strength: 100, Speed: 10, Energy: 100, Lifetime: 15, Tile: Fox}
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
