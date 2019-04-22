package game

import "math/rand"

func (wg *WorldGenerator) generateOutdoor(levelName string) *Level {
	level := NewLevel(LevelTypeOutdoor)
	level.Map = make([][]Tile, WorldHeight)
	for i := range level.Map {
		level.Map[i] = make([]Tile, WorldWidth)
	}

	for y := 0; y < WorldHeight; y++ {
		for x := 0; x < WorldWidth; x++ {
			t := DirtFloor
			level.Map[y][x] = t
		}
	}

	wg.generateOcean(level)
	wg.generateTrees(level, 10000)
	wg.generateMonsters(level, 1000)
	wg.generateGrottos(level)
	wg.generateCities(level)

	wg.g.Levels[levelName] = level
	return level
}

func (wg *WorldGenerator) generateOcean(level *Level) {
	// Ocean top
	for y := 0; y < OceanY-1; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for x := OceanX - 1; x < WorldWidth-OceanX; x++ {
		o := &Object{Rune: rune(OceanTopSide), Blocking: false}
		o.Pos = Pos{x, OceanY - 1}
		level.Objects[Pos{x, OceanY - 1}] = o
	}

	// Ocean bottom
	for y := WorldHeight - OceanY; y < WorldHeight; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for x := OceanX - 1; x < WorldWidth-OceanX; x++ {
		o := &Object{Rune: rune(OceanDownSide), Blocking: false}
		o.Pos = Pos{x, WorldHeight - OceanY}
		level.Objects[Pos{x, WorldHeight - OceanY}] = o
	}

	// Ocean left
	for y := 0; y < WorldHeight; y++ {
		for x := 0; x < OceanX-1; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for y := OceanY - 1; y < WorldHeight-OceanY; y++ {
		o := &Object{Rune: rune(OceanLeftSide), Blocking: false}
		o.Pos = Pos{OceanX - 1, y}
		level.Objects[Pos{OceanX - 1, y}] = o
	}

	// Ocean right
	for y := 0; y < WorldHeight; y++ {
		for x := WorldWidth - OceanX; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for y := OceanY - 1; y < WorldHeight-OceanY; y++ {
		o := &Object{Rune: rune(OceanRightSide), Blocking: false}
		o.Pos = Pos{WorldWidth - OceanX, y}
		level.Objects[Pos{WorldWidth - OceanX, y}] = o
	}
}

func (wg *WorldGenerator) generateTrees(level *Level, nbTrees int) {
	for i := 0; i < nbTrees; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		pos := Pos{X: x, Y: y}
		o := &Object{Rune: rune(Tree), Blocking: true}
		o.Pos = pos
		_, oe := level.Objects[pos]
		if !oe {
			level.Objects[pos] = o
		}
	}
}

func (wg *WorldGenerator) generateMonsters(level *Level, nbMonsters int) {
	bestiary := Bestiary()
	for i := 0; i < nbMonsters; i++ {
		x := rand.Intn(len(level.Map[0]))
		y := rand.Intn(len(level.Map))
		m := rand.Intn(len(bestiary))
		pos := Pos{X: x, Y: y}

		mt := bestiary[m]
		_, oe := level.Objects[pos]
		if !oe {
			level.Monsters[pos] = NewMonster(mt, pos)
		}
	}
}
