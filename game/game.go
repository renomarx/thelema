package game

import (
	"log"
)

// Length of a map case (in pixels)
const CaseLen = 32

var EM *EventManager

type Game struct {
	GameDir       string
	Level         *Level
	Levels        map[string]*Level
	Books         map[string]*OBook
	QuestsObjects map[string]*QuestObject
	FightingRing  *FightingRing
	Paused        bool
	Running       bool
	Playing       bool
	input         *Input
	input2        *Input
	menu          *Menu
	GG            *GameGenerator
	Config        *Config
}

func (g *Game) GetInput() *Input {
	return g.input
}

func (g *Game) GetInput2() *Input {
	return g.input2
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
	Shadow
	Meditate
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
	Rune      string
	Static    bool
	Blocking  bool
	IsPlaying bool
}

func (o *Object) GetTile() Tile {
	return Tile(o.Rune)
}

type MovingObject struct {
	Object
	Xb int
	Yb int
}

func NewGame(gameDir string) *Game {
	game := &Game{Paused: false, Running: true, Playing: false, GameDir: gameDir}
	game.LoadConfig()
	game.LoadMenu()
	game.menu.IsOpen = true
	game.input = &Input{Typ: StayStill}
	game.input2 = &Input{Typ: None}
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

func (g *Game) UpdateLevel() {
	input := g.input
	if g.Level.Paused {
		g.HandleInputPlayerMenu()
	} else {
		g.handleInput()
		g.Level.handleMap()
		if input.Typ == Select {
			g.OpenPlayerMenu()
		}
	}
	if input.Typ == Escape {
		g.OpenMenu()
	}
}

func (g *Game) handleInput() {
	level := g.Level
	p := level.Player
	if !p.IsPlaying {
		p.IsPlaying = true
		go func(p *Player) {
			p.Update(g)
			p.IsPlaying = false
		}(p)
	}
}
