package game

func isThereAMonster(level *Level, pos Pos) bool {
	Mux.Lock()
	_, exists := level.Monsters[pos]
	Mux.Unlock()
	return exists
}

func isThereAPnj(level *Level, pos Pos) bool {
	Mux.Lock()
	_, exists := level.Pnjs[pos]
	Mux.Unlock()
	return exists
}

func isThereAnInvocation(level *Level, pos Pos) bool {
	Mux.Lock()
	_, exists := level.Invocations[pos]
	Mux.Unlock()
	return exists
}

func isThereABlockingObject(level *Level, pos Pos) bool {
	if obj, ok := level.Objects[pos]; ok {
		// There is an object !
		return obj.Blocking
	}
	return false
}

func canGo(level *Level, pos Pos) bool {
	if isThereABlockingObject(level, pos) {
		return false
	}
	if isThereAMonster(level, pos) {
		return false
	}
	if isThereAPnj(level, pos) {
		return false
	}
	if pos.Y >= 0 && pos.Y < len(level.Map) {
		if pos.X >= 0 && pos.X < len(level.Map[pos.Y]) {
			return level.Map[pos.Y][pos.X] != StoneWall && level.Map[pos.Y][pos.X] != Blank && level.Map[pos.Y][pos.X] != 0
		}
	}
	return false
}

func openDoor(g *Game, pos Pos) {
	level := g.Level
	t := level.Map[pos.Y][pos.X]
	switch t {
	case DoorClosed:
		level.Map[pos.Y][pos.X] = DoorOpened
		g.GetEventManager().Dispatch(&Event{Action: ActionOpenDoor})
	}
}

func closeDoor(g *Game, pos Pos) {
	level := g.Level
	t := level.Map[pos.Y][pos.X]
	switch t {
	case DoorOpened:
		level.Map[pos.Y][pos.X] = DoorClosed
		g.GetEventManager().Dispatch(&Event{Action: ActionCloseDoor})
	}
}

func (c *Character) moveFromTo(from Pos, to Pos) {
	c.Pos = to
	if from.Y == c.Pos.Y {
		if from.X < c.Pos.X {
			c.LookAt = Right
			c.Xb = CaseLen
			go c.moveRight()
		} else if from.X > c.Pos.X {
			c.LookAt = Left
			c.Xb = -1 * CaseLen
			go c.moveLeft()
		}
	}
	if from.X == c.Pos.X {
		if from.Y < c.Pos.Y {
			c.LookAt = Down
			c.Yb = CaseLen
			go c.moveDown()
		} else if from.Y > c.Pos.Y {
			c.LookAt = Up
			c.Yb = -1 * CaseLen
			go c.moveUp()
		}
	}
}
