package game

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	FirstLevelName = "arcanea/home"
	PlayerInitialX = 7
	PlayerInitialY = 5
)

func (g *Game) GenerateWorld() {
	g.loadBooks()
	firstLevel := g.loadLevels()
	g.loadPortals()
	g.loadPnjsVIP()
	g.loadQuestsObjects()
	g.Level = firstLevel
}

func (g *Game) LoadPlayer(p *Player) {
	gameDir := g.GameDir
	p.X = PlayerInitialX
	p.Y = PlayerInitialY
	p.LoadQuests(gameDir)
	p.LoadPlayerMenu()
	g.Level.Player = p
}

func (g *Game) loadLevels() *Level {
	g.Levels = make(map[string]*Level)

	worldName := "world"
	l := g.LoadMapTemplate(worldName, worldName)
	g.Levels[worldName] = l
	g.loadCities()
	g.loadDungeons()

	return g.Levels[FirstLevelName]
}

func (g *Game) loadCities() {
	maps := LoadFilenames(g.GameDir + "/maps/cities")
	for _, filemap := range maps {
		fileArr := strings.Split(filemap, ".")
		if len(fileArr) == 2 {
			levelName := fileArr[0]
			if fileArr[1] == "map" {
				mapName := "cities/" + levelName
				l := g.LoadMapTemplate(mapName, levelName)
				g.generatePnjs(l, rand.Intn(20)+1)
				g.Levels[levelName] = l
				g.loadHouses(mapName, levelName)
			}
		}
	}
}

func (g *Game) loadHouses(path, cityName string) {
	if _, err := os.Stat(g.GameDir + "/maps/" + path); os.IsNotExist(err) {
		return
	}
	maps := LoadFilenames(g.GameDir + "/maps/" + path)
	for _, filemap := range maps {
		fileArr := strings.Split(filemap, ".")
		if len(fileArr) == 2 {
			levelName := cityName + "/" + fileArr[0]
			if fileArr[1] == "map" {
				mapName := path + "/" + fileArr[0]
				l := g.LoadMapTemplate(mapName, levelName)
				g.Levels[levelName] = l
			}
		}
	}
}

func (g *Game) loadDungeons() {
	maps := LoadFilenames(g.GameDir + "/maps/dungeons")
	for _, filemap := range maps {
		fileArr := strings.Split(filemap, ".")
		if len(fileArr) == 2 {
			levelName := fileArr[0]
			if fileArr[1] == "map" {
				mapName := "dungeons/" + levelName
				l := g.LoadMapTemplate(mapName, levelName)
				g.generatePnjs(l, rand.Intn(20)+1)
				g.Levels[fileArr[0]] = l
			}
		}
	}
}

func (g *Game) loadPortals() {
	filepath := g.GameDir + "/maps/portals.txt"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		levels := strings.Split(line, ">")
		if len(levels) < 2 {
			panic("Bad line format, expecting [a-Z]+,[0-9]+,[0-9]+>[a-Z]+,[0-9]+,[0-9]+: " + line)
		}
		level0Arr := strings.Split(levels[0], ",")
		if len(level0Arr) < 3 {
			panic("Bad line format, expecting [a-Z]+,[0-9]+,[0-9]+>[a-Z]+,[0-9]+,[0-9]+: " + line)
		}
		x0, _ := strconv.Atoi(level0Arr[1])
		y0, _ := strconv.Atoi(level0Arr[2])
		level1Arr := strings.Split(levels[1], ",")
		if len(level1Arr) < 3 {
			panic("Bad line format, expecting [a-Z]+,[0-9]+,[0-9]+>[a-Z]+,[0-9]+,[0-9]+: " + line)
		}
		x1, _ := strconv.Atoi(level1Arr[1])
		y1, _ := strconv.Atoi(level1Arr[2])
		g.addBidirectionalPortal(level0Arr[0], Pos{X: x0, Y: y0}, level1Arr[0], Pos{X: x1, Y: y1})
	}
}

func (g *Game) addBidirectionalPortal(srcName string, srcPos Pos, dstName string, dstPos Pos) {
	srcLevel, e := g.Levels[srcName]
	if !e {
		panic("Level " + srcName + " does not exist.")
	}
	dstLevel, e := g.Levels[dstName]
	if !e {
		panic("Level " + dstName + " does not exist.")
	}
	srcLevel.AddPortal(srcPos, &Portal{LevelTo: dstName, PosTo: dstPos})
	dstLevel.AddPortal(dstPos, &Portal{LevelTo: srcName, PosTo: srcPos})
}

