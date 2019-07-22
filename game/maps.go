package game

import (
	"bufio"
	"log"
	"os"
)

func (g *Game) LoadMapTemplate(mapName string, levelName string) *Level {
	dirpath := g.GameDir
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

	level := NewLevel()
	level.InitMaps(len(levelLines), longestRow)
	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		// Re-compose line to handle utf8
		var utf8line []rune
		for _, c := range line {
			utf8line = append(utf8line, c)
		}
		var t Tile
		t = DirtFloor
		for x, c := range utf8line {
			switch Tile(c) {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case DirtFloor:
				t = DirtFloor
			case CityFloor:
				t = CityFloor
			case HerbFloor:
				level.Map[y][x].MonstersProbability = 6
				t = HerbFloor
			case DoorOpened:
				level.Map[y][x].Object = &Object{Rune: rune(DoorOpened)}
			case Upstairs:
				level.Map[y][x].Object = &Object{Rune: rune(Upstairs)}
			case Downstairs:
				level.Map[y][x].Object = &Object{Rune: rune(Downstairs)}
			case CityEntry:
				level.Map[y][x].Object = &Object{Rune: rune(CityEntry)}
			case CityOut:
				level.Map[y][x].Object = &Object{Rune: rune(CityOut)}
			case HouseDoor:
				level.Map[y][x].Object = &Object{Rune: rune(HouseDoor)}
			case PrisonDoor:
				level.Map[y][x].Object = &Object{Rune: rune(PrisonDoor)}
			default:
				o := &Object{Rune: c, Blocking: true}
				o.Pos = Pos{x, y}
				level.Map[y][x].Object = o
			}
			level.Map[y][x].T = t
		}
	}

	return level
}
