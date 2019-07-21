package game

import "math/rand"

func (wg *WorldGenerator) generateOutdoor(levelName string) *Level {
	level := NewLevel(LevelTypeOutdoor)
	level.Name = levelName
	level.InitMaps(WorldHeight, WorldWidth)

	initialMagnitude := 300
	magnitude := initialMagnitude
	for y := 0; y < WorldHeight; y++ {
		for x := 0; x < WorldWidth; x++ {
			t := DirtFloor
			p := rand.Intn(magnitude)
			if p < 50 {
				t = HerbFloor
				level.Map[y][x].MonstersProbability = 6
				magnitude -= initialMagnitude / 50
				if magnitude <= 0 {
					magnitude = initialMagnitude
				}
			}
			level.Map[y][x].T = t
		}
	}

	wg.generateOcean(level)
	wg.generateTrees(level, 20000)
	wg.generateBooks(level, 100)
	wg.generateGrottos(level, 1000)
	wg.generateCities(level, 50)
	objects := []Tile{
		Fruits,
	}
	wg.generateUsables(level, objects, 100)

	wg.g.Levels[levelName] = level
	return level
}

func (wg *WorldGenerator) generateOcean(level *Level) {
	// Ocean top
	for y := 0; y < OceanY-1; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Map[y][x].Object = o
		}
	}
	for x := OceanX - 1; x < WorldWidth-OceanX; x++ {
		o := &Object{Rune: rune(OceanTopSide), Blocking: false}
		o.Pos = Pos{x, OceanY - 1}
		level.Map[o.Y][o.X].Object = o
	}

	// Ocean bottom
	for y := WorldHeight - OceanY; y < WorldHeight; y++ {
		for x := 0; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Map[y][x].Object = o
		}
	}
	for x := OceanX - 1; x < WorldWidth-OceanX; x++ {
		o := &Object{Rune: rune(OceanDownSide), Blocking: false}
		o.Pos = Pos{x, WorldHeight - OceanY}
		level.Map[o.Y][o.X].Object = o
	}

	// Ocean left
	for y := 0; y < WorldHeight; y++ {
		for x := 0; x < OceanX-1; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Map[y][x].Object = o
		}
	}
	for y := OceanY - 1; y < WorldHeight-OceanY; y++ {
		o := &Object{Rune: rune(OceanLeftSide), Blocking: false}
		o.Pos = Pos{OceanX - 1, y}
		level.Map[o.Y][o.X].Object = o
	}

	// Ocean right
	for y := 0; y < WorldHeight; y++ {
		for x := WorldWidth - OceanX; x < WorldWidth; x++ {
			o := &Object{Rune: rune(Ocean), Blocking: true}
			o.Pos = Pos{x, y}
			level.Map[y][x].Object = o
		}
	}
	for y := OceanY - 1; y < WorldHeight-OceanY; y++ {
		o := &Object{Rune: rune(OceanRightSide), Blocking: false}
		o.Pos = Pos{WorldWidth - OceanX, y}
		level.Map[o.Y][o.X].Object = o
	}
}
