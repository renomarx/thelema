package game

import (
	"math"
	"sort"
)

type PriorityPos struct {
	Pos
	Priority int
}

type PriorityArray []PriorityPos

func (p PriorityArray) Len() int {
	return len(p)
}
func (p PriorityArray) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PriorityArray) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}

func (level *Level) getNeighbors(pos Pos, v Viewer) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{X: pos.X - 1, Y: pos.Y}
	if v.CanSee(level, left) {
		neighbors = append(neighbors, left)
	}
	right := Pos{X: pos.X + 1, Y: pos.Y}
	if v.CanSee(level, right) {
		neighbors = append(neighbors, right)
	}
	up := Pos{X: pos.X, Y: pos.Y - 1}
	if v.CanSee(level, up) {
		neighbors = append(neighbors, up)
	}
	down := Pos{X: pos.X, Y: pos.Y + 1}
	if v.CanSee(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

type Viewer interface {
	CanSee(level *Level, pos Pos) bool
}

func (level *Level) astar(start Pos, goal Pos, v Viewer) []Pos {
	var result []Pos
	frontier := make(PriorityArray, 0, 8)
	frontier = append(frontier, PriorityPos{start, 1})
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	for len(frontier) > 0 {
		sort.Stable(frontier)
		current := frontier[0]
		if current.Pos == goal {
			p := current.Pos
			result = append([]Pos{p}, result...)
			for p != start {
				p = cameFrom[p]
				result = append([]Pos{p}, result...)
			}
			break
		}
		frontier = frontier[1:]
		for _, next := range level.getNeighbors(current.Pos, v) {
			newCost := costSoFar[current.Pos] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = append(frontier, PriorityPos{next, priority})
				cameFrom[next] = current.Pos
			}
		}
	}
	return result
}
