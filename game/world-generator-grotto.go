package game

import "math/rand"
import "strconv"

func (wg *WorldGenerator) generateGrottos(level *Level, nbGrottos int) {
	nbTemplates := 2 // TODO load as much templates as there are
	grottoNumber := 0
	for i := 0; i < nbGrottos; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		m := rand.Intn(nbTemplates) + 1
		pos := Pos{X: x, Y: y}

		o := level.Map[pos.Y][pos.X].Object
		if o != nil {
			level.Map[pos.Y][pos.X].Object = &Object{Rune: rune(Downstairs)} // TODO : maybe use a different grotto entry
			mapName := "grotto/grotto" + strconv.Itoa(m)
			wg.generateGrotto(level, pos, mapName, grottoNumber)
			grottoNumber++
		}
	}
}

func (wg *WorldGenerator) generateGrotto(level *Level, pos Pos, mapName string, grottoNumber int) {
	levelName := "Grotto " + strconv.Itoa(grottoNumber)
	nl, npos := wg.LoadMapTemplate(mapName, LevelTypeGrotto, levelName)
	nl.Name = levelName
	wg.g.Levels[levelName] = nl
	level.AddPortal(pos, &Portal{LevelTo: levelName, PosTo: npos})
	nl.AddPortal(npos, &Portal{LevelTo: WorldName, PosTo: pos})

	wg.generateMonsters(nl, BestiaryUnderworld(), rand.Intn(20)+1)
	wg.generateEnnemies(nl, CreaturesUnderworld(), rand.Intn(10)+1)
	objects := []Tile{
		Senzu,
	}
	wg.generateUsables(nl, objects, rand.Intn(5)+1)
	wg.generateBooks(nl, rand.Intn(2))
}
