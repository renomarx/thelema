package game

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func (wg *WorldGenerator) LoadMapTemplate(mapName string, levelType string, levelName string) (*Level, Pos) {
	dirpath := wg.g.GameDir
	filename := dirpath + "/maps/" + mapName + ".map"
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
	houseNumber := 0
	nbHouseTemplates := 1 // TODO load as much templates as there are
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
			if levelType == LevelTypeHouse {
				t = GreenFloor
			}
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
			case GreenFloor:
				t = GreenFloor
			case Upstairs:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(Upstairs)}
				initialPos = Pos{X: x, Y: y}
			case CityOut:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(CityOut)}
				initialPos = Pos{X: x, Y: y}
			case HouseDoor:
				level.Objects[Pos{x, y}] = &Object{Rune: rune(HouseDoor)}
				if levelType == LevelTypeHouse {
					initialPos = Pos{X: x, Y: y}
				} else {
					m := rand.Intn(nbHouseTemplates) + 1
					mapName := "house/house" + strconv.Itoa(m)
					wg.generateHouse(level, Pos{X: x, Y: y}, mapName, houseNumber, levelName)
					houseNumber++
				}
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
