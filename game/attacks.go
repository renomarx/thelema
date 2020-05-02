package game

type Attack struct {
	UID            string
	Strength       int
	Speed          int
	EnergyCost     int
	Name           string
	Range          int
	Accuracy       int
	Type           AttackType
	MagickUID      string
	MagickElement  MagickElement
	MagickCategory MagickCategory
}
type AttackType string

const AttackTypePhysical AttackType = "PHYSICAL"
const AttackTypeMagick AttackType = "MAGICK"

func Attacks() []Attack {
	return []Attack{
		Attack{
			UID:      "charge",
			Speed:    20,
			Strength: 20,
			Name:     "Charge",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 95,
		},
		Attack{
			UID:      "faint",
			Speed:    15,
			Strength: 25,
			Name:     "Feinte",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 100,
		},
		Attack{
			UID:      "bear_assault",
			Speed:    9,
			Strength: 40,
			Name:     "Charge de l'ours",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 80,
		},
		Attack{
			UID:      "fast_attack",
			Speed:    40,
			Strength: 15,
			Name:     "Coup rapide",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 100,
		},
		Attack{
			UID:      "counter",
			Speed:    1,
			Strength: 50,
			Name:     "Contre",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 100,
		},
		Attack{
			UID:      "fury",
			Speed:    7,
			Strength: 40,
			Name:     "Furie",
			Type:     AttackTypePhysical,
			Range:    3,
			Accuracy: 75,
		},
		Attack{
			UID:      "vital_points",
			Speed:    2,
			Strength: 70,
			Name:     "Points vitaux",
			Type:     AttackTypePhysical,
			Range:    1,
			Accuracy: 50,
		},
	}
}

func (att *Attack) GetPower(p *Character) int {
	power := 0
	switch att.Type {
	case AttackTypePhysical:
		power = att.Strength * p.CalculateAttackScore() / 10
	case AttackTypeMagick:
		power = (att.Strength*p.CalculatePowerAttackScore() + p.GetElementalAffinity(att.MagickElement)) / 10
	}
	return power
}

func (att *Attack) GetSpeed(p *Character) int {
	power := 0
	switch att.Type {
	case AttackTypePhysical:
		power = att.Speed * p.Dexterity.Current / 10
	case AttackTypeMagick:
		power = att.Speed * p.Intelligence.Current / 10
	}
	return power
}
