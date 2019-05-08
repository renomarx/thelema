package game

import (
	"encoding/gob"
	"os"
)

func (g *Game) InitSlots() {
	gob.Register(Food{})
}

func SaveGame(g *Game, slot string) {
	if g.Level.Player.IsDead() {
		return
	}
	filepath := g.generateSlotFilepath(slot)
	Mux.Lock()
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Saving game..."})
	err := writeGob(filepath, g)
	if err != nil {
		panic(err)
	}
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Saved!"})
	Mux.Unlock()
}

func LoadGame(g *Game, slot string) {
	eventManager := g.GetEventManager()
	filepath := g.generateSlotFilepath(slot)
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Loading game..."})
	lg := NewGame(g.GameDir)
	err := readGob(filepath, lg)
	*g = *lg
	if err != nil {
		panic(err)
	}
	g.SetEventManager(eventManager)
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Loaded."})
}

func (g *Game) generateSlotFilepath(slot string) string {
	return g.GameDir + "/../saves/" + slot + ".gob"
}

func writeGob(filepath string, object interface{}) error {
	file, err := os.Create(filepath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filepath string, object interface{}) error {
	file, err := os.Open(filepath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
