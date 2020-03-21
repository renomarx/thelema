package main

import (
	"path/filepath"
	"thelema/game"
	"thelema/uipixel"
)

func main() {

	absPath, _ := filepath.Abs("game")

	game.EM = game.NewEventManager()

	g := game.NewGame(absPath)
	g.InitSlots()

	// ui := ui2d.NewUI(g)
	// game.EM.Subscribe(ui)
	//
	// go ui.WatchInput()
	// go ui.Run()

	go g.Run()

	ui := uipixel.UI{
		WindowWidth:  800.0,
		WindowHeight: 600.0,
	}
	ui.Run()
}
