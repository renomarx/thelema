package game

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
)

type GameSlot struct {
	Level         string
	Levels        map[string]*LevelSlot
	Books         map[string]*OBook
	QuestsObjects map[rune]*QuestObject
	Paused        bool
	Running       bool
	Playing       bool
}
type LevelSlot struct {
	Name        string
	Width       int
	Height      int
	Type        string
	Player      *Player
	Map         [][]Tile
	Portals     map[Pos]*Portal
	Monsters    map[Pos]*Monster
	Objects     map[Pos]*Object
	Effects     map[Pos]*Effect
	Projectiles map[Pos]*Projectile
	Pnjs        map[Pos]*Pnj
	Invocations map[Pos]*Invoked
	Friends     map[Pos]*Friend
	Enemies     map[Pos]*Enemy
	Paused      bool
	PRay        int
}

func (gs *GameSlot) FromGame(g *Game) {
	gs.Paused = g.Paused
	gs.Running = g.Running
	gs.Playing = g.Playing
	gs.Level = g.Level.Name
	gs.Levels = make(map[string]*LevelSlot)
	for _, level := range g.Levels {
		ls := &LevelSlot{}
		ls.FromLevel(level)
		gs.Levels[ls.Name] = ls
	}
	gs.Books = g.Books
	gs.QuestsObjects = g.QuestsObjects
}

func (gs *GameSlot) ToGame(g *Game) {
	g.Levels = make(map[string]*Level)
	for _, ls := range gs.Levels {
		level := &Level{}
		ls.ToLevel(level)
		g.Levels[ls.Name] = level
		if ls.Name == gs.Level {
			g.Level = level
		}
	}
	g.Books = gs.Books
	g.QuestsObjects = gs.QuestsObjects
}

func (ls *LevelSlot) FromLevel(l *Level) {
	ls.Name = l.Name
	ls.Width = l.Width
	ls.Height = l.Height
	ls.Type = l.Type
	ls.Player = l.Player
	ls.Paused = l.Paused
	ls.PRay = l.PRay
	ls.Map = l.Map
	ls.Portals = make(map[Pos]*Portal)
	for y, row := range l.Portals {
		for x, m := range row {
			if m != nil {
				ls.Portals[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Monsters = make(map[Pos]*Monster)
	for y, row := range l.Monsters {
		for x, m := range row {
			if m != nil {
				ls.Monsters[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Objects = make(map[Pos]*Object)
	for y, row := range l.Objects {
		for x, m := range row {
			if m != nil {
				ls.Objects[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Effects = make(map[Pos]*Effect)
	for y, row := range l.Effects {
		for x, m := range row {
			if m != nil {
				ls.Effects[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Projectiles = make(map[Pos]*Projectile)
	for y, row := range l.Projectiles {
		for x, m := range row {
			if m != nil {
				ls.Projectiles[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Pnjs = make(map[Pos]*Pnj)
	for y, row := range l.Pnjs {
		for x, m := range row {
			if m != nil {
				ls.Pnjs[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Invocations = make(map[Pos]*Invoked)
	for y, row := range l.Invocations {
		for x, m := range row {
			if m != nil {
				ls.Invocations[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Friends = make(map[Pos]*Friend)
	for y, row := range l.Friends {
		for x, m := range row {
			if m != nil {
				ls.Friends[Pos{X: x, Y: y}] = m
			}
		}
	}
	ls.Enemies = make(map[Pos]*Enemy)
	for y, row := range l.Enemies {
		for x, m := range row {
			if m != nil {
				ls.Enemies[Pos{X: x, Y: y}] = m
			}
		}
	}
}

func (ls *LevelSlot) ToLevel(l *Level) {
	l.Name = ls.Name
	l.Type = ls.Type
	l.Player = ls.Player
	l.Paused = ls.Paused
	l.PRay = ls.PRay
	l.InitMaps(ls.Height, ls.Width)
	l.Map = ls.Map
	for p, m := range ls.Portals {
		l.Portals[p.Y][p.X] = m
	}
	for p, m := range ls.Monsters {
		l.Monsters[p.Y][p.X] = m
	}
	for p, m := range ls.Objects {
		l.Objects[p.Y][p.X] = m
	}
	for p, m := range ls.Effects {
		l.Effects[p.Y][p.X] = m
	}
	for p, m := range ls.Projectiles {
		l.Projectiles[p.Y][p.X] = m
	}
	for p, m := range ls.Pnjs {
		l.Pnjs[p.Y][p.X] = m
	}
	for p, m := range ls.Invocations {
		l.Invocations[p.Y][p.X] = m
	}
	for p, m := range ls.Friends {
		l.Friends[p.Y][p.X] = m
	}
	for p, m := range ls.Enemies {
		l.Enemies[p.Y][p.X] = m
	}
}

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
	gs := &GameSlot{}
	gs.FromGame(g)
	err := writeGob(filepath, gs)
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
	gs := &GameSlot{}
	err := readGob(filepath, gs)
	gs.ToGame(lg)
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
