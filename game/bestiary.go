package game

type MonsterType struct {
	Name        string
	Tile        Tile
	Hitpoints   int
	Strength    int
	Speed       int
	VisionRange int
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
