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
	wg.generateTrees(level)
	wg.generateMonsters(level)

	wg.g.Levels[levelName] = level
	return level
}

func (wg *WorldGenerator) generateOcean(level *Level) {
	for y := 0; y < 9; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for x := 9; x < WorldWidth-10; x++ {
		o := &Object{Rune: rune(OceanTopSide), Blocking: true}
		o.Pos = Pos{x, 9}
		level.Objects[Pos{x, 9}] = o
	}

	for y := WorldHeight - 9; y < WorldHeight; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for x := 9; x < WorldWidth-10; x++ {
		o := &Object{Rune: rune(OceanDownSide), Blocking: true}
		o.Pos = Pos{x, WorldHeight - 10}
		level.Objects[Pos{x, WorldHeight - 10}] = o
	}

	for y := 0; y < WorldHeight; y++ {
		for x := 0; x < 9; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for y := 9; y < WorldHeight-10; y++ {
		o := &Object{Rune: rune(OceanLeftSide), Blocking: true}
		o.Pos = Pos{9, y}
		level.Objects[Pos{9, y}] = o
	}

	for y := 0; y < WorldHeight; y++ {
		for x := WorldWidth - 9; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Objects[Pos{x, y}] = o
		}
	}
	for y := 9; y < WorldHeight-10; y++ {
		o := &Object{Rune: rune(OceanRightSide), Blocking: true}
		o.Pos = Pos{WorldWidth - 10, y}
		level.Objects[Pos{WorldWidth - 10, y}] = o
	}
}

func (wg *WorldGenerator) generateTrees(level *Level) {
	nbTrees := 5000
	for i := 0; i < nbTrees; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		pos := Pos{X: x, Y: y}
		o := &Object{Rune: rune(Tree), Blocking: true}
		o.Pos = pos
		level.Objects[pos] = o
	}
}

func (wg *WorldGenerator) generateMonsters(level *Level) {
	bestiary := Bestiary()
	nbMonsters := 500
	for i := 0; i < nbMonsters; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		m := rand.Intn(len(bestiary))
		pos := Pos{X: x, Y: y}

		mt := bestiary[m]
		level.Monsters[pos] = NewMonster(mt, pos)
	}
}
