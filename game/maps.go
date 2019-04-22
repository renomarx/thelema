package game

import (
	"bufio"
	"log"
	"os"
)

func (g *Game) LoadMapTemplate(levelName string, levelType string) (*Level, Pos) {
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

	level := NewLevel(levelType)
	level.Map = make([][]Tile, len(levelLines))
	initialPos := Pos{X: 1, Y: 1}
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
			case Upstairs:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(Upstairs)} // TODO : maybe use a different grotto out
				initialPos = Pos{X: x, Y: y}
			case CityOut:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(CityOut)} // TODO : maybe use a different city out
				initialPos = Pos{X: x, Y: y}
			default:
				o := &Object{Rune: c, Blocking: true}
				o.Pos = Pos{x, y}
				level.Objects[Pos{x, y}] = o
			}
			level.Map[y][x] = t
		}
	}

	return level, initialPos
}
