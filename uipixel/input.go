package uipixel

import (
	"thelema/game"

	"github.com/faiface/pixel/pixelgl"
)

func (ui *UI) LoadKeymap() {
	ui.Keymap = make(map[string]pixelgl.Button)
	ui.Keymap["a"] = pixelgl.KeyA
	ui.Keymap["z"] = pixelgl.KeyZ
	ui.Keymap["e"] = pixelgl.KeyE
	ui.Keymap["r"] = pixelgl.KeyR
	ui.Keymap["t"] = pixelgl.KeyT
	ui.Keymap["y"] = pixelgl.KeyY
	ui.Keymap["u"] = pixelgl.KeyU
	ui.Keymap["i"] = pixelgl.KeyI
	ui.Keymap["o"] = pixelgl.KeyO
	ui.Keymap["p"] = pixelgl.KeyP
	ui.Keymap["q"] = pixelgl.KeyQ
	ui.Keymap["s"] = pixelgl.KeyS
	ui.Keymap["d"] = pixelgl.KeyD
	ui.Keymap["f"] = pixelgl.KeyF
	ui.Keymap["g"] = pixelgl.KeyG
	ui.Keymap["h"] = pixelgl.KeyH
	ui.Keymap["j"] = pixelgl.KeyJ
	ui.Keymap["k"] = pixelgl.KeyK
	ui.Keymap["l"] = pixelgl.KeyL
	ui.Keymap["m"] = pixelgl.KeyM
	ui.Keymap["w"] = pixelgl.KeyW
	ui.Keymap["x"] = pixelgl.KeyX
	ui.Keymap["c"] = pixelgl.KeyC
	ui.Keymap["v"] = pixelgl.KeyV
	ui.Keymap["b"] = pixelgl.KeyB
	ui.Keymap["n"] = pixelgl.KeyN

	ui.Keymap["up"] = pixelgl.KeyUp
	ui.Keymap["down"] = pixelgl.KeyDown
	ui.Keymap["left"] = pixelgl.KeyLeft
	ui.Keymap["right"] = pixelgl.KeyRight

	ui.Keymap["tab"] = pixelgl.KeyTab
	ui.Keymap["escape"] = pixelgl.KeyEscape
	ui.Keymap["space"] = pixelgl.KeySpace
	ui.Keymap["l-alt"] = pixelgl.KeyLeftAlt
	ui.Keymap["r-alt"] = pixelgl.KeyRightAlt
	ui.Keymap["l-ctrl"] = pixelgl.KeyLeftControl
	ui.Keymap["r-ctrl"] = pixelgl.KeyRightControl
	ui.Keymap["l-shift"] = pixelgl.KeyLeftShift
	ui.Keymap["r-shift"] = pixelgl.KeyRightShift
}

func (ui *UI) GetInput() {

	input := ui.Game.GetInput()
	input2 := ui.Game.GetInput2()

	input.Typ = game.StayStill

	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Up]) {
		input.Typ = game.Up
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Down]) {
		input.Typ = game.Down
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Left]) {
		input.Typ = game.Left
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Right]) {
		input.Typ = game.Right
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Action]) {
		input.Typ = game.Action
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Power]) {
		input.Typ = game.Power
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Escape]) {
		input.Typ = game.Escape
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Select]) {
		input.Typ = game.Select
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Speed]) {
		input.Typ = game.SpeedUp
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Shadow]) {
		input.Typ = game.Shadow
	}
	if ui.win.JustPressed(ui.Keymap[ui.Game.Config.Keymap.Meditate]) {
		input.Typ = game.Meditate
	}

	if ui.win.JustReleased(ui.Keymap[ui.Game.Config.Keymap.Speed]) {
		input2.Typ = game.None
	}
	if ui.win.JustReleased(ui.Keymap[ui.Game.Config.Keymap.Shadow]) {
		input2.Typ = game.None
	}
	if ui.win.JustReleased(ui.Keymap[ui.Game.Config.Keymap.Meditate]) {
		input2.Typ = game.None
	}
}

func (ui *UI) WatchInput() {
	for {
		ui.GetInput()
	}
}
