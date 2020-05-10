package game

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	FirstLevelName = "arcanea/home"
	PlayerInitialX = 7
	PlayerInitialY = 5
)

func (g *Game) GenerateWorld() {
	firstLevel := g.loadLevels()
	g.loadPortals()
	g.loadNpcsVIP()
	g.loadBooks()
	g.loadQuestsObjects()
	g.LoadMonsters()
	g.Level = firstLevel
}

func (g *Game) LoadPlayer(p *Player) {
	gameDir := g.DataDir
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
	maps := LoadFilenames(g.DataDir + "/maps/cities")
	for _, filemap := range maps {
		fileArr := strings.Split(filemap, ".")
		if len(fileArr) == 2 {
			levelName := fileArr[0]
			if fileArr[1] == "map" {
				mapName := "cities/" + levelName
				l := g.LoadMapTemplate(mapName, levelName)
				g.generateNpcs(l, rand.Intn(len(l.Map))+1)
				g.Levels[levelName] = l
				g.loadHouses(mapName, levelName)
			}
		}
	}
}

func (g *Game) loadHouses(path, cityName string) {
	if _, err := os.Stat(g.DataDir + "/maps/" + path); os.IsNotExist(err) {
		return
	}
	maps := LoadFilenames(g.DataDir + "/maps/" + path)
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
	maps := LoadFilenames(g.DataDir + "/maps/dungeons")
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
	filepath := g.DataDir + "/maps/portals.txt"
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

func (g *Game) loadNpcsVIP() {
	npcNames := LoadFilenames(g.DataDir + "/npcs")
	for _, filename := range npcNames {
		fileArr := strings.Split(filename, ".")
		if len(fileArr) == 2 && fileArr[1] == "yaml" {
			p := Pos{}
			npc := NewNpc(p, fileArr[0])
			filename := g.DataDir + "/npcs/" + npc.Name + ".yaml"
			level, pos := npc.LoadNpc(filename)

			l, exists := g.Levels[level]
			if !exists {
				panic("Level " + level + " does not exist")
			}
			npc.Pos = pos
			l.Map[pos.Z][pos.Y][pos.X].Npc = npc
		}
	}
}

func (g *Game) loadBooks() {
	filename := g.DataDir + "/books/books.yaml"
	yamlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()
	byteValue, _ := ioutil.ReadAll(yamlFile)

	books := make(map[string]BookInfo)
	err = yaml.Unmarshal(byteValue, &books)
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
		pos := Pos{X: bookInfo.PosX, Y: bookInfo.PosY, Z: bookInfo.PosZ}
		physicalObj := &Object{Rune: tile, Blocking: true}
		physicalObj.Pos = pos
		l.Map[pos.Z][pos.Y][pos.X].Object = physicalObj
	}
}

func (g *Game) loadBookFromFile(filename string, bookInfo *BookInfo) *OBook {
	filepath := g.DataDir + "/books/" + filename + ".txt"
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
	filename := g.DataDir + "/quests/objects.yaml"
	yamlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()

	byteValue, _ := ioutil.ReadAll(yamlFile)

	objects := make(map[string]*QuestObject)

	err = yaml.Unmarshal(byteValue, &objects)
	if err != nil {
		log.Fatal(err)
	}

	objectsByRune := make(map[string]*QuestObject)
	for key, obj := range objects {
		l, exists := g.Levels[obj.Level]
		if !exists {
			log.Fatal("Level " + obj.Level + " does not exist")
		}
		pos := l.GetRandomFreePos(0) // FIXME
		if pos != nil {
			physicalObj := &Object{Rune: key, Blocking: true}
			physicalObj.Pos = *pos
			l.Map[pos.Z][pos.Y][pos.X].Object = physicalObj
			objectsByRune[key] = obj
		} else {
			log.Fatal("No place left on level " + obj.Level)
		}
	}

	g.QuestsObjects = objectsByRune
}

func (p *Player) LoadQuests(dirpath string) {
	filename := dirpath + "/quests/quests.yaml"
	yamlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()

	byteValue, _ := ioutil.ReadAll(yamlFile)

	quests := make(map[string]*Quest)

	yaml.Unmarshal(byteValue, &quests)

	p.Quests = quests
}

func (g *Game) generateNpcs(l *Level, nbNpcs int) {
	npcNames := []string{
		"warrior",
		"doctor",
		"policeman",
		"artist",
		"lord",
		"monk",
	} // TODO : different number for each type
	npcVoices := map[string]string{
		"warrior":   VoiceMaleStandard,
		"doctor":    VoiceFemaleStandard,
		"policeman": VoiceMaleStandard,
		"artist":    VoiceFemaleStandard,
		"lord":      VoiceMaleStandard,
		"monk":      VoiceMaleStandard,
	} // TODO : better sex handling
	for i := 0; i < nbNpcs; i++ {
		j := i % len(npcNames)
		pos := l.GetRandomFreePos(0) // FIXME
		if pos != nil {
			npc := NewNpc(*pos, npcNames[j])
			npc.Voice = npcVoices[npcNames[j]]
			filename := g.DataDir + "/npcs/common/" + npc.Name + ".yaml"
			npc.LoadNpc(filename)
			l.Map[pos.Z][pos.Y][pos.X].Npc = npc
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
		pos := l.GetRandomFreePos(0) // FIXME
		if pos != nil {
			physicalObj := &Object{Rune: string(objects[j]), Blocking: true}
			physicalObj.Pos = *pos
			l.Map[pos.Z][pos.Y][pos.X].Object = physicalObj
		}
	}
}
