package game

import (
	"log"
	"sync"
)

// Length of a map case (in pixels)
const CaseLen = 32

type Game struct {
	GameDir       string
	Level         *Level
	Levels        map[string]*Level
	Books         map[string]*OBook
	QuestsObjects map[rune]*QuestObject
	Paused        bool
	Running       bool
	Playing       bool
	input         *Input
	input2        *Input
	eventManager  *EventManager
	menu          *Menu
	GG            *GameGenerator
	Config        *Config
	Mux           *sync.Mutex
}

func (g *Game) GetInput() *Input {
	return g.input
}

func (g *Game) GetInput2() *Input {
	return g.input2
}

func (g *Game) GetEventManager() *EventManager {
	return g.eventManager
}

func (g *Game) SetEventManager(em *EventManager) {
	g.eventManager = em
}

func (g *Game) GetMenu() *Menu {
	return g.menu
}

type InputType int

const (
	Up InputType = iota
	Down
	Left
	Right
	Quit
	Action
	Power
	StayStill
	Escape
	Select

	None
	SpeedUp
)

type Input struct {
	Typ InputType
}

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
}

type Object struct {
	Entity
	Rune      rune
	Blocking  bool
	IsPlaying bool
}

type MovingObject struct {
	Object
	Xb int
	Yb int
}

func NewGame(gameDir string) *Game {
	game := &Game{Paused: false, Running: true, Playing: false, GameDir: gameDir}
	game.Mux = &sync.Mutex{}
	game.LoadConfig()
	game.LoadMenu()
	game.menu.IsOpen = true
	game.input = &Input{Typ: StayStill}
	game.input2 = &Input{Typ: None}
	game.eventManager = NewEventManager()
	return game
}

func (g *Game) Run() {

	for g.Running {
		input := g.input
		if input.Typ == Quit {
			log.Println("Quit")
			break
		}
		if g.menu.IsOpen {
			g.HandleInputMenu()
		} else {
			g.UpdateLevel()
		}
	}
}
