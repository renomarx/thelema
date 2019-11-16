package game

func isThereAPnj(level *Level, pos Pos) bool {
	return level.GetPnj(pos.X, pos.Y) != nil
}

func isThereAPlayerCharacter(level *Level, pos Pos) bool {
	p := level.Player
	if p != nil && p.X == pos.X && p.Y == pos.Y {
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
	if isThereAPnj(level, pos) {
		return false
	}
	return true
}

func isInsideMap(level *Level, pos Pos) bool {
	if pos.Y >= 0 && pos.Y < len(level.Map) {
		if pos.X >= 0 && pos.X < len(level.Map[pos.Y]) {
			return level.Map[pos.Y][pos.X].T != Blank && level.Map[pos.Y][pos.X].T != ""
		}
	}
	return false
}

func isThereAPortalAround(level *Level, pos Pos) bool {
	if level.Map[pos.Y][pos.X].Portal != nil {
		return true
	}
	if isInsideMap(level, Pos{X: pos.X - 1, Y: pos.Y}) {
		if level.Map[pos.Y][pos.X-1].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X + 1, Y: pos.Y}) {
		if level.Map[pos.Y][pos.X+1].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X, Y: pos.Y - 1}) {
		if level.Map[pos.Y-1][pos.X].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X, Y: pos.Y + 1}) {
		if level.Map[pos.Y+1][pos.X].Portal != nil {
			return true
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
			o.Rune = string(DoorOpened)
			o.Blocking = false
			EM.Dispatch(&Event{Action: ActionOpenDoor})
		}
	}
}

func closeDoor(g *Game, pos Pos) {
	level := g.Level
	o := level.GetObject(pos.X, pos.Y)
	if o != nil {
		switch Tile(o.Rune) {
		case DoorOpened:
			o.Rune = string(DoorClosed)
			o.Blocking = true
			EM.Dispatch(&Event{Action: ActionCloseDoor})
		}
	}
}
