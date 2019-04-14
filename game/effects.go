package game

import "time"

type Effect struct {
	Object
	Size int
}

const ExplosionSizeSmall = "SMALL"
const ExplosionSizeMedium = "MEDIUM"
const ExplosionSizeLarge = "LARGE"

func (g *Game) MakeExplosion(p Pos, size int, lifetime int) {
	esize := ExplosionSizeSmall
	if size >= 50 {
		esize = ExplosionSizeMedium
		if size >= 100 {
			esize = ExplosionSizeLarge
		}
	}
	g.GetEventManager().Dispatch(&Event{Action: ActionExplode, Payload: map[string]string{"size": esize}})
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
