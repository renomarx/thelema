package game

func isThereAMonster(level *Level, pos Pos) bool {
	return level.Monsters[pos.Y][pos.X] != nil
}

func isThereAnEnemy(level *Level, pos Pos) bool {
	return level.Enemies[pos.Y][pos.X] != nil
}

func isThereAPnj(level *Level, pos Pos) bool {
	return level.Pnjs[pos.Y][pos.X] != nil
}

func isThereAnInvocation(level *Level, pos Pos) bool {
	return level.Invocations[pos.Y][pos.X] != nil
}

func isThereAFriend(level *Level, pos Pos) bool {
	return level.Friends[pos.Y][pos.X] != nil
}

func isThereAPlayerCharacter(level *Level, pos Pos) bool {
	p := level.Player
	if p != nil && p.X == pos.X && p.Y == pos.Y {
		return true
	}
	if level.Friends[pos.Y][pos.X] != nil {
		return true
	}
	if level.Invocations[pos.Y][pos.X] != nil {
		return true
	}
	return false
}

func isThereAnEnemyCharacter(level *Level, pos Pos) bool {
	if level.Monsters[pos.Y][pos.X] != nil {
		return true
	}
	if level.Enemies[pos.Y][pos.X] != nil {
		return true
	}
	return false
}

func isThereABlockingObject(level *Level, pos Pos) bool {
	if level.Objects[pos.Y][pos.X] == nil {
		return false
	}
	return level.Objects[pos.Y][pos.X].Blocking
}

func canGo(level *Level, pos Pos) bool {
	if !isInsideMap(level, pos) {
		return false
	}
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
	return true
}

func isInsideMap(level *Level, pos Pos) bool {
	if pos.Y >= 0 && pos.Y < len(level.Map) {
		if pos.X >= 0 && pos.X < len(level.Map[pos.Y]) {
			return level.Map[pos.Y][pos.X] != Blank && level.Map[pos.Y][pos.X] != 0
		}
	}
	return true
}

func openDoor(g *Game, pos Pos) {
	level := g.Level
	o := level.Objects[pos.Y][pos.X]
	if o != nil {
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
	o := level.Objects[pos.Y][pos.X]
	if o != nil {
		switch Tile(o.Rune) {
		case DoorOpened:
			o.Rune = rune(DoorClosed)
			o.Blocking = true
			g.GetEventManager().Dispatch(&Event{Action: ActionCloseDoor})
		}
	}
}
