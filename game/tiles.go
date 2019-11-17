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
	Rat      Tile = "rat"
	Spider   Tile = "spider"
	Fox      Tile = "fox"
	Snake    Tile = "snake"
	Cat      Tile = "cat"
	Eagle    Tile = "eagle"
	Wolf     Tile = "wolf"
	Bear     Tile = "bear"
	Scorpion Tile = "scorpion"
	Bat      Tile = "bat"
	Elephant Tile = "elephant"
	Daemon   Tile = "daemon"
	Angel    Tile = "angel"
	Spirit   Tile = "spirit"
	Dragon   Tile = "dragon"
	God      Tile = "god"

	// Quests Objects
	Alcohol Tile = "alcohol"
	Herbs   Tile = "herbs"

	// Usables
	Gold   Tile = "gold"
	Book   Tile = "b"
	Senzu  Tile = "senzu"
	Bread  Tile = "bread"
	Fruits Tile = "fruits"
	Water  Tile = "water"
	Steak  Tile = "steak"

	// Powers & effects
	Explosion  Tile = "explosion"
	Flames     Tile = "flames"
	Storm      Tile = "storm"
	Healing    Tile = "healing"
	Teleport   Tile = "teleport"
	Skull      Tile = "skull"
	Necromancy Tile = "necromancy"
	Calm       Tile = "calm"
)
