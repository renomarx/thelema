package game

import (
	"bufio"
	"log"
	"os"
)

const (
	PlayerInitialX = 87
	PlayerInitialY = 14
	WorldHeight    = 500
	WorldWidth     = 1000
	OceanX         = 20
	OceanY         = 10
	WorldName      = "world"
)

type WorldGenerator struct {
	g *Game
}

func (g *Game) GenerateWorld() {
	g.loadBooks()
	firstLevel := g.loadLevels()
	g.loadPnjs()
	g.Level = firstLevel
}

func (g *Game) LoadPlayer(p *Player) {
	gameDir := g.GameDir
	p.X = PlayerInitialX
	p.Y = PlayerInitialY
	p.LoadQuests(gameDir)
	p.LoadQuestsObjects(gameDir)
	p.LoadPlayerMenu()
	g.Level.Player = p
}

func (g *Game) loadLevels() *Level {
	g.Levels = make(map[string]*Level)

	wg := &WorldGenerator{g: g}
	firstLevel := wg.generateOutdoor(WorldName)
	return firstLevel
}

func (g *Game) loadPnjs() {
	pnjNames := []string{
		"jason",
		"sarah",
	} // TODO : load automatically from dialogs
	pnjVoices := map[string]string{
		"jason": VoiceMaleStandard,
		"sarah": VoiceFemaleStandard,
	}
	pnjRunes := map[string]Tile{
		"jason": Jason,
		"sarah": Sarah,
	}
	i := 0
	for _, l := range g.Levels {
		pos := l.GetRandomFreePos()
		pnj := NewPnj(pos, rune(pnjRunes[pnjNames[i]]), pnjNames[i], pnjVoices[pnjNames[i]])
		pnj.LoadDialogs(g.GameDir)
		l.Pnjs[pos] = pnj
		i++
		if i >= len(pnjNames) {
			break
		}
	}
}

func (g *Game) loadBooks() {
	g.Books = make(map[string]*OBook)
	g.Books["cats"] = g.loadBookFromFile("cats", []string{})
	g.Books["invocat"] = g.loadBookFromFile("invocat", []string{PowerInvocation})
	g.Books["rats"] = g.loadBookFromFile("rats", []string{})
	g.Books["spiders"] = g.loadBookFromFile("spiders", []string{})
}

func (g *Game) loadBookFromFile(filename string, powers []string) *OBook {
	filepath := g.GameDir + "/books/" + filename + ".txt"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	title := ""
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(lines) > 0 {
		title = lines[0]
	}

	return &OBook{Title: title, Text: lines, Powers: powers, Rune: rune(Book)}
}
