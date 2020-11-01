package main

import (
	"path/filepath"
	"runtime"

	"github.com/renomarx/thelema/game"
	"github.com/renomarx/thelema/ui2d"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	absPath, _ := filepath.Abs("data")

	game.EM = game.NewEventManager()

	g := game.NewGame(absPath)
	g.InitSlots()

	ui := ui2d.NewUI(g)
	game.EM.Subscribe(ui)

	go g.Run()

	go ui.WatchInput()
	ui.Run()
}
