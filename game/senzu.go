package game

type USenzu struct {
	Object *Object
}

func (s USenzu) GetObject() *Object {
	return s.Object
}

func (s USenzu) Use(g *Game) {
	p := g.Level.Player
	p.Health.Current = p.Health.Initial
	p.Energy.Current = p.Energy.Initial

	g.GetEventManager().Dispatch(&Event{
		Action:  ActionEat,
		Message: "Health and energy regenerated"})
}
