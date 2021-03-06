package game

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

func newPlayer(health, energy, stg, dex, bea, wil, intel, cha, rg int) *Player {
	player := &Player{}
	player.Character = NewCharacter()
	player.Speed.Init(10)
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
	attacks := Attacks()
	player.Attacks = attacks[:2]
	// TODO : load only first powers
	player.NewPower(PowerBrutalStrength)
	player.NewPower(PowerQuickening)
	player.NewPower(PowerRockBody)
	player.NewPower(PowerCharm)
	player.NewPower(PowerGlaciation)
	player.NewPower(PowerHealing)
	player.NewPower(PowerFlames)
	player.NewPower(PowerStorm)
	player.NewPower(PowerLightness)
	player.NewPower(PowerInvocation)
	player.NewPower(PowerCalm)
	player.CurrentPower = player.Powers[PowerHealing]

	player.LoadPlayerMenu()
	player.Inventory = NewInventory()
	player.Library = NewLibrary()

	return player
}

func NewAsmodeus() *Player {
	p := newPlayer(300, 200, 30, 20, 20, 30, 20, 20, 1)
	p.Name = "asmodeus"
	p.Affinity = "Asmodeus"
	return p
}

func NewKali() *Player {
	p := newPlayer(200, 300, 20, 30, 20, 20, 30, 20, 1)
	p.Name = "kali"
	p.Affinity = "Kali"
	return p
}

func NewBaal() *Player {
	p := newPlayer(300, 300, 20, 20, 20, 20, 20, 20, 1)
	p.Name = "baal"
	p.Affinity = "Baal"
	return p
}

func NewLilith() *Player {
	p := newPlayer(225, 225, 25, 25, 25, 25, 25, 25, 1)
	p.Name = "lilith"
	p.Affinity = "Lilith"
	return p
}

func NewDagon() *Player {
	p := newPlayer(200, 200, 30, 30, 20, 30, 30, 20, 1)
	p.Name = "dagon"
	p.Affinity = "Dagon"
	return p
}

func NewLucifer() *Player {
	p := newPlayer(200, 400, 20, 20, 20, 20, 20, 20, 1)
	p.Name = "lucifer"
	p.Affinity = "Lucifer"
	return p
}

func NewHecate() *Player {
	p := newPlayer(200, 300, 20, 20, 30, 20, 20, 30, 1)
	p.Name = "hecate"
	p.Affinity = "Hecate"
	return p
}
