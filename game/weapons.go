package game

import "time"

const (
	WeaponTypeDagger = "knife"
	WeaponTypeWand   = "wand"
	WeaponTypeBow    = "bow"
	WeaponTypeSpear  = "spear"
)

type Weapon struct {
	Typ             string
	Name            string
	Tile            Tile
	Damages         int
	Speed           int
	MagickalDamages int
}

func (w *Weapon) adaptSpeed() {
	time.Sleep(time.Duration(CharacterDeltaTime/w.Speed) * time.Millisecond)
}

func (p *Player) NextWeapon() {
	idx := 0
	for i, w := range p.Weapons {
		if w.Typ == string(p.Weapon.Typ) {
			idx = i
			break
		}
	}
	idx++
	if idx >= len(p.Weapons) {
		idx = len(p.Weapons) - 1
	}
	p.Weapon = p.Weapons[idx]
}

func (p *Player) LastWeapon() {
	idx := 0
	for i, w := range p.Weapons {
		if w.Typ == string(p.Weapon.Typ) {
			idx = i
			break
		}
	}
	idx--
	if idx <= 0 {
		idx = 0
	}
	p.Weapon = p.Weapons[idx]
}
