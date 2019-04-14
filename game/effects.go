package game

import "time"

type Effect struct {
	Object
	Size int
}

func (g *Game) MakeExplosion(p Pos, size int, lifetime int) {
	g.GetEventManager().Dispatch(&Event{
		Action: ActionExplode})
	level := g.Level
	eff := &Effect{}
	eff.Rune = rune(Explosion)
	eff.Blocking = false
	eff.Size = size

	Mux.Lock()
	level.Effects[p] = eff
	Mux.Unlock()
	go func() {
		time.Sleep(time.Duration(lifetime) * time.Millisecond)
		Mux.Lock()
		delete(level.Effects, p)
		Mux.Unlock()
	}()
}
