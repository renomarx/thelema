package game

const FightingMenuCmdAttack = "Attaquer"
const FightingMenuCmdInventory = "Inventaire"
const PlayerMenuCmdRun = "Fuir"

func (g *Game) LoadFightingMenu() {
	menu := &Menu{IsOpen: false}
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: FightingMenuCmdAttack, Highlighted: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: FightingMenuCmdInventory})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdRun})
	g.FightingMenu = menu
}

func (g *Game) OpenFightingMenu() {
	g.DispatchEventMenu(ActionMenuOpen)
	g.FightingMenu.IsOpen = true
	adaptMenuSpeed()
}

func (g *Game) CloseFightingMenu() {
	menu := g.FightingMenu
	menu.ClearSelected()
	menu.IsOpen = false
}

func (g *Game) HandleInputFightingMenu() {
	input := g.input
	menu := g.FightingMenu
	sidx := menu.GetSelectedIndex()
	if sidx < 0 {
		switch input.Typ {
		case Up:
			g.DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceUp()
		case Down:
			g.DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceDown()
		case Action:
			g.DispatchEventMenu(ActionMenuOpen)
			c := menu.ConfirmChoice()
			switch c.Cmd {
			case FightingMenuCmdAttack:
				// TODO
				g.FightingRing.SelectedPlayerAction = "attack:sword"
				g.CloseFightingMenu()
				adaptMenuSpeed()
			case FightingMenuCmdInventory:
				// TODO
				g.CloseFightingMenu()
				adaptMenuSpeed()
			case PlayerMenuCmdRun:
				// TODO
				g.FightingRing.SelectedPlayerAction = "run"
				g.DispatchEventMenu(ActionMenuClose)
				g.CloseFightingMenu()
				adaptMenuSpeed()
			}
			adaptMenuSpeed()
		default:
		}
	} else {
		sc := menu.Choices[sidx]
		switch sc.Cmd {
		case FightingMenuCmdAttack:
			// TODO
			g.Level.Player.Library.HandleInput(g)
		case FightingMenuCmdInventory:
			// TODO
			g.Level.Player.Inventory.HandleInput(g)
		default:
		}
	}
}
