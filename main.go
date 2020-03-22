package main

import (
	"path/filepath"
	"thelema/game"
	"thelema/uipixel"
)

func main() {

	absPath, _ := filepath.Abs("data")

	game.EM = game.NewEventManager()

	g := game.NewGame(absPath)
	g.InitSlots()

	// ui := ui2d.NewUI(g)

	ui := uipixel.NewUI(g)
	game.EM.Subscribe(ui)

	// go ui.WatchInput()

	go g.Run()

	ui.Run()
}
