package game

import "time"

func GeneratePlayers() []*Player {
	var players []*Player
	players = append(players, NewAsmodeus())
	players = append(players, NewKali())
	players = append(players, NewBaal())
	players = append(players, NewLilith())
	players = append(players, NewDagon())
	players = append(players, NewLucifer())
	players = append(players, NewHecate())
	return players
}

func newPlayer(speed, health, energy, stg, dex, bea, wil, intel, cha, rg int, weapon string) *Player {
	player := &Player{}
	player.X = 0
	player.Y = 0
	player.Xb = 0
	player.Yb = 0
	player.Speed.Init(speed)
	player.Health.Init(health)
	player.Energy.Init(energy)
	player.Strength.Init(stg)
	player.Dexterity.Init(dex)
	player.Beauty.Init(bea)
	player.Will.Init(wil)
	player.Intelligence.Init(intel)
	player.Charisma.Init(cha)
	player.RegenerationSpeed.Init(rg)
	player.Luck.Init(20)
	player.IsAttacking = false
	player.IsPowerAttacking = false
	player.IsTalking = false
	player.IsTaking = false
	player.Inventory = NewInventory()
	player.Library = NewLibrary()
	player.Powers = make(map[string]*PlayerPower)
	// FIXME
	player.newPowerRaw(PowerHealing)
	player.newPowerRaw(PowerEnergyBall)
	player.newPowerRaw(PowerFlames)
	player.newPowerRaw(PowerStorm)
	player.newPowerRaw(PowerInvocation)
	player.CurrentPower = player.Powers[string(PowerEnergyBall)]
	player.LastRegenerationTime = time.Now()
	player.LoadPlayerMenu()
	player.Weapons = append(player.Weapons, &Weapon{Tile: Dagger, Name: "Dague", Typ: WeaponTypeDagger, Damages: 7, Speed: 20})
	player.Weapons = append(player.Weapons, &Weapon{Tile: Bow, Name: "Arc", Typ: WeaponTypeBow, Damages: 10, Speed: 5})
	player.Weapons = append(player.Weapons, &Weapon{Tile: Wand, Name: "Sceptre", Typ: WeaponTypeWand, Speed: 10, MagickalDamages: 20})
	player.Weapons = append(player.Weapons, &Weapon{Tile: Spear, Name: "Lance", Typ: WeaponTypeSpear, Damages: 20, Speed: 12})
	switch weapon {
	case WeaponTypeDagger:
		player.Weapon = player.Weapons[0]
	case WeaponTypeBow:
		player.Weapon = player.Weapons[1]
	case WeaponTypeWand:
		player.Weapon = player.Weapons[2]
	case WeaponTypeSpear:
		player.Weapon = player.Weapons[3]
	}

	return player
}

func NewAsmodeus() *Player {
	p := newPlayer(10, 300, 200, 30, 20, 20, 30, 20, 20, 1, WeaponTypeSpear)
	p.Name = "asmodeus"
	p.Affinity = "Asmodeus"
	return p
}

func NewKali() *Player {
	p := newPlayer(10, 200, 300, 20, 30, 20, 20, 30, 20, 1, WeaponTypeWand)
	p.Name = "kali"
	p.Affinity = "Kali"
	return p
}

func NewBaal() *Player {
	p := newPlayer(10, 300, 300, 20, 20, 20, 20, 20, 20, 1, WeaponTypeDagger)
	p.Name = "baal"
	p.Affinity = "Baal"
	return p
}

func NewLilith() *Player {
	p := newPlayer(10, 225, 225, 25, 25, 25, 25, 25, 25, 1, WeaponTypeBow)
	p.Name = "lilith"
	p.Affinity = "Lilith"
	return p
}

func NewDagon() *Player {
	p := newPlayer(10, 200, 200, 30, 30, 20, 30, 30, 20, 1, WeaponTypeWand)
	p.Name = "dagon"
	p.Affinity = "Dagon"
	return p
}

func NewLucifer() *Player {
	p := newPlayer(10, 200, 400, 20, 20, 20, 20, 20, 20, 1, WeaponTypeWand)
	p.Name = "lucifer"
	p.Affinity = "Lucifer"
	return p
}

func NewHecate() *Player {
	p := newPlayer(10, 200, 300, 20, 20, 30, 20, 20, 30, 1, WeaponTypeWand)
	p.Name = "hecate"
	p.Affinity = "Hecate"
	return p
}
