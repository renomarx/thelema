package game

const PlayerMenuCmdCharacter = "Personnage"
const PlayerMenuCmdInventory = "Inventaire"
const PlayerMenuCmdLibrary = "Livres"
const PlayerMenuCmdQuests = "Journal"
const PlayerMenuCmdMap = "Carte"

func (p *Player) LoadPlayerMenu() {
	menu := &Menu{IsOpen: false}
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdCharacter, Highlighted: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdInventory})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdLibrary})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdQuests})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdMap})
	p.Menu = menu
}

func (g *Game) OpenPlayerMenu() {
	DispatchEventMenu(ActionMenuOpen)
	g.Level.Player.Menu.IsOpen = true
	g.Level.Paused = true
	adaptMenuSpeed()
}

func (g *Game) ClosePlayerMenu() {
	menu := g.Level.Player.Menu
	menu.ClearSelected()
	menu.IsOpen = false
	g.Level.Paused = false
}

func (g *Game) HandleInputPlayerMenu() {
	input := g.input
	menu := g.Level.Player.Menu
	sidx := menu.GetSelectedIndex()
	if sidx < 0 {
		switch input.Typ {
		case Up:
			DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceUp()
		case Down:
			DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceDown()
		case Action:
			DispatchEventMenu(ActionMenuOpen)
			c := menu.ConfirmChoice()
			switch c.Cmd {
			case PlayerMenuCmdInventory:
				g.Level.Player.Inventory.Open()
			case PlayerMenuCmdLibrary:
				g.Level.Player.Library.IsOpen = true
			case PlayerMenuCmdQuests:
				g.Level.Player.QuestMenuOpen = true
			case PlayerMenuCmdCharacter:
				g.Level.Player.CharacterMenuOpen = true
			case PlayerMenuCmdMap:
				g.Level.Player.MapMenuOpen = true
			}
			adaptMenuSpeed()
		case Select:
			DispatchEventMenu(ActionMenuClose)
			g.ClosePlayerMenu()
			adaptMenuSpeed()
		default:
		}
	} else {
		sc := menu.Choices[sidx]
		switch sc.Cmd {
		case PlayerMenuCmdInventory:
			g.Level.Player.Inventory.HandleInput(g)
		case PlayerMenuCmdLibrary:
			g.Level.Player.Library.HandleInput(g)
		case PlayerMenuCmdQuests:
			switch input.Typ {
			case Select:
				g.Level.Player.QuestMenuOpen = false
				DispatchEventMenu(ActionMenuClose)
				menu.ClearSelected()
				adaptMenuSpeed()
			default:
			}
		case PlayerMenuCmdCharacter:
			switch input.Typ {
			case Right:
				g.Level.Player.NextPower()
				DispatchEventMenu(ActionMenuSelect)
				adaptMenuSpeed()
			case Left:
				g.Level.Player.LastPower()
				DispatchEventMenu(ActionMenuSelect)
				adaptMenuSpeed()
			case Select:
				g.Level.Player.CharacterMenuOpen = false
				DispatchEventMenu(ActionMenuClose)
				menu.ClearSelected()
				adaptMenuSpeed()
			}
		case PlayerMenuCmdMap:
			switch input.Typ {
			case Select:
				g.Level.Player.MapMenuOpen = false
				DispatchEventMenu(ActionMenuClose)
				menu.ClearSelected()
				adaptMenuSpeed()
			default:
			}
		default:
		}
	}
}
