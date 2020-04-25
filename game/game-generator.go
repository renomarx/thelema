package game

type GameGenerator struct {
	IsOpen     bool
	Players    []*Player
	currentIdx int
}

func NewGameGenerator() *GameGenerator {
	gg := &GameGenerator{}
	gg.Players = GeneratePlayers()
	return gg
}

func (gg *GameGenerator) ChoiceRight() {
	gg.currentIdx = gg.currentIdx + 1
	if gg.currentIdx >= len(gg.Players) {
		gg.currentIdx = len(gg.Players) - 1
	}
}

func (gg *GameGenerator) ChoiceLeft() {
	gg.currentIdx = gg.currentIdx - 1
	if gg.currentIdx <= 0 {
		gg.currentIdx = 0
	}
}

func (gg *GameGenerator) ConfirmChoice(g *Game) {
	p := gg.Players[gg.currentIdx]
	gg.Close()
	g.GenerateWorld()
	g.LoadPlayer(p)
	g.Playing = true
	g.GG = nil
	g.CloseMenu()

	EM.Dispatch(&Event{
		Action:  ActionChangeLevel,
		Payload: map[string]string{"levelName": g.Level.Name},
		Message: "Premier niveau chargé"})
}

func (gg *GameGenerator) Close() {
	gg.IsOpen = false
}

func (gg *GameGenerator) IsHighlighted(idx int) bool {
	return gg.currentIdx == idx
}

func (gg *GameGenerator) GetCurrentPlayer() *Player {
	return gg.Players[gg.currentIdx]
}

func (gg *GameGenerator) HandleInput(g *Game) {
	input := g.GetInput()
	switch input.Typ {
	case Right:
		DispatchEventMenu(ActionMenuSelect)
		gg.ChoiceRight()
		adaptMenuSpeed()
	case Left:
		DispatchEventMenu(ActionMenuSelect)
		gg.ChoiceLeft()
		adaptMenuSpeed()
	case Action:
		DispatchEventMenu(ActionMenuConfirm)
		gg.ConfirmChoice(g)
		adaptMenuSpeed()
	case Escape:
		DispatchEventMenu(ActionMenuClose)
		gg.Close()
		adaptMenuSpeed()
	default:
	}
}
