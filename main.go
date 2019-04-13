package main

import (
	"path/filepath"
	"thelema/game"
	"thelema/ui2d"
)

func main() {
	absPath, _ := filepath.Abs("game")

	game := game.NewGame(absPath)
	game.InitSlots()

	ui := ui2d.NewUI(game)
	go ui.WatchInput()
	go ui.Run()

	game.Run()
}
