package game

import "time"

type Effect struct {
	Object
	Size int
}

func (level *Level) MakeExplosion(p Pos, size int, lifetime int) {
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
