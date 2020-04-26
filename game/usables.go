package game

import "fmt"

type Usable interface {
	Use(g *Game)
	GetObject() *Object
}

func NewUsable(o *Object) Usable {
	switch Tile(o.Rune) {
	case Senzu:
		return Food{Health: 100, Energy: 100, Object: o}
	case Fruits:
		return Food{Health: 20, Energy: 0, Object: o}
	case Bread:
		return Food{Health: 30, Energy: 0, Object: o}
	case Water:
		return Food{Health: 0, Energy: 100, Object: o}
	case Steak:
		return Food{Health: 50, Energy: 0, Object: o}
	default:
		return nil
	}
}

type Food struct {
	Object *Object
	Health int
	Energy int
}

func (f Food) GetObject() *Object {
	return f.Object
}

func (f Food) Use(g *Game) {
	p := g.Level.Player
	p.Health.Restore(f.Health)
	p.Energy.Restore(f.Energy)
	EM.Dispatch(&Event{
		Action:  ActionEat,
		Message: fmt.Sprintf("Santé régénerée de %d", f.Health)})
}
