package main

import (
	"path/filepath"
	"thelema/game"
	"thelema/ui2d"
)

func main() {
	absPath, _ := filepath.Abs("game")

	game.EM = game.NewEventManager()

	g := game.NewGame(absPath)
	g.InitSlots()

	ui := ui2d.NewUI(g)
	game.EM.Subscribe(ui)

	go ui.WatchInput()
	go ui.Run()

	g.Run()
}
