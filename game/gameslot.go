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
	EM.Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Sauvegarde..."})
	err := writeGob(filepath, g)
	if err != nil {
		panic(err)
	}
	EM.Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Sauvegardé!"})
}

func LoadGame(g *Game, slot string) {
	filepath := g.generateSlotFilepath(slot)
	EM.Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Chargement..."})
	lg := NewGame(g.DataDir)
	err := readGob(filepath, lg)
	*g = *lg
	if err != nil {
		panic(err)
	}
	EM.Dispatch(&Event{
		Action:  ActionMenuConfirm,
		Message: "Chargé."})
}

func (g *Game) generateSlotFilepath(slot string) string {
	return g.DataDir + "/../saves/" + slot + ".sav"
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
