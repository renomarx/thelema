package game

const PlayerMenuCmdCharacter = "Personnage"
const PlayerMenuCmdInventory = "Inventaire"
const PlayerMenuCmdLibrary = "Bibliothèque"
const PlayerMenuCmdQuests = "Journal"

func (p *Player) LoadPlayerMenu() {
	menu := &Menu{IsOpen: false}
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdCharacter, Highlighted: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdInventory})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdLibrary})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdQuests})
	p.Menu = menu
}

func (g *Game) OpenPlayerMenu() {
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
			menu.ChoiceUp()
		case Down:
			menu.ChoiceDown()
		case Action:
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
			}
			adaptMenuSpeed()
		case Select:
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
				menu.ClearSelected()
				adaptMenuSpeed()
			default:
			}
		case PlayerMenuCmdCharacter:
			switch input.Typ {
			case Right:
				g.Level.Player.NextPower()
				adaptMenuSpeed()
			case Left:
				g.Level.Player.LastPower()
				adaptMenuSpeed()
			case Select:
				g.Level.Player.CharacterMenuOpen = false
				menu.ClearSelected()
				adaptMenuSpeed()
			}
		default:
		}
	}
}
