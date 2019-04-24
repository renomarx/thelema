package game

import "math/rand"
import "strconv"

func (wg *WorldGenerator) generateCities(level *Level) {
	nbCities := 100
	nbTemplates := 1 // TODO load as much templates as there are
	cityNumber := 0
	for i := 0; i < nbCities; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		m := rand.Intn(nbTemplates) + 1
		pos := Pos{X: x, Y: y}

		_, oe := level.Objects[pos]
		if !oe {
			level.Objects[pos] = &Object{Rune: rune(CityEntry)} // TODO : maybe use a different city entry
			cityName := "city/city" + strconv.Itoa(m)
			wg.generateCity(level, pos, cityName, cityNumber)
			cityNumber++
		}
	}
}

func (wg *WorldGenerator) generateCity(level *Level, pos Pos, cityName string, cityNumber int) {
	nl, npos := wg.g.LoadMapTemplate(cityName, LevelTypeCity)
	levelName := "grotto" + strconv.Itoa(cityNumber)
	wg.g.Levels[levelName] = nl
	level.AddPortal(pos, &Portal{LevelTo: levelName, PosTo: npos})
	nl.AddPortal(npos, &Portal{LevelTo: WorldName, PosTo: pos})

	wg.generatePnjs(nl, rand.Intn(20)+1)
}
