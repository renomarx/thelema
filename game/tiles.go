package game

type Tile rune

const (
	StoneWall      Tile = '#'
	DirtFloor      Tile = '.'
	HerbFloor      Tile = ':'
	GreenFloor     Tile = '*'
	DoorClosed     Tile = '|'
	DoorOpened     Tile = '/'
	Tree           Tile = '&'
	Blank          Tile = ' '
	Statue         Tile = 'é'
	Window         Tile = '('
	WhiteWall      Tile = '$'
	Clothe         Tile = 'à'
	Upstairs       Tile = '^'
	Downstairs     Tile = '~'
	Ocean          Tile = 'o'
	OceanTopSide   Tile = 'ô'
	OceanLeftSide  Tile = 'ò'
	OceanRightSide Tile = 'ó'
	OceanDownSide  Tile = 'õ'
	CityEntry      Tile = '>'
	CityOut        Tile = '<'
	HouseDoor      Tile = '='
	PrisonDoor     Tile = '!'

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

	// Quests Objects
	Alcohol Tile = 'c'
	Herbs   Tile = 'v'

	// Usables
	Gold      Tile = 'g'
	Explosion Tile = 'x'
	Book      Tile = 'b'
	Senzu     Tile = 'z'
	Bread     Tile = 'i'
	Fruits    Tile = 'u'
	Water     Tile = 'l'
	Steak     Tile = 'k'

	// Powers
	Flames     Tile = 'f'
	Storm      Tile = 't'
	Healing    Tile = 'e'
	Teleport   Tile = 'y'
	Skull      Tile = 'q'
	Necromancy Tile = 'ç'

	// Weapons
	Dagger Tile = 'd'
	Bow    Tile = 'a'
	Wand   Tile = 'w'
	Spear  Tile = 's'
	Arrow  Tile = 'r'
)
