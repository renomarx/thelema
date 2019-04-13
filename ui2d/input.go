package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) LoadKeymap() {
	ui.Keymap = make(map[string]sdl.Keycode)
	ui.Keymap["a"] = sdl.K_a
	ui.Keymap["z"] = sdl.K_z
	ui.Keymap["e"] = sdl.K_e
	ui.Keymap["r"] = sdl.K_r
	ui.Keymap["t"] = sdl.K_t
	ui.Keymap["y"] = sdl.K_y
	ui.Keymap["u"] = sdl.K_u
	ui.Keymap["i"] = sdl.K_i
	ui.Keymap["o"] = sdl.K_o
	ui.Keymap["p"] = sdl.K_p
	ui.Keymap["q"] = sdl.K_q
	ui.Keymap["s"] = sdl.K_s
	ui.Keymap["d"] = sdl.K_d
	ui.Keymap["f"] = sdl.K_f
	ui.Keymap["g"] = sdl.K_g
	ui.Keymap["h"] = sdl.K_h
	ui.Keymap["j"] = sdl.K_j
	ui.Keymap["k"] = sdl.K_k
	ui.Keymap["l"] = sdl.K_l
	ui.Keymap["m"] = sdl.K_m
	ui.Keymap["w"] = sdl.K_w
	ui.Keymap["x"] = sdl.K_x
	ui.Keymap["c"] = sdl.K_c
	ui.Keymap["v"] = sdl.K_v
	ui.Keymap["b"] = sdl.K_b
	ui.Keymap["n"] = sdl.K_n

	ui.Keymap["up"] = sdl.K_UP
	ui.Keymap["down"] = sdl.K_DOWN
	ui.Keymap["left"] = sdl.K_LEFT
	ui.Keymap["right"] = sdl.K_RIGHT

	ui.Keymap["tab"] = sdl.K_TAB
	ui.Keymap["escape"] = sdl.K_ESCAPE
	ui.Keymap["space"] = sdl.K_SPACE
	ui.Keymap["l-alt"] = sdl.K_LALT
	ui.Keymap["r-alt"] = sdl.K_RALT
	ui.Keymap["l-ctrl"] = sdl.K_LCTRL
	ui.Keymap["r-ctrl"] = sdl.K_RCTRL
	ui.Keymap["l-shift"] = sdl.K_LSHIFT
	ui.Keymap["r-shift"] = sdl.K_RSHIFT
}

func (ui *UI) GetInput() {
	input := ui.Game.GetInput()
	input2 := ui.Game.GetInput2()
	event := sdl.PollEvent()
	if event != nil {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			input.Typ = game.Quit
		case *sdl.KeyboardEvent:

			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case ui.Keymap[ui.Game.Config.Keymap.Up]:
					input.Typ = game.Up
				case ui.Keymap[ui.Game.Config.Keymap.Down]:
					input.Typ = game.Down
				case ui.Keymap[ui.Game.Config.Keymap.Left]:
					input.Typ = game.Left
				case ui.Keymap[ui.Game.Config.Keymap.Right]:
					input.Typ = game.Right
				case ui.Keymap[ui.Game.Config.Keymap.Action]:
					input.Typ = game.Action
				case ui.Keymap[ui.Game.Config.Keymap.Power]:
					input.Typ = game.Power
				case ui.Keymap[ui.Game.Config.Keymap.Escape]:
					input.Typ = game.Escape
				case ui.Keymap[ui.Game.Config.Keymap.Select]:
					input.Typ = game.Select
				case ui.Keymap[ui.Game.Config.Keymap.Speed]:
					input2.Typ = game.SpeedUp
				}
				ui.LastKeyDown = e.Keysym.Sym
			}

			if e.Type == sdl.KEYUP {
				if e.Keysym.Sym == ui.Keymap[ui.Game.Config.Keymap.Speed] {
					input2.Typ = game.None
				} else {
					if ui.LastKeyDown == e.Keysym.Sym {
						input.Typ = game.StayStill
					}
				}
			}
		default:
		}
	}
}

func (ui *UI) WatchInput() {
	for {
		ui.GetInput()
	}
}
