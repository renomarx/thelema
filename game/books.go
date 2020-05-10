package game

type BookInfo struct {
	Level          string   `yaml:"level"`
	PosX           int      `yaml:"posX"`
	PosY           int      `yaml:"posY"`
	PosZ           int      `yaml:"posZ"`
	PowersGiven    []string `yaml:"powers_given"`
	StepsFinishing []string `yaml:"steps_finishing"`
	StepsBeginning []string `yaml:"steps_beginning"`
}

type OBook struct {
	Rune           string
	Title          string
	Text           []string
	Powers         []string
	StepsFinishing []string `yaml:"steps_finishing"`
	StepsBeginning []string `yaml:"steps_beginning"`
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
	book, e := g.Books[o.Rune]
	if e {
		p.Library.AddBook(book)
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
		pp := g.Level.Player.NewPower(powername)
		EM.Dispatch(&Event{
			Message: "Vous avez appris la magie: '" + pp.Name + "' avec ce livre!",
			Action:  ActionPower,
			Payload: map[string]string{"type": powername}})
	}
	for _, stID := range b.StepsBeginning {
		g.beginStep(stID)
	}
	for _, stID := range b.StepsFinishing {
		g.finishStep(stID)
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
