package game

type Tile string

const (
	StoneWall    Tile = "#"
	DirtFloor    Tile = "."
	HerbFloor    Tile = ":"
	CityFloor    Tile = "*"
	DoorClosed   Tile = "|"
	DoorOpened   Tile = "/"
	Tree         Tile = "&"
	Blank        Tile = " "
	Statue       Tile = "é"
	Window       Tile = "("
	WhiteWall    Tile = "$"
	Clothe       Tile = "à"
	Upstairs     Tile = "^"
	Downstairs   Tile = "~"
	Ocean        Tile = "o"
	CityEntry    Tile = ">"
	CityOut      Tile = "<"
	HouseDoor    Tile = "="
	PrisonDoor   Tile = "!"
	DungeonEntry Tile = "]"
	DungeonOut   Tile = "["

	// Monsters and invocations
	Rat      Tile = "R"
	Spider   Tile = "S"
	Fox      Tile = "F"
	Snake    Tile = "Ƨ"
	Cat      Tile = "C"
	Eagle    Tile = "A"
	Wolf     Tile = "W"
	Bear     Tile = "B"
	Scorpion Tile = "P"
	Bat      Tile = "T"
	Elephant Tile = "E"
	Daemon   Tile = "M"
	Angel    Tile = "L"
	Spirit   Tile = "I"
	Dragon   Tile = "D"
	God      Tile = "G"

	// Quests Objects
	Alcohol Tile = "c"
	Herbs   Tile = "v"

	// Usables
	Gold   Tile = "g"
	Book   Tile = "b"
	Senzu  Tile = "z"
	Bread  Tile = "p"
	Fruits Tile = "u"
	Water  Tile = "w"
	Steak  Tile = "s"

	// Powers & effects
	Explosion  Tile = "x"
	Flames     Tile = "f"
	Storm      Tile = "t"
	Healing    Tile = "h"
	Teleport   Tile = "y"
	Skull      Tile = "q"
	Necromancy Tile = "ç"
	Calm       Tile = "m"
)
