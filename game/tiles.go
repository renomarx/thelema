package game

type Tile rune

const (
	StoneWall      Tile = '#'
	DirtFloor      Tile = '.'
	DoorClosed     Tile = '|'
	DoorOpened     Tile = '/'
	Tree           Tile = '&'
	Blank          Tile = ' '
	Statue         Tile = 'é'
	Upstairs       Tile = '^'
	Downstairs     Tile = '~'
	Ocean          Tile = 'o'
	OceanTopSide   Tile = 'ô'
	OceanLeftSide  Tile = 'ò'
	OceanRightSide Tile = 'ó'
	OceanDownSide  Tile = 'õ'
	CityEntry      Tile = '>'
	CityOut        Tile = '<'

	// Houses
	HouseWall        Tile = 'h'
	HouseDoor        Tile = '='
	HouseTop         Tile = 'Ħ'
	HouseTopRight    Tile = 'Ĥ'
	HouseTopLeft     Tile = 'ĥ'
	HouseBottom      Tile = 'Ḫ'
	HouseBottomRight Tile = 'Ḩ'
	HouseBottomLeft  Tile = 'ḩ'
	HouseRight       Tile = 'Ḧ'
	HouseLeft        Tile = 'ḧ'

	// Monsters and invocations
	Rat      Tile = 'R'
	Spider   Tile = 'S'
	Fox      Tile = 'F'
	Snake    Tile = 'Ƨ'
	Cat      Tile = 'C'
	Eagle    Tile = 'A'
	Wolf     Tile = 'W'
	Bear     Tile = 'B'
	Scorpion Tile = 'P'
	Bat      Tile = 'T'
	Elephant Tile = 'E'
	Daemon   Tile = 'M'
	Angel    Tile = 'L'
	Spirit   Tile = 'I'
	Dragon   Tile = 'D'
	God      Tile = 'G'

	// Usables
	Senzu     Tile = 'z'
	Explosion Tile = 'x'
	Book      Tile = 'b'

	// Powers
	Energyball Tile = 'p'

	// Weapons
	Dagger Tile = 'd'
	Bow    Tile = 'a'
	Wand   Tile = 'w'
	Spear  Tile = 's'
)
