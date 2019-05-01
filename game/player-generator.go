package game

import "time"

func GeneratePlayers() []*Player {
	var players []*Player
	players = append(players, NewSayan())
	players = append(players, NewMonk())
	players = append(players, NewNamek())
	players = append(players, NewCyborg())
	return players
}

func newPlayer(speed, health, energy, strength, rg int) *Player {
	player := &Player{}
	player.X = 0
	player.Y = 0
	player.Xb = 0
	player.Yb = 0
	player.IsMoving = false
	player.Speed.Init(speed)
	player.Health.Init(health)
	player.Energy.Init(energy)
	player.Strength.Init(strength)
	player.RegenerationSpeed.Init(rg)
	player.IsAttacking = false
	player.IsPowerAttacking = false
	player.IsTalking = false
	player.IsTaking = false
	player.Inventory = NewInventory()
	player.Library = NewLibrary()
	player.Powers = make(map[string]*PlayerPower)
	player.Powers[string(PowerEnergyBall)] = &PlayerPower{Type: PowerEnergyBall, Strength: 50, Speed: 10, Energy: 10, Tile: Energyball}
	player.CurrentPower = player.Powers[string(PowerEnergyBall)]
	player.LastRegenerationTime = time.Now()
	player.LoadPlayerMenu()

	return player
}

func NewSayan() *Player {
	p := newPlayer(12, 120, 120, 25, 1)
	p.Name = "sayan"
	return p
}

func NewMonk() *Player {
	p := newPlayer(14, 100, 100, 20, 1)
	p.Name = "monk"
	return p
}

func NewNamek() *Player {
	p := newPlayer(10, 100, 100, 10, 2)
	p.Name = "namek"
	return p
}

func NewCyborg() *Player {
	p := newPlayer(8, 100, 200, 30, 1)
	p.Name = "cyborg"
	return p
}
