package game

type MonsterType struct {
	Name        string
	Tile        Tile
	Health      int
	Energy      int
	Stats       int
	Speed       int
	Luck        int
	VisionRange int
	Probability int
	Weapon      *Weapon
}

func Bestiary() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Tile: Rat, Name: "Rat", Health: 50, Energy: 50, Stats: 5, Speed: 10, Luck: 10, VisionRange: 10, Probability: 100},
		&MonsterType{Tile: Spider, Name: "Spider", Health: 50, Energy: 50, Stats: 10, Speed: 5, Luck: 10, VisionRange: 10, Probability: 100},
		&MonsterType{Tile: Snake, Name: "Snake", Health: 100, Energy: 200, Stats: 10, Speed: 5, Luck: 0, VisionRange: 10, Probability: 80},
		&MonsterType{Tile: Cat, Name: "Cat", Health: 100, Energy: 300, Stats: 10, Speed: 10, Luck: 20, VisionRange: 10, Probability: 70},
		&MonsterType{Tile: Eagle, Name: "Eagle", Health: 150, Energy: 200, Stats: 10, Speed: 7, Luck: 25, VisionRange: 12, Probability: 50},
		&MonsterType{Tile: Wolf, Name: "Wolf", Health: 300, Energy: 300, Stats: 20, Speed: 8, Luck: 10, VisionRange: 10, Probability: 30},
		&MonsterType{Tile: Bear, Name: "Bear", Health: 500, Energy: 200, Stats: 40, Speed: 7, Luck: 20, VisionRange: 10, Probability: 20},
		&MonsterType{Tile: Elephant, Name: "Elephant", Health: 800, Energy: 200, Stats: 80, Speed: 7, Luck: 20, VisionRange: 8, Probability: 20},
	}
}

func BestiaryUnderworld() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Tile: Rat, Name: "Rat", Health: 50, Energy: 50, Stats: 5, Speed: 10, Luck: 10, VisionRange: 7, Probability: 100},
		&MonsterType{Tile: Spider, Name: "Spider", Health: 50, Energy: 50, Stats: 10, Speed: 5, Luck: 10, VisionRange: 7, Probability: 100},
		&MonsterType{Tile: Bat, Name: "Bat", Health: 50, Energy: 50, Stats: 10, Speed: 10, Luck: 10, VisionRange: 10, Probability: 80},
		&MonsterType{Tile: Scorpion, Name: "Scorpion", Health: 100, Energy: 50, Stats: 5, Speed: 5, Luck: 10, VisionRange: 5, Probability: 50},
		&MonsterType{Tile: Bear, Name: "Bear", Health: 500, Energy: 200, Stats: 40, Speed: 7, Luck: 20, VisionRange: 7, Probability: 20},
	}
}

func Creatures() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Name: "orc_thief", Health: 400, Energy: 100, Stats: 60, Speed: 10, Luck: 20, VisionRange: 12, Probability: 100,
			Weapon: &Weapon{Tile: Dagger, Name: "Dague", Typ: WeaponTypeDagger, Damages: 7, Speed: 20}},
		&MonsterType{Name: "orc_guard", Health: 800, Energy: 100, Stats: 80, Speed: 10, Luck: 20, VisionRange: 12, Probability: 20,
			Weapon: &Weapon{Tile: Spear, Name: "Lance", Typ: WeaponTypeSpear, Damages: 20, Speed: 12}},
	}
}

func CreaturesUnderworld() []*MonsterType {
	return []*MonsterType{
		&MonsterType{Name: "skeleton_warrior", Health: 400, Energy: 100, Stats: 60, Speed: 10, Luck: 20, VisionRange: 12, Probability: 100,
			Weapon: &Weapon{Tile: Spear, Name: "Lance", Typ: WeaponTypeSpear, Damages: 20, Speed: 12}},
		&MonsterType{Name: "skeleton_sorcerer", Health: 200, Energy: 1000, Stats: 60, Speed: 10, Luck: 10, VisionRange: 15, Probability: 40,
			Weapon: &Weapon{Tile: Dagger, Name: "Dague", Typ: WeaponTypeDagger, Damages: 7, Speed: 20}},
		&MonsterType{Name: "skeleton_lord", Health: 1000, Energy: 200, Stats: 80, Speed: 10, Luck: 10, VisionRange: 12, Probability: 20,
			Weapon: &Weapon{Tile: Spear, Name: "Lance", Typ: WeaponTypeSpear, Damages: 20, Speed: 12}},
	}
}
