package game

import (
	"bufio"
	"log"
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
	level.InitMaps(len(levelLines), longestRow)
	initialPos := Pos{X: 1, Y: 1}
	houseTemplates := LoadFilenames(wg.g.GameDir + "/maps/house")
	nbHouseTemplates := len(houseTemplates)
	houseNumber := 0
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
			switch levelType {
			case LevelTypeHouse, LevelTypeCity:
				t = GreenFloor
			case LevelTypeGrotto:
				t = DirtFloor
				level.Map[y][x].MonstersProbability = 10
			}
			switch Tile(c) {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case DirtFloor, GreenFloor:
			case DoorOpened:
				level.Map[y][x].Object = &Object{Rune: rune(DoorOpened)}
			case Upstairs:
				level.Map[y][x].Object = &Object{Rune: rune(Upstairs)}
				initialPos = Pos{X: x, Y: y}
			case CityOut:
				level.Map[y][x].Object = &Object{Rune: rune(CityOut)}
				initialPos = Pos{X: x, Y: y}
			case HouseDoor:
				level.Map[y][x].Object = &Object{Rune: rune(HouseDoor)}
				if levelType == LevelTypeHouse {
					initialPos = Pos{X: x, Y: y}
				} else {
					m := (houseNumber % nbHouseTemplates) + 1
					mapName := "house/house" + strconv.Itoa(m)
					wg.generateHouse(level, Pos{X: x, Y: y}, mapName, houseNumber, levelName)
					houseNumber++
				}
			case PrisonDoor:
				level.Map[y][x].Object = &Object{Rune: rune(PrisonDoor)}
				if levelType == LevelTypeHouse {
					initialPos = Pos{X: x, Y: y}
				} else {
					mapName := "place/prison"
					wg.generatePrison(level, Pos{X: x, Y: y}, mapName, houseNumber, levelName)
				}
			default:
				o := &Object{Rune: c, Blocking: true}
				o.Pos = Pos{x, y}
				level.Map[y][x].Object = o
			}
			level.Map[y][x].T = t
		}
	}

	return level, initialPos
}
