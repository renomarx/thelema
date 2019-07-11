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
