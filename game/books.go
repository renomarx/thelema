package game

type OBook struct {
	Rune   string
	Title  string
	Score  int
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
	for _, b := range l.Books {
		if book.Title == b.Title {
			return
		}
	}
	l.Books = append(l.Books, book)
}

func (p *Player) AddBook(o *Object, g *Game) bool {
	// TODO : not taking book when already have it
	if Tile(o.Rune) != Book {
		return false
	}
	var firstBook *OBook
	firstScore := 1000
	for _, book := range g.Books {
		if book.Score <= firstScore {
			firstBook = book
			firstScore = book.Score
			if book.Score < 1000 {
				book.Score += 1
			}
		}
	}
	if firstBook != nil {
		p.Library.AddBook(firstBook)
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
		EM.Dispatch(&Event{Action: ActionReadBook})
		l.ChoiceRight()
		l.ConfirmChoice(g)
		adaptMenuSpeed()
	case Left:
		EM.Dispatch(&Event{Action: ActionReadBook})
		l.ChoiceLeft()
		l.ConfirmChoice(g)
		adaptMenuSpeed()
	case Select:
		EM.Dispatch(&Event{Action: ActionMenuClose})
		l.Close(g)
		adaptMenuSpeed()
	default:
	}
}
