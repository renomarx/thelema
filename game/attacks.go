package game

type Attack struct {
	Damages    int
	Speed      int
	EnergyCost int
	Name       string
	Range      int
	Type       AttackType
	MagickType PowerType
}
type AttackType string

const AttackTypePhysical AttackType = "PHYSICAL"
const AttackTypeMagick AttackType = "MAGICK"

func Attacks() []*Attack {
	return []*Attack{
		&Attack{
			Speed:   20,
			Damages: 20,
			Name:    "Sword attack",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   15,
			Damages: 25,
			Name:    "Faint",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   9,
			Damages: 30,
			Name:    "Bear Assault",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   30,
			Damages: 30,
			Name:    "Snake attack",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   1,
			Damages: 40,
			Name:    "Counter",
			Type:    AttackTypePhysical,
			Range:   1,
		},
		&Attack{
			Speed:   7,
			Damages: 40,
			Name:    "Hurricane madness",
			Type:    AttackTypePhysical,
			Range:   3,
		},
		&Attack{
			Speed:   2,
			Damages: 100,
			Name:    "Fatal strike",
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

func (att *Attack) GetSpeed(p *Player) int {
	power := 0
	switch att.Type {
	case AttackTypePhysical:
		power = att.Speed * p.Dexterity.Current / 10
	case AttackTypeMagick:
		power = att.Speed * p.Intelligence.Current / 10
	}
	return power
}
