package game

type Attack struct {
	Damages    int
	Speed      int
	EnergyCost int
	Name       string
	Range      int
	Type       AttackType
	MagickType string
}
type AttackType string

const AttackTypePhysical AttackType = "PHYSICAL"
const AttackTypeMagick AttackType = "MAGICK"

func Attacks() []*Attack {
	return []*Attack{
		&Attack{
			Speed:   20,
			Damages: 20,
			Name:    "Charge",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   15,
			Damages: 25,
			Name:    "Feinte",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   9,
			Damages: 30,
			Name:    "Charge de l'ours",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   30,
			Damages: 30,
			Name:    "Coup rapide",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   1,
			Damages: 40,
			Name:    "Contre",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   7,
			Damages: 40,
			Name:    "Attaques furie",
			Type:    AttackTypePhysical,
			Range:   3,
		},
		&Attack{
			Speed:   2,
			Damages: 100,
			Name:    "Points vitaux",
			Type:    AttackTypePhysical,
			Range:   1,
		},
	}
}

func (att *Attack) GetPower(p *Character) int {
	power := 0
	switch att.Type {
	case AttackTypePhysical:
		power = att.Damages * p.CalculateAttackScore() / 10
	case AttackTypeMagick:
		power = att.Damages * p.CalculatePowerAttackScore() / 10
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
