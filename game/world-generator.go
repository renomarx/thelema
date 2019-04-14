package game

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	PlayerInitialX = 3
	PlayerInitialY = 2
)

func (g *Game) GenerateWorld() {
	firstLevel := g.loadLevels()
	g.Level = firstLevel
	g.loadLevelPortals("/maps/world.txt")
	g.loadBooks()
}

func (g *Game) LoadPlayer(p *Player) {
	gameDir := g.GameDir
	p.X = PlayerInitialX
	p.Y = PlayerInitialY
	p.LoadQuests(gameDir)
	p.LoadQuestsObjects(gameDir)
	p.LoadPlayerMenu()
	g.Level.Player = p
}

func (g *Game) loadLevels() *Level {
	g.Levels = make(map[string]*Level)
	firstLevel := g.loadLevelFromFile("level1", LevelTypeGrotto)
	g.loadLevelFromFile("level2", LevelTypeGrotto)
	g.loadLevelFromFile("world", LevelTypeOutdoor)
	return firstLevel
}

func (g *Game) loadLevelPortals(filepath string) {
	g.loadLevelPortalsFromFile(g.GameDir + filepath)
}

func (g *Game) loadLevelFromFile(levelName string, levelType string) *Level {
	dirpath := g.GameDir
	filename := dirpath + "/maps/" + levelName + ".map"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var levelLines []string
	longestRow := 0
	for scanner.Scan() {
		line := scanner.Text()
		levelLines = append(levelLines, line)
		if len(line) > longestRow {
			longestRow = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	level := &Level{}
	level.Type = levelType
	level.Map = make([][]Tile, len(levelLines))
	level.Monsters = make(map[Pos]*Monster)
	level.Objects = make(map[Pos]*Object)
	level.Effects = make(map[Pos]*Effect)
	level.Projectiles = make(map[Pos]*Projectile)
	level.Pnjs = make(map[Pos]*Pnj)
	level.Invocations = make(map[Pos]*Invoked)
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}
	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]

		// Re-compose line to handle utf8
		var utf8line []rune
		for _, c := range line {
			utf8line = append(utf8line, c)
		}

		for x, c := range utf8line {
			var t Tile
			t = DirtFloor
			switch Tile(c) {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case StoneWall:
				t = StoneWall
			case DoorClosed:
				t = DoorClosed
			case DoorOpened:
				t = DoorOpened
			case DirtFloor:
				t = DirtFloor
			case Spider:
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
			case Rat:
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
			case Downstairs:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(Downstairs)}
			case Upstairs:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(Upstairs)}
			case Jason:
				level.Pnjs[Pos{x, y}] = NewPnj(Pos{x, y}, rune(Jason), "jason")
			case Sarah:
				level.Pnjs[Pos{x, y}] = NewPnj(Pos{x, y}, rune(Sarah), "sarah")
			default:
				o := &Object{Rune: c, Blocking: true}
				o.Pos = Pos{x, y}
				level.Objects[Pos{x, y}] = o
			}
			level.Map[y][x] = t
		}
	}

	level.loadPnjsDialogs(dirpath)

	g.Levels[levelName] = level
	return level
}

func (level *Level) loadPnjsDialogs(dirpath string) {
	for _, pnj := range level.Pnjs {
		pnj.LoadDialogs(dirpath)
	}
}

func (game *Game) loadLevelPortalsFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		levels := strings.Split(line, " ")
		levelFrom := strings.Split(levels[0], ",")
		levelTo := strings.Split(levels[1], ",")

		gameLevelFrom := game.Levels[levelFrom[0]]
		if gameLevelFrom == nil {
			panic("Level with name " + levelFrom[0] + " not found")
		}
		gameLevelTo := game.Levels[levelTo[0]]
		if gameLevelTo == nil {
			panic("Level with name " + levelTo[0] + " not found")
		}
		xFrom, _ := strconv.Atoi(levelFrom[1])
		yFrom, _ := strconv.Atoi(levelFrom[2])
		posFrom := Pos{X: xFrom, Y: yFrom}
		xTo, _ := strconv.Atoi(levelTo[1])
		yTo, _ := strconv.Atoi(levelTo[2])
		posTo := Pos{X: xTo, Y: yTo}
		gameLevelFrom.AddPortal(posFrom, &Portal{LevelTo: levelTo[0], PosTo: posTo})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) loadBooks() {
	g.Books = make(map[string]*OBook)
	g.Books["cats"] = g.loadBookFromFile("cats", []string{})
	g.Books["invocat"] = g.loadBookFromFile("invocat", []string{PowerInvocation})
	g.Books["rats"] = g.loadBookFromFile("rats", []string{})
	g.Books["spiders"] = g.loadBookFromFile("spiders", []string{})
}

func (g *Game) loadBookFromFile(filename string, powers []string) *OBook {
	filepath := g.GameDir + "/books/" + filename + ".txt"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	title := ""
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(lines) > 0 {
		title = lines[0]
	}

	return &OBook{Title: title, Text: lines, Powers: powers, Rune: rune(Book)}
}
