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
	firstLevel := g.loadLevels()
	g.loadPortals()
	g.loadPnjsVIP()
	g.loadBooks()
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
	g.generateObjects(l, rand.Intn(100))

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
				g.generatePnjs(l, rand.Intn(len(l.Map))+1)
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
				g.generateObjects(l, rand.Intn(20))
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
		key := ""
		if len(levels) >= 3 && levels[2] != "" {
			key = levels[2]
		}
		g.addBidirectionalPortal(level0Arr[0], Pos{X: x0, Y: y0}, level1Arr[0], Pos{X: x1, Y: y1}, key)
	}
}

func (g *Game) addBidirectionalPortal(srcName string, srcPos Pos, dstName string, dstPos Pos, key string) {
	srcLevel, e := g.Levels[srcName]
	if !e {
		panic("Level " + srcName + " does not exist.")
	}
	dstLevel, e := g.Levels[dstName]
	if !e {
		panic("Level " + dstName + " does not exist.")
	}
	srcLevel.AddPortal(srcPos, &Portal{LevelTo: dstName, PosTo: dstPos, Key: key})
	dstLevel.AddPortal(dstPos, &Portal{LevelTo: srcName, PosTo: srcPos, Key: key})
}

func (g *Game) loadPnjsVIP() {
	pnjNames := LoadFilenames(g.GameDir + "/pnjs")
	for _, filename := range pnjNames {
		fileArr := strings.Split(filename, ".")
		if len(fileArr) == 2 && fileArr[1] == "json" {
			p := Pos{}
			pnj := NewPnj(p, fileArr[0])
			filename := g.GameDir + "/pnjs/" + pnj.Name + ".json"
			level, pos := pnj.LoadPnj(filename)

			l, exists := g.Levels[level]
			if !exists {
				panic("Level " + level + " does not exist")
			}
			pnj.Pos = pos
			l.Map[pos.Y][pos.X].Pnj = pnj
		}
	}
}

func (g *Game) loadBooks() {
	filename := g.GameDir + "/books/books.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	books := make(map[string]BookInfo)
	err = json.Unmarshal(byteValue, &books)
	if err != nil {
		log.Fatal(err)
	}

	g.Books = make(map[string]*OBook)
	for tile, bookInfo := range books {
		g.Books[tile] = g.loadBookFromFile(tile, &bookInfo)
		levelName := bookInfo.Level
		l, exists := g.Levels[levelName]
		if !exists {
			log.Fatal("Level " + levelName + " does not exist")
		}
		pos := Pos{X: bookInfo.PosX, Y: bookInfo.PosY}
		physicalObj := &Object{Rune: tile, Blocking: true}
		physicalObj.Pos = pos
		l.Map[pos.Y][pos.X].Object = physicalObj
	}
}

func (g *Game) loadBookFromFile(filename string, bookInfo *BookInfo) *OBook {
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

	return &OBook{Title: title, Text: lines, Powers: bookInfo.PowersGiven, Rune: filename, Quest: bookInfo.Quest}
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

	objectsByRune := make(map[string]*QuestObject)
	for key, obj := range objects {
		l, exists := g.Levels[obj.Level]
		if !exists {
			log.Fatal("Level " + obj.Level + " does not exist")
		}
		pos := l.GetRandomFreePos()
		if pos != nil {
			physicalObj := &Object{Rune: key, Blocking: true}
			physicalObj.Pos = *pos
			l.Map[pos.Y][pos.X].Object = physicalObj
			objectsByRune[key] = obj
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
			pnj := NewPnj(*pos, pnjNames[j])
			pnj.Voice = pnjVoices[pnjNames[j]]
			filename := g.GameDir + "/pnjs/common/" + pnj.Name + ".json"
			pnj.LoadPnj(filename)
			l.Map[pos.Y][pos.X].Pnj = pnj
		}
	}
}

func (g *Game) generateObjects(l *Level, nbObjects int) {
	objects := []Tile{
		Fruits,
		Senzu,
		Bread,
		Water,
		Steak,
	}
	for i := 0; i < nbObjects; i++ {
		j := rand.Intn(42) % len(objects)
		pos := l.GetRandomFreePos()
		if pos != nil {
			physicalObj := &Object{Rune: string(objects[j]), Blocking: true}
			physicalObj.Pos = *pos
			l.Map[pos.Y][pos.X].Object = physicalObj
		}
	}
}
