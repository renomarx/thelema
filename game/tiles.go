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
	Upstairs       Tile = 'u'
	Downstairs     Tile = 'd'
	Explosion      Tile = 'x'
	Book           Tile = 'b'
	Ocean          Tile = 'o'
	OceanTopSide   Tile = 'ô'
	OceanLeftSide  Tile = 'ò'
	OceanRightSide Tile = 'ó'
	OceanDownSide  Tile = 'õ'
)

// Monsters and invocations
const (
	Rat    Tile = 'R'
	Spider Tile = 'S'
	Fox    Tile = 'F'
)

// Houses : TODO : only one tile, move others in ui2d
const (
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
)

// Pnjs
const (
	Jason Tile = 'J'
	Sarah Tile = 'A'
)

// Usables
const (
	Senzu Tile = 'z'
)

// Powers
const (
	Energyball Tile = 'p'
)
