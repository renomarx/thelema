package game

import "time"

type Menu struct {
	Choices []*MenuChoice
	IsOpen  bool
}

type MenuChoice struct {
	Cmd         string
	Highlighted bool
	Selected    bool
	Disabled    bool
}

const MenuCmdNew = "Nouvelle partie"
const MenuCmdSave = "Sauvegarder"
const MenuCmdLoad = "Charger"
const MenuCmdQuit = "Quitter"

func (menu *Menu) GetHighlightedIndex() int {
	for i := 0; i < len(menu.Choices); i++ {
		if menu.Choices[i].Highlighted {
			return i
		}
	}
	return 0
}

func (menu *Menu) ClearHighlight() {
	for j := 0; j < len(menu.Choices); j++ {
		menu.Choices[j].Highlighted = false
	}
}

func (menu *Menu) SetHighlightedIndex(i int) *MenuChoice {
	menu.ClearHighlight()
	len := len(menu.Choices)
	if i < 0 {
		i = 0
	}
	if i >= len {
		i = len - 1
	}
	idx := i
	menu.Choices[idx].Highlighted = true
	return menu.Choices[idx]
}

func (menu *Menu) GetHighlightedChoice() *MenuChoice {
	idx := menu.GetHighlightedIndex()
	return menu.Choices[idx]
}

func (menu *Menu) ClearSelected() {
	for _, c := range menu.Choices {
		c.Selected = false
	}
}

func (menu *Menu) ConfirmChoice() *MenuChoice {
	menu.ClearSelected()
	choice := menu.GetHighlightedChoice()
	choice.Selected = true
	adaptMenuSpeed()
	return choice
}

func (menu *Menu) GetSelectedIndex() int {
	for i := 0; i < len(menu.Choices); i++ {
		if menu.Choices[i].Selected {
			return i
		}
	}
	return -1
}

func (menu *Menu) ChoiceUp() {
	choiceIdx := menu.GetHighlightedIndex()
	c := menu.SetHighlightedIndex(choiceIdx - 1)
	if c.Disabled {
		menu.ChoiceUp()
		return
	}
	adaptMenuSpeed()
}

func (menu *Menu) ChoiceDown() {
	choiceIdx := menu.GetHighlightedIndex()
	c := menu.SetHighlightedIndex(choiceIdx + 1)
	if c.Disabled {
		menu.ChoiceDown()
		return
	}
	adaptMenuSpeed()
}

func adaptMenuSpeed() {
	time.Sleep(200 * time.Millisecond)
}

func (g *Game) LoadMenu() {
	menu := &Menu{}
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: MenuCmdNew, Highlighted: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: MenuCmdSave, Disabled: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: MenuCmdLoad})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: MenuCmdQuit})
	g.menu = menu
}

func (g *Game) OpenMenu() {
	g.ClosePlayerMenu()
	g.menu.Choices[1].Disabled = !g.Playing || g.Level.Player.IsDead()
	g.menu.IsOpen = true
	g.Paused = true
	adaptMenuSpeed()
}

func (g *Game) CloseMenu() {
	if g.Playing {
		g.menu.IsOpen = false
		g.Paused = false
		adaptMenuSpeed()
	}
}

func (g *Game) HandleInputMenu() {
	if g.GG != nil && g.GG.IsOpen {
		g.GG.HandleInput(g)
	} else {
		input := g.input
		switch input.Typ {
		case Up:
			g.menu.ChoiceUp()
		case Down:
			g.menu.ChoiceDown()
		case Action:
			c := g.menu.ConfirmChoice()
			switch c.Cmd {
			case MenuCmdNew:
				g.GG = NewGameGenerator()
				g.GG.IsOpen = true
			case MenuCmdSave:
				SaveGame(g, "slot1")
				// TODO
			case MenuCmdLoad:
				LoadGame(g, "slot1")
				// TODO
				g.Playing = true
				g.menu.Choices[1].Disabled = false
				g.CloseMenu()

				g.GetEventManager().Dispatch(&Event{
					Type:    PlayerEventsType,
					Action:  ActionChangeLevel,
					Payload: map[string]string{"levelType": g.Level.Type},
					Message: "Loaded level"})
			case MenuCmdQuit:
				g.Running = false
			default:
			}
		case Escape:
			g.CloseMenu()
		default:
		}
	}
}
