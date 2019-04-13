package game

type Usable interface {
	Use(g *Game)
	GetObject() *Object
}

type InventoryUsable struct {
	Usable
	Highlighted bool
}

const InventoryUsablesMax = 100
const InitialGoldAmount = 400

type Inventory struct {
	Usables      []InventoryUsable
	QuestObjects map[rune]*Object
	Gold         int
	IsOpen       bool
}

func NewInventory() *Inventory {
	iv := &Inventory{IsOpen: false}
	iv.QuestObjects = make(map[rune]*Object)
	iv.Gold = InitialGoldAmount
	return iv
}

func (iv *Inventory) GetHighlightedIndex() int {
	for i := 0; i < len(iv.Usables); i++ {
		if iv.Usables[i].Highlighted {
			return i
		}
	}
	return -1
}

func (iv *Inventory) ClearHighlight() {
	for j := 0; j < len(iv.Usables); j++ {
		iv.Usables[j].Highlighted = false
	}
}

func (iv *Inventory) SetHighlightedIndex(i int) {
	iv.ClearHighlight()
	len := len(iv.Usables)
	if i < 0 {
		i = 0
	}
	if i >= len {
		i = len - 1
	}
	if i >= 0 {
		idx := i
		iv.Usables[idx].Highlighted = true
	}
}

func (iv *Inventory) GetCurrentChoice() InventoryUsable {
	idx := iv.GetHighlightedIndex()
	return iv.Usables[idx]
}

func (iv *Inventory) ConfirmChoice(g *Game) {
	idx := iv.GetHighlightedIndex()
	if idx >= 0 {
		iv.UseUsable(idx, g)
	}
}

func (iv *Inventory) ChoiceUp() {
	choiceIdx := iv.GetHighlightedIndex()
	iv.SetHighlightedIndex(choiceIdx - 1)
}

func (iv *Inventory) ChoiceDown() {
	choiceIdx := iv.GetHighlightedIndex()
	iv.SetHighlightedIndex(choiceIdx + 1)
}

func (iv *Inventory) TakeUsable(o *Object) bool {
	if len(iv.Usables) >= InventoryUsablesMax {
		return false
	}
	switch Tile(o.Rune) {
	case Senzu:
		u := USenzu{}
		u.Object = o
		iu := InventoryUsable{Highlighted: false}
		iu.Usable = u
		iv.Usables = append(iv.Usables, iu)
	default:
		return false
	}
	return true
}

func (iv *Inventory) UseUsable(i int, g *Game) {
	u := iv.Usables[i]
	u.Use(g)
	iv.deleteUsable(i)
}

func (iv *Inventory) deleteUsable(i int) {
	iv.Usables = append(iv.Usables[:i], iv.Usables[i+1:]...)
}

func (iv *Inventory) Open() {
	iv.IsOpen = true
}

func (iv *Inventory) Close(g *Game) {
	iv.IsOpen = false
	menu := g.Level.Player.Menu
	for _, c := range menu.Choices {
		if c.Cmd == PlayerMenuCmdInventory {
			c.Selected = false
		}
	}
}

func (iv *Inventory) HandleInput(g *Game) {
	input := g.GetInput()
	switch input.Typ {
	case Up:
		iv.ChoiceUp()
		adaptMenuSpeed()
	case Down:
		iv.ChoiceDown()
		adaptMenuSpeed()
	case Action:
		iv.ConfirmChoice(g)
		adaptMenuSpeed()
	case Select:
		iv.Close(g)
		adaptMenuSpeed()
	default:
	}
}
