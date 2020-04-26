package game

type Tile string

const (
	Blank        Tile = " "
	Wall         Tile = "#"
	Door         Tile = "="
	Floor        Tile = "-"
	MonsterFloor Tile = "*"
	Ocean             = "o"

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
	Necromancy Tile = "necromancy"
	Calm       Tile = "calm"

	// Icons
	BrutalStrengthIcon Tile = "brutal_strength_icon"
	QuickeningIcon     Tile = "quickening_icon"
	StormIcon          Tile = "storm_icon"
	HealingIcon        Tile = "healing_icon"
	InvocationIcon     Tile = "invocation_icon"
	NecromancyIcon     Tile = "necromancy_icon"
	CalmIcon           Tile = "calm_icon"
	FlamesIcon         Tile = "flames_icon"
)
