package game

import (
	"bufio"
	"log"
	"math/rand"
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
	g.loadPnjsVIP()
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

func (g *Game) loadPnjsVIP() {
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
		if l.Type == LevelTypeCity {
			pos := l.GetRandomFreePos()
			if pos != nil {
				pnj := NewPnj(*pos, rune(pnjRunes[pnjNames[i]]), pnjNames[i], pnjVoices[pnjNames[i]])
				filename := g.GameDir + "/dialogs/" + pnj.Name + ".json"
				pnj.LoadDialogs(filename)
				l.Pnjs[*pos] = pnj
			}
			i++
			if i >= len(pnjNames) {
				break
			}
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

func (wg *WorldGenerator) generateTrees(level *Level, nbTrees int) {
	for i := 0; i < nbTrees; i++ {
		x := rand.Intn(WorldWidth)
		y := rand.Intn(WorldHeight)
		pos := Pos{X: x, Y: y}
		o := &Object{Rune: rune(Tree), Blocking: true}
		o.Pos = pos
		_, oe := level.Objects[pos]
		if !oe {
			level.Objects[pos] = o
		}
	}
}

func (wg *WorldGenerator) generateMonsters(level *Level, nbMonsters int) {
	bestiary := Bestiary()
	for i := 0; i < nbMonsters; i++ {
		x := rand.Intn(len(level.Map[0]))
		y := rand.Intn(len(level.Map))
		m := rand.Intn(len(bestiary))
		pos := Pos{X: x, Y: y}

		mt := bestiary[m]
		if canGo(level, pos) {
			level.Monsters[pos] = NewMonster(mt, pos)
		}
	}
}

func (wg *WorldGenerator) generateBooks(level *Level, nbBooks int) {
	for i := 0; i < nbBooks; i++ {
		x := rand.Intn(len(level.Map[0]))
		y := rand.Intn(len(level.Map))
		pos := Pos{X: x, Y: y}

		if canGo(level, pos) {
			b := &Object{Rune: rune(Book), Blocking: true}
			b.Pos = pos
			level.Objects[pos] = b
		}
	}
}

func (wg *WorldGenerator) generatePnjs(l *Level, nbPnjs int) {
	pnjNames := []string{
		"warrior",
		"doctor",
		"policeman",
		"artist",
		"lord",
		"monk",
	} // TODO : different number for each type
	pnjVoices := map[string]string{
		"warrior":   VoiceMaleStandard,
		"doctor":    VoiceFemaleStandard,
		"policeman": VoiceMaleStandard,
		"artist":    VoiceFemaleStandard,
		"lord":      VoiceMaleStandard,
		"monk":      VoiceMaleStandard,
	} // TODO : better sex handling
	pnjRunes := map[string]Tile{
		"warrior":   Warrior,
		"doctor":    Doctor,
		"policeman": Policeman,
		"artist":    Artist,
		"lord":      Lord,
		"monk":      Monk,
	}
	for i := 0; i < nbPnjs; i++ {
		j := i % len(pnjNames)
		pos := l.GetRandomFreePos()
		if pos != nil {
			pnj := NewPnj(*pos, rune(pnjRunes[pnjNames[j]]), pnjNames[j], pnjVoices[pnjNames[j]])
			filename := wg.g.GameDir + "/dialogs/common/" + pnj.Name + ".json"
			pnj.LoadDialogs(filename)
			l.Pnjs[*pos] = pnj
		}
	}
}
