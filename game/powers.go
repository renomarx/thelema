package game

import (
	"log"
	"sort"
	"time"
)

const PowerBrutalStrength = "brutal_strength"
const PowerQuickening = "quickening"
const PowerRockBody = "rock_body"

const PowerCharm = "charm"
const PowerGlaciation = "glaciation"

const PowerStorm = "storm"
const PowerLightness = "lightness"

const PowerKiai = "kiai"
const PowerFireball = "fireball"

const PowerPoison = "poison"
const PowerRoot = "root"
const PowerWeakening = "weakening"
const PowerVampirism = "vampirism"

const PowerHealing = "healing"
const PowerRecovery = "recovery"

const PowerCalm = "calm"
const PowerMadness = "madness"
const PowerFear = "fear"

const PowerMeditation = "meditation"
const PowerConcentration = "concentration"

const PowerInvocation = "invocation"
const PowerBanishment = "banishment"

const PowerCurse = "curse"
const PowerPossession = "possession"

const PowerLevitation = "levitation"
const PowerParalysis = "paralysis"

const PowerHypnosis = "hypnosis"
const PowerProtection = "protection"

const PowerInvisibility = "invisibility"
const PowerPrediction = "prediction"
const PowerSplit = "split"

const PowerCreation = "creation"
const PowerDevotion = "devotion"

const PowerDomination = "domination"
const PowerDeadSpeaking = "dead_speaking"

const PowerEarthquake = "earthquake"
const PowerGolem = "golem"

const PowerRevolution = "revolution"
const PowerPeace = "peace"

const PowerTempest = "tempest"
const PowerDivination = "divination"

const PowerTeleportation = "teleportation"
const PowerDisintegration = "disintegration"
const PowerFlames = "flames"

const PowerResurrection = "resurrection"
const PowerDeath = "death"

const MagickCategoryPhysical = "physical"
const MagickCategoryAstral = "astral"
const MagickCategoryMental = "mental"
const MagickCategoryHigh = "high"

const MagickElementEarth = "earth"
const MagickElementWater = "water"
const MagickElementAir = "air"
const MagickElementFire = "fire"
const MagickElementEther = "ether"

type Power struct {
	UID         string
	Tile        Tile
	Strength    int
	Speed       int
	Energy      int
	Range       int
	Description string
	Name        string
	IsAttack    bool
	IsEffect    bool
	Category    MagickCategory
	Element     MagickElement
}

type Magicks []Power

type MagickCategory string

type MagickElement string

func (pow *Power) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/pow.Speed) * time.Millisecond)
}

func (m Magicks) GetPower(uid string) *Power {
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
		Power{
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
		Power{
			UID:         PowerQuickening,
			Speed:       20,
			Energy:      20,
			Strength:    2,
			Tile:        QuickeningIcon,
			Name:        "Hâte",
			Description: "Augmente votre vitesse.",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementEarth,
		},
		Power{
			UID:         PowerRockBody,
			Speed:       20,
			Energy:      40,
			Strength:    5,
			Tile:        RockBodyIcon,
			Name:        "Corps de pierre",
			Description: "Augmente votre défense.",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementEarth,
		},
		// Water
		Power{
			UID:         PowerCharm,
			Speed:       20,
			Energy:      10,
			Strength:    5,
			Tile:        CharmIcon,
			Range:       1,
			Description: "Augmente votre beauté.",
			Name:        "Charme",
			IsAttack:    false,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementWater,
		},
		Power{
			UID:         PowerGlaciation,
			Speed:       10,
			Energy:      30,
			Strength:    20,
			Tile:        GlaciationIcon,
			Range:       1,
			Description: "Glace votre cible.",
			Name:        "Glaciation",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementWater,
		},
		// Air
		Power{
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
		Power{
			UID:         PowerLightness,
			Speed:       20,
			Energy:      40,
			Strength:    3,
			Tile:        LightnessIcon,
			Name:        "Légèreté",
			Description: "Augmente votre esquive.",
			IsAttack:    true,
			Category:    MagickCategoryPhysical,
			Element:     MagickElementAir,
		},

		// Astral
		// Earth
		Power{
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
		Power{
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
		Power{
			UID:         PowerInvocation,
			Speed:       5,
			Energy:      50,
			Strength:    20,
			Tile:        InvocationIcon,
			Description: "Invoque une créature astrale.",
			Name:        "Invocation",
			IsAttack:    true,
			Category:    MagickCategoryAstral,
			Element:     MagickElementFire,
		},

		// Mental
		// Ether
		Power{
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
		Power{
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

func (p *Player) NewPower(powername string) *Power {
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
