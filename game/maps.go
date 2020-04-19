package game

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func (g *Game) LoadMapTemplate(mapName, levelName string) *Level {
	dirpath := g.GameDir

	level := NewLevel()
	level.Name = levelName

	z := 0
	filename := dirpath + "/maps/" + mapName + ".map"
	for FileExists(filename) {
		g.doLoadMapTemplate(filename, z, level)
		z++
		filename = dirpath + "/maps/" + mapName + strconv.Itoa(z) + ".map"
	}

	return level
}

func (g *Game) doLoadMapTemplate(filename string, z int, level *Level) {
	isDungeon := strings.Contains(filename, "/dungeons/")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var m [][]Case

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		// Re-compose line to handle utf8
		var utf8line []rune
		for _, c := range line {
			utf8line = append(utf8line, c)
		}
		var row []Case
		for x, c := range utf8line {
			var ca Case
			if isDungeon {
				ca.MonstersProbability = 10
			}
			var t Tile
			t = Floor
			switch Tile(c) {
			case "", " ", "\t", "\n", "\r":
				t = Blank
			case Floor:
				t = Floor
			case MonsterFloor:
				ca.MonstersProbability = 10
				t = MonsterFloor
			case Door:
				ca.Object = &Object{Rune: string(c), Static: true}
			default:
				o := &Object{Rune: string(c), Static: true, Blocking: true}
				o.Pos = Pos{X: x, Y: y, Z: z}
				ca.Object = o
			}
			ca.T = t
			row = append(row, ca)
		}
		m = append(m, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	level.Map = append(level.Map, m)

}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
