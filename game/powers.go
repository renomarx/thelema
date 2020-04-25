package game

import (
	"log"
	"sort"
	"time"
)

const PowerInvocation = "invocation"
const PowerFlames = "flames"
const PowerStorm = "storm"
const PowerHealing = "healing"
const PowerDeadSpeaking = "dead_speaking"
const PowerCalm = "calm"

const MagickCategoryPhysical = "physical"
const MagickCategoryAstral = "astral"
const MagickCategoryMental = "mental"
const MagickCategoryHigh = "high"
const MagickCategoryMeta = "meta"

const MagickElementEarth = "earth"
const MagickElementWater = "water"
const MagickElementAir = "air"
const MagickElementFire = "fire"
const MagickElementEther = "ether"

type PlayerPower struct {
	UID         string
	Tile        Tile
	Strength    int
	Speed       int
	Energy      int
	Range       int
	Description string
	Name        string
	IsAttack    bool
	Category    MagickCategory
	Element     MagickElement
}

type Magicks []PlayerPower

type MagickCategory string

type MagickElement string

func (pow *PlayerPower) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/pow.Speed) * time.Millisecond)
}

func (m Magicks) GetPower(uid string) *PlayerPower {
	for _, pow := range m {
		if pow.UID == uid {
			return &pow
		}
	}
	return nil
}

func Powers() Magicks {
	return Magicks{
		PlayerPower{
			UID:         PowerStorm,
			Speed:       15,
			Energy:      30,
			Strength:    30,
			Tile:        Storm,
			Range:       1,
			Description: "Crée un éclair entre vous et la cible.",
			Name:        "Eclair",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementAir,
		},
		PlayerPower{
			UID:         PowerHealing,
			Speed:       4,
			Energy:      20,
			Strength:    20,
			Tile:        Healing,
			Description: "Régénère la santé.",
			Name:        "Soin",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementEarth,
		},
		PlayerPower{
			UID:         PowerCalm,
			Tile:        Calm,
			Speed:       10,
			Energy:      20,
			Strength:    20,
			Range:       1,
			Description: "Calme un être vivant.",
			Name:        "Calme",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementWater,
		},
		PlayerPower{
			UID:         PowerInvocation,
			Speed:       5,
			Energy:      50,
			Tile:        Spirit,
			Description: "Invoque une créature astrale.",
			Name:        "Invocation",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementFire,
		},
		PlayerPower{
			UID:         PowerDeadSpeaking,
			Tile:        Skull,
			Speed:       5,
			Energy:      20,
			Description: "Rappelle l'esprit d'un mort.",
			Name:        "Nécromancie",
			Category:    MagickCategoryMental,
			Element:     MagickElementEther,
		},
		PlayerPower{
			UID:         PowerFlames,
			Speed:       3,
			Energy:      100,
			Strength:    60,
			Tile:        Flames,
			Range:       5,
			Description: "Crée un gigantesque incendie sur vos cibles.",
			Name:        "Incendie",
			IsAttack:    true,
			Category:    MagickCategoryHigh,
			Element:     MagickElementFire,
		},
	}
}

func (p *Player) NewPower(powername string, g *Game) *PlayerPower {
	p.newPowerRaw(powername)
	return p.Powers[powername]
}

func (p *Player) newPowerRaw(powername string) {
	_, exists := p.Powers[powername]
	if !exists {
		magicks := Powers()
		pow := magicks.GetPower(powername)
		if pow == nil {
			log.Printf("Error: power %s does not exist", powername)
		} else {
			p.Powers[powername] = pow
		}
	}
}

func (p *Player) NextPower() {
	currentPowerIdx := 0
	powernames := p.GetSortedPowernames()
	for i, powername := range powernames {
		if powername == string(p.CurrentPower.UID) {
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
		if powername == string(p.CurrentPower.UID) {
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
		if powername == string(p.CurrentPower.UID) {
			currentPowerIdx = i
			break
		}
	}
	return currentPowerIdx
}
