package game

type MonsterType struct {
	Name           string
	Tile           Tile
	Health         int
	Energy         int
	Strength       int
	Speed          int
	Luck           int
	VisionRange    int
	Probability    int
	Aggressiveness int
}

func Bestiary() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Tile: Rat, Name: "Rat", Health: 50, Energy: 50, Strength: 5, Speed: 10, Luck: 10, VisionRange: 10, Probability: 100, Aggressiveness: 200},
		&MonsterType{Tile: Spider, Name: "Spider", Health: 50, Energy: 50, Strength: 10, Speed: 5, Luck: 10, VisionRange: 10, Probability: 100, Aggressiveness: 100},
		&MonsterType{Tile: Snake, Name: "Snake", Health: 100, Energy: 200, Strength: 10, Speed: 5, Luck: 0, VisionRange: 10, Probability: 80, Aggressiveness: 75},
		&MonsterType{Tile: Cat, Name: "Cat", Health: 100, Energy: 300, Strength: 10, Speed: 10, Luck: 20, VisionRange: 10, Probability: 70, Aggressiveness: 50},
		&MonsterType{Tile: Eagle, Name: "Eagle", Health: 150, Energy: 200, Strength: 10, Speed: 7, Luck: 25, VisionRange: 12, Probability: 50, Aggressiveness: 80},
		&MonsterType{Tile: Wolf, Name: "Wolf", Health: 300, Energy: 300, Strength: 20, Speed: 8, Luck: 10, VisionRange: 10, Probability: 30, Aggressiveness: 400},
		&MonsterType{Tile: Bear, Name: "Bear", Health: 500, Energy: 200, Strength: 40, Speed: 7, Luck: 20, VisionRange: 10, Probability: 20, Aggressiveness: 200},
		&MonsterType{Tile: Elephant, Name: "Elephant", Health: 800, Energy: 200, Strength: 80, Speed: 7, Luck: 20, VisionRange: 8, Probability: 20, Aggressiveness: 100},
	}
}

func BestiaryUnderworld() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Tile: Rat, Name: "Rat", Health: 50, Energy: 50, Strength: 5, Speed: 10, Luck: 10, VisionRange: 7, Probability: 100, Aggressiveness: 200},
		&MonsterType{Tile: Spider, Name: "Spider", Health: 50, Energy: 50, Strength: 10, Speed: 5, Luck: 10, VisionRange: 7, Probability: 100, Aggressiveness: 100},
		&MonsterType{Tile: Bat, Name: "Bat", Health: 50, Energy: 50, Strength: 10, Speed: 10, Luck: 10, VisionRange: 10, Probability: 80, Aggressiveness: 300},
		&MonsterType{Tile: Scorpion, Name: "Scorpion", Health: 100, Energy: 50, Strength: 5, Speed: 5, Luck: 10, VisionRange: 5, Probability: 50, Aggressiveness: 75},
		&MonsterType{Tile: Bear, Name: "Bear", Health: 500, Energy: 200, Strength: 40, Speed: 7, Luck: 20, VisionRange: 7, Probability: 20, Aggressiveness: 200},
	}
}
