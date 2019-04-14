package game

type USenzu struct {
	Object *Object
}

func (s USenzu) GetObject() *Object {
	return s.Object
}

func (s USenzu) Use(g *Game) {
	p := g.Level.Player
	p.Hitpoints.Current = p.Hitpoints.Initial
	p.Energy.Current = p.Energy.Initial

	g.GetEventManager().Dispatch(&Event{
		Type:    PlayerEventsType,
		Action:  ActionEat,
		Message: "Health and energy regenerated"})
}
