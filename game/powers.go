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
const PowerBrutalStrength = "brutal_strength"
const PowerQuickening = "quickening"

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
		// Physical
		//Earth
		PlayerPower{
			UID:         PowerBrutalStrength,
			Speed:       20,
			Energy:      30,
			Strength:    2,
			Tile:        BrutalStrengthIcon,
			Name:        "Force brute",
			Description: "Augmente votre force.",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementEarth,
		},
		PlayerPower{
			UID:         PowerQuickening,
			Speed:       20,
			Energy:      30,
			Strength:    2,
			Tile:        QuickeningIcon,
			Name:        "Hâte",
			Description: "Augmente votre vitesse.",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementEarth,
		},
		// Air
		PlayerPower{
			UID:         PowerStorm,
			Speed:       15,
			Energy:      30,
			Strength:    30,
			Tile:        StormIcon,
			Range:       1,
			Description: "Electrocute la cible.",
			Name:        "Eclair",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementAir,
		},

		// Astral
		// Earth
		PlayerPower{
			UID:         PowerHealing,
			Speed:       4,
			Energy:      20,
			Strength:    20,
			Tile:        HealingIcon,
			Description: "Régénère la santé.",
			Name:        "Soin",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementEarth,
		},
		// Water
		PlayerPower{
			UID:         PowerCalm,
			Tile:        CalmIcon,
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
		// Fire
		PlayerPower{
			UID:         PowerInvocation,
			Speed:       5,
			Energy:      50,
			Tile:        InvocationIcon,
			Description: "Invoque une créature astrale.",
			Name:        "Invocation",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementFire,
		},

		// Mental
		// Ether
		PlayerPower{
			UID:         PowerDeadSpeaking,
			Tile:        NecromancyIcon,
			Speed:       5,
			Energy:      20,
			Description: "Rappelle l'esprit d'un mort.",
			Name:        "Nécromancie",
			Category:    MagickCategoryMental,
			Element:     MagickElementEther,
		},

		// High
		// Fire
		PlayerPower{
			UID:         PowerFlames,
			Speed:       3,
			Energy:      200,
			Strength:    90,
			Tile:        FlamesIcon,
			Range:       5,
			Description: "Crée un gigantesque incendie sur vos cibles.",
			Name:        "Incendie",
			IsAttack:    true,
			Category:    MagickCategoryHigh,
			Element:     MagickElementFire,
		},
	}
}

func (p *Player) NewPower(powername string) *PlayerPower {
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
	return p.Powers[powername]
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
