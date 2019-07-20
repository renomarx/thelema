package game

const FightingMenuCmdAttack = "Attaquer"
const FightingMenuCmdInventory = "Inventaire"
const PlayerMenuCmdRun = "Fuir"

func (fr *FightingRing) LoadFightingMenu() {
	menu := &Menu{IsOpen: false}
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: FightingMenuCmdAttack, Highlighted: true})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: FightingMenuCmdInventory})
	menu.Choices = append(menu.Choices, &MenuChoice{Cmd: PlayerMenuCmdRun})
	fr.Menu = menu
}

func (fr *FightingRing) OpenFightingMenu() {
	DispatchEventMenu(ActionMenuOpen)
	fr.Menu.IsOpen = true
	adaptMenuSpeed()
}

func (fr *FightingRing) CloseFightingMenu() {
	menu := fr.Menu
	menu.ClearSelected()
	menu.IsOpen = false
}

func (fr *FightingRing) HandleInputFightingMenu(input *Input) {
	menu := fr.Menu
	sidx := menu.GetSelectedIndex()
	if sidx < 0 {
		switch input.Typ {
		case Left:
			DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceUp()
		case Right:
			DispatchEventMenu(ActionMenuSelect)
			menu.ChoiceDown()
		case Action:
			DispatchEventMenu(ActionMenuOpen)
			c := menu.ConfirmChoice()
			switch c.Cmd {
			case FightingMenuCmdAttack:
				fr.AttacksMenuOpen = true
				adaptMenuSpeed()
			case FightingMenuCmdInventory:
				// TODO
				fr.CloseFightingMenu()
				adaptMenuSpeed()
			case PlayerMenuCmdRun:
				fr.SelectedPlayerAction = "run"
				DispatchEventMenu(ActionMenuClose)
				fr.CloseFightingMenu()
				adaptMenuSpeed()
			}
			adaptMenuSpeed()
		default:
		}
	} else {
		sc := menu.Choices[sidx]
		switch sc.Cmd {
		case FightingMenuCmdAttack:

			if fr.AttackTargetSelectionOpen {
				switch input.Typ {
				case Right:
					fr.NextTarget()
					DispatchEventMenu(ActionMenuSelect)
					adaptMenuSpeed()
				case Left:
					fr.LastTarget()
					DispatchEventMenu(ActionMenuSelect)
					adaptMenuSpeed()
				case Power:
					fr.AttackTargetSelectionOpen = false
				case Action:
					fr.SelectedPlayerAction = "attack"
					fr.AttackTargetSelectionOpen = false
					fr.AttacksMenuOpen = false
					DispatchEventMenu(ActionMenuClose)
					menu.ClearSelected()
					fr.CloseFightingMenu()
					adaptMenuSpeed()
				}

			} else {
				switch input.Typ {
				case Right:
					fr.NextPossibleAttack()
					DispatchEventMenu(ActionMenuSelect)
					adaptMenuSpeed()
				case Left:
					fr.LastPossibleAttack()
					DispatchEventMenu(ActionMenuSelect)
					adaptMenuSpeed()
				case Power:
					fr.AttacksMenuOpen = false
					DispatchEventMenu(ActionMenuClose)
					menu.ClearSelected()
					adaptMenuSpeed()
				case Action:
					if fr.GetSelectedAttack().Range == 0 {
						fr.SelectedPlayerAction = "attack"
						fr.AttacksMenuOpen = false
						DispatchEventMenu(ActionMenuClose)
						menu.ClearSelected()
						fr.CloseFightingMenu()
						adaptMenuSpeed()
					} else {
						fr.AttackTargetSelectionOpen = true
						adaptMenuSpeed()
					}
				}

			}
		case FightingMenuCmdInventory:
			// TODO
		default:
		}
	}
}
