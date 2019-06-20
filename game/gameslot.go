package game

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"log"
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
	return g.GameDir + "/../saves/" + slot + ".sav"
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

func writeJSON(filepath string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(filepath)
	if err == nil {
		zw := gzip.NewWriter(file)
		_, err := zw.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	file.Close()
	return err
}

func readJSON(filepath string, object interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	var b []byte
	zr, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal("error reader", err)
	}
	zr.Read(b)
	json.Unmarshal(b, object)
	return err
}