func (g *Game) loadPnjsVIP() {
	pnjNames := []string{
		"jason",
		"sarah",
		"nathaniel",
	}
	pnjVoices := map[string]string{
		"jason":     VoiceMaleStandard,
		"sarah":     VoiceFemaleStandard,
		"nathaniel": VoiceMaleStandard,
	}
	for _, name := range pnjNames {
		p := Pos{}
		pnj := NewPnj(p, name, pnjVoices[name])
		filename := g.GameDir + "/pnjs/" + pnj.Name + ".json"
		pnj.LoadDialogs(filename)

		l, exists := g.Levels[pnj.Dialog.Level]
		if !exists {
			log.Fatal("Level " + pnj.Dialog.Level + " does not exist")
		}
		pos := l.GetRandomFreePos()
		if pos == nil {
			log.Fatal("No place left on level " + pnj.Dialog.Level)
		}
		pnj.Pos = *pos
		l.Map[pos.Y][pos.X].Pnj = pnj
	}
}

func (g *Game) loadBooks() {
	g.Books = make(map[string]*OBook)

	books := LoadFilenames(g.GameDir + "/books")
	for _, bookFile := range books {
		book := strings.Split(bookFile, ".")
		bookName := book[0]
		powers := []string{}
		switch bookName {
		case "invocat":
			powers = append(powers, PowerInvocation)
		case "dead_speaking":
			powers = append(powers, PowerDeadSpeaking)
		}
		g.Books[bookName] = g.loadBookFromFile(bookName, powers)
	}
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

func (g *Game) generateUsables(level *Level, objects []Tile, nb int) {
	for i := 0; i < nb; i++ {
		x := rand.Intn(len(level.Map[0]))
		y := rand.Intn(len(level.Map))
		m := rand.Intn(len(objects))
		pos := Pos{X: x, Y: y}

		mt := objects[m]
		if canGo(level, pos) {
			b := &Object{Rune: rune(mt), Blocking: true}
			b.Pos = pos
			level.Map[pos.Y][pos.X].Object = b
		}
	}
}

func (g *Game) generateBooks(level *Level, nbBooks int) {
	for i := 0; i < nbBooks; i++ {
		x := rand.Intn(len(level.Map[0]))
		y := rand.Intn(len(level.Map))
		pos := Pos{X: x, Y: y}

		if canGo(level, pos) {
			b := &Object{Rune: rune(Book), Blocking: true}
			b.Pos = pos
			level.Map[pos.Y][pos.X].Object = b
		}
	}
}

func (g *Game) loadQuestsObjects() {
	filename := g.GameDir + "/quests/objects.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	objects := make(map[string]*QuestObject)

	err = json.Unmarshal(byteValue, &objects)
	if err != nil {
		log.Fatal(err)
	}

	objectsByRune := make(map[rune]*QuestObject)
	for key, obj := range objects {
		l, exists := g.Levels[obj.Level]
		if !exists {
			log.Fatal("Level " + obj.Level + " does not exist")
		}
		pos := l.GetRandomFreePos()
		if pos != nil {
			rune := rune(key[0])
			physicalObj := &Object{Rune: rune, Blocking: true}
			physicalObj.Pos = *pos
			l.Map[pos.Y][pos.X].Object = physicalObj
			objectsByRune[rune] = obj
		} else {
			log.Fatal("No place left on level " + obj.Level)
		}
	}

	g.QuestsObjects = objectsByRune
}

func (g *Game) generatePnjs(l *Level, nbPnjs int) {
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
	for i := 0; i < nbPnjs; i++ {
		j := i % len(pnjNames)
		pos := l.GetRandomFreePos()
		if pos != nil {
			pnj := NewPnj(*pos, pnjNames[j], pnjVoices[pnjNames[j]])
			filename := g.GameDir + "/pnjs/common/" + pnj.Name + ".json"
			pnj.LoadDialogs(filename)
			l.Map[pos.Y][pos.X].Pnj = pnj
		}
	}
}
