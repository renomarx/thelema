package game

type MonsterType struct {
	Name        string
	Tile        Tile
	Hitpoints   int
	Strength    int
	Speed       int
	VisionRange int
}

func Bestiary() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Tile: Rat, Name: "Rat", Hitpoints: 50, Strength: 5, Speed: 10, VisionRange: 5},
		&MonsterType{Tile: Spider, Name: "Spider", Hitpoints: 50, Strength: 10, Speed: 5, VisionRange: 5},
	}
}
