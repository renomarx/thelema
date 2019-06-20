package game

func isThereAMonster(level *Level, pos Pos) bool {
	return level.GetMonster(pos.X, pos.Y) != nil
}

func isThereAnEnemy(level *Level, pos Pos) bool {
	return level.GetEnemy(pos.X, pos.Y) != nil
}

func isThereAPnj(level *Level, pos Pos) bool {
	return level.GetPnj(pos.X, pos.Y) != nil
}

func isThereAnInvocation(level *Level, pos Pos) bool {
	return level.GetInvocation(pos.X, pos.Y) != nil
}

func isThereAFriend(level *Level, pos Pos) bool {
	return level.GetFriend(pos.X, pos.Y) != nil
}

func isThereAPlayerCharacter(level *Level, pos Pos) bool {
	p := level.Player
	if p != nil && p.X == pos.X && p.Y == pos.Y {
		return true
	}
	if level.GetFriend(pos.X, pos.Y) != nil {
		return true
	}
	if level.GetInvocation(pos.X, pos.Y) != nil {
		return true
	}
	return false
}

func isThereAnEnemyCharacter(level *Level, pos Pos) bool {
	if level.GetMonster(pos.X, pos.Y) != nil {
		return true
	}
	if level.GetEnemy(pos.X, pos.Y) != nil {
		return true
	}
	return false
}

func isThereABlockingObject(level *Level, pos Pos) bool {
	o := level.GetObject(pos.X, pos.Y)
	if o == nil {
		return false
	}
	return o.Blocking
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
			return level.Map[pos.Y][pos.X].T != Blank && level.Map[pos.Y][pos.X].T != 0
		}
	}
	return false
}

func openDoor(g *Game, pos Pos) {
	level := g.Level
	o := level.GetObject(pos.X, pos.Y)
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
	o := level.GetObject(pos.X, pos.Y)
	if o != nil {
		switch Tile(o.Rune) {
		case DoorOpened:
			o.Rune = rune(DoorClosed)
			o.Blocking = true
			g.GetEventManager().Dispatch(&Event{Action: ActionCloseDoor})
		}
	}
}
