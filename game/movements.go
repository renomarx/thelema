package game

func isThereAMonster(level *Level, pos Pos) bool {
	_, exists := level.Monsters[pos]
	return exists
}

func isThereAnEnemy(level *Level, pos Pos) bool {
	_, exists := level.Enemies[pos]
	return exists
}

func isThereAPnj(level *Level, pos Pos) bool {
	_, exists := level.Pnjs[pos]
	return exists
}

func isThereAnInvocation(level *Level, pos Pos) bool {
	_, exists := level.Invocations[pos]
	return exists
}

func isThereAFriend(level *Level, pos Pos) bool {
	_, exists := level.Friends[pos]
	return exists
}

func isThereAPlayerCharacter(level *Level, pos Pos) bool {
	p := level.Player
	if p != nil && p.X == pos.X && p.Y == pos.Y {
		return true
	}
	if _, exists := level.Friends[pos]; exists {
		return true
	}
	if _, exists := level.Invocations[pos]; exists {
		return true
	}
	return false
}

func isThereAnEnemyCharacter(level *Level, pos Pos) bool {
	if _, exists := level.Monsters[pos]; exists {
		return true
	}
	if _, exists := level.Enemies[pos]; exists {
		return true
	}
	return false
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
	if isThereAnEnemy(level, pos) {
		return false
	}
	if isThereAPnj(level, pos) {
		return false
	}
	if pos.Y >= 0 && pos.Y < len(level.Map) {
		if pos.X >= 0 && pos.X < len(level.Map[pos.Y]) {
			return level.Map[pos.Y][pos.X] != Blank && level.Map[pos.Y][pos.X] != 0
		}
	}
	return false
}

func openDoor(g *Game, pos Pos) {
	level := g.Level
	o, e := level.Objects[pos]
	if e {
		switch Tile(o.Rune) {
		case DoorClosed:
			o.Rune = rune(DoorOpened)
			o.Blocking = false
			g.GetEventManager().Dispatch(&Event{Action: ActionOpenDoor})
		}
	}
}

func closeDoor(g *Game, pos Pos) {
	level := g.Level
	o, e := level.Objects[pos]
	if e {
		switch Tile(o.Rune) {
		case DoorOpened:
			o.Rune = rune(DoorClosed)
			o.Blocking = true
			g.GetEventManager().Dispatch(&Event{Action: ActionCloseDoor})
		}
	}
}
