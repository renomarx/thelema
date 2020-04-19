package game

func isThereAPnj(level *Level, pos Pos) bool {
	return level.GetPnj(pos) != nil
}

func isThereAPlayerCharacter(level *Level, pos Pos) bool {
	p := level.Player
	if p != nil && p.X == pos.X && p.Y == pos.Y && p.Z == pos.Z {
		return true
	}
	return false
}

func isThereABlockingObject(level *Level, pos Pos) bool {
	o := level.GetObject(pos)
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
	if pos.Z >= 0 && pos.Z < len(level.Map) {
		if pos.Y >= 0 && pos.Y < len(level.Map[pos.Z]) {
			if pos.X >= 0 && pos.X < len(level.Map[pos.Z][pos.Y]) {
				return level.Map[pos.Z][pos.Y][pos.X].T != Blank && level.Map[pos.Z][pos.Y][pos.X].T != ""
			}
		}
	}
	return false
}

func isThereAPortalAround(level *Level, pos Pos) bool {
	if level.Map[pos.Z][pos.Y][pos.X].Portal != nil {
		return true
	}
	if isInsideMap(level, Pos{X: pos.X - 1, Y: pos.Y, Z: pos.Z}) {
		if level.Map[pos.Z][pos.Y][pos.X-1].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X + 1, Y: pos.Y, Z: pos.Z}) {
		if level.Map[pos.Z][pos.Y][pos.X+1].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X, Y: pos.Y - 1, Z: pos.Z}) {
		if level.Map[pos.Z][pos.Y-1][pos.X].Portal != nil {
			return true
		}
	}
	if isInsideMap(level, Pos{X: pos.X, Y: pos.Y + 1, Z: pos.Z}) {
		if level.Map[pos.Z][pos.Y+1][pos.X].Portal != nil {
			return true
		}
	}
	return false
}
