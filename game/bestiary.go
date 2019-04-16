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

func NewRat(p Pos) *Monster {
	mt := &MonsterType{}
	mt.Tile = Rat
	mt.Name = "Rat"
	mt.Hitpoints = 50
	mt.Strength = 5
	mt.Speed = 10
	mt.VisionRange = 5
	return NewMonster(mt, p)
}

func NewSpider(p Pos) *Monster {
	mt := &MonsterType{}
	mt.Tile = Spider
	mt.Name = "Spider"
	mt.Hitpoints = 50
	mt.Strength = 10
	mt.Speed = 5
	mt.VisionRange = 5
	return NewMonster(mt, p)
}
