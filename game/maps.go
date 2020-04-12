package game

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func (g *Game) LoadMapTemplate(mapName, levelName string) *Level {
	dirpath := g.GameDir
	filename := dirpath + "/maps/" + mapName + ".map"
	isDungeon := strings.Contains(filename, "/dungeons/")
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
	level.Name = levelName
	level.InitMaps(len(levelLines), longestRow)
	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		// Re-compose line to handle utf8
		var utf8line []rune
		for _, c := range line {
			utf8line = append(utf8line, c)
		}
		for x, c := range utf8line {
			if isDungeon {
				level.Map[y][x].MonstersProbability = 10
			}
			var t Tile
			t = Floor
			switch Tile(c) {
			case "", " ", "\t", "\n", "\r":
				t = Blank
			case Floor:
				t = Floor
			case MonsterFloor:
				level.Map[y][x].MonstersProbability = 10
				t = MonsterFloor
			case Door:
				level.Map[y][x].Object = &Object{Rune: string(c), Static: true}
			default:
				o := &Object{Rune: string(c), Static: true, Blocking: true}
				o.Pos = Pos{x, y}
				level.Map[y][x].Object = o
			}
			level.Map[y][x].T = t
		}
	}

	return level
}
