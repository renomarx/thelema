package game

import "math/rand"
import "strconv"

func (wg *WorldGenerator) generateCities(level *Level, nbCities int) {
	nbTemplates := 1 // TODO load as much templates as there are
	cityNumber := 0
	for i := 0; i < nbCities; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		m := rand.Intn(nbTemplates) + 1
		pos := Pos{X: x, Y: y}

		o := level.Objects[pos.Y][pos.X]
		if o == nil {
			level.Objects[pos.Y][pos.X] = &Object{Rune: rune(CityEntry)} // TODO : maybe use a different city entry
			mapName := "city/city" + strconv.Itoa(m)
			wg.generateCity(level, pos, mapName, cityNumber)
			cityNumber++
		}
	}
}

func (wg *WorldGenerator) generateCity(level *Level, pos Pos, mapName string, cityNumber int) {
	levelName := "City " + strconv.Itoa(cityNumber)
	nl, npos := wg.LoadMapTemplate(mapName, LevelTypeCity, levelName)
	nl.Name = levelName
	wg.g.Levels[levelName] = nl
	level.AddPortal(pos, &Portal{LevelTo: levelName, PosTo: npos})
	nl.AddPortal(npos, &Portal{LevelTo: WorldName, PosTo: pos})

	wg.generatePnjs(nl, rand.Intn(20)+1)
	objects := []Tile{
		Bread,
		Water,
	}
	wg.generateUsables(nl, objects, rand.Intn(5)+1)
}
