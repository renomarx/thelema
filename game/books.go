package game

type OBook struct {
	Rune   rune
	Title  string
	Text   []string
	Powers []string
}

type Library struct {
	IsOpen     bool
	Books      []*OBook
	currentIdx int
}

func NewLibrary() *Library {
	lib := &Library{}
	lib.currentIdx = 0
	return lib
}

func (l *Library) AddBook(book *OBook) {
	l.Books = append(l.Books, book)
}

func (p *Player) AddBook(o *Object, g *Game) bool {
	if Tile(o.Rune) != Book {
		return false
	}
	for key, book := range g.Books {
		p.Library.AddBook(book)
		delete(g.Books, key)
		return true
	}

	return false
}

func (l *Library) ChoiceRight() {
	l.currentIdx = l.currentIdx + 1
	if l.currentIdx >= len(l.Books) {
		l.currentIdx = len(l.Books) - 1
	}
}

func (l *Library) ChoiceLeft() {
	l.currentIdx = l.currentIdx - 1
	if l.currentIdx <= 0 {
		l.currentIdx = 0
	}
}

func (l *Library) ConfirmChoice(g *Game) {
	b := l.GetCurrentBook()
	if b == nil {
		return
	}
	for _, powername := range b.Powers {
		g.Level.Player.NewPower(powername, g)
	}
}

func (l *Library) Close(g *Game) {
	l.IsOpen = false
	menu := g.Level.Player.Menu
	for _, c := range menu.Choices {
		if c.Cmd == PlayerMenuCmdLibrary {
			c.Selected = false
		}
	}
}

func (l *Library) IsHighlighted(idx int) bool {
	return l.currentIdx == idx
}

func (l *Library) GetCurrentBook() *OBook {
	if len(l.Books) == 0 {
		return nil
	}
	return l.Books[l.currentIdx]
}

func (l *Library) HandleInput(g *Game) {
	input := g.GetInput()
	switch input.Typ {
	case Right:
		g.GetEventManager().Dispatch(&Event{Action: ActionReadBook})
		l.ChoiceRight()
		l.ConfirmChoice(g)
		adaptMenuSpeed()
	case Left:
		g.GetEventManager().Dispatch(&Event{Action: ActionReadBook})
		l.ChoiceLeft()
		l.ConfirmChoice(g)
		adaptMenuSpeed()
	case Select:
		g.GetEventManager().Dispatch(&Event{Action: ActionMenuClose})
		l.Close(g)
		adaptMenuSpeed()
	default:
	}
}
