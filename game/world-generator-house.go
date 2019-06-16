package game

import "math/rand"
import "strconv"

func (wg *WorldGenerator) generateHouse(level *Level, pos Pos, mapName string, houseNumber int, cityName string) {
	levelName := cityName + " - House " + strconv.Itoa(houseNumber)
	nl, npos := wg.LoadMapTemplate(mapName, LevelTypeHouse, levelName)
	nl.Name = levelName
	wg.g.Levels[levelName] = nl
	level.AddPortal(pos, &Portal{LevelTo: levelName, PosTo: npos})
	nl.AddPortal(npos, &Portal{LevelTo: cityName, PosTo: pos})

	if levelName != FirstLevelName {
		wg.generatePnjs(nl, rand.Intn(2))
		objects := []Tile{
			Bread,
			Water,
		}
		wg.generateUsables(nl, objects, rand.Intn(3))
	}
}
