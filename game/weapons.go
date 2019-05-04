package game

import "time"

const (
	WeaponTypeDagger = "knife"
	WeaponTypeWand   = "wand"
	WeaponTypeBow    = "bow"
	WeaponTypeSpear  = "spear"
)

type Weapon struct {
	Typ     string
	Name    string
	Tile    rune
	Damages int
	Speed   int
}

func (w *Weapon) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/w.Speed) * time.Millisecond)
}
