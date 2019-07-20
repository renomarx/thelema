package game

import "time"

func (fr *FightingRing) MakeFlame(p Pos, damages int, lifetime int) {
	eff := NewFlame(p, damages)
	eff.TileIdx = 0
	fr.CurrentEffect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond / 4)
	eff.TileIdx = 1
	time.Sleep(time.Duration(lifetime) * time.Millisecond / 4)
	eff.TileIdx = 2
	time.Sleep(time.Duration(lifetime) * time.Millisecond / 4)
	eff.TileIdx = 3
	time.Sleep(time.Duration(lifetime) * time.Millisecond / 4)
	fr.CurrentEffect = nil
}

func (fr *FightingRing) MakeStorm(p Pos, damages int, dir InputType, lifetime int) {
	eff := NewStorm(p, damages, dir)
	fr.CurrentEffect = eff
	time.Sleep(time.Duration(lifetime) * time.Millisecond)
	fr.CurrentEffect = nil
}
