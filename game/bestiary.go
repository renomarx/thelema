package game

import (
	"log"
	"math/rand"
)

type MonsterType struct {
	Name           string
	Tile           Tile
	Health         int
	Energy         int
	Strength       int
	Speed          int
	Luck           int
	VisionRange    int
	Probability    int
	Aggressiveness int
}

func TRat() *MonsterType {
	return &MonsterType{
		Tile:           Rat,
		Name:           "Rat",
		Health:         50,
		Energy:         50,
		Strength:       10,
		Speed:          20,
		Luck:           10,
		VisionRange:    10,
		Probability:    100,
		Aggressiveness: 200,
	}
}
func TSpider() *MonsterType {
	return &MonsterType{
		Tile:           Spider,
		Name:           "Migale",
		Health:         50,
		Energy:         50,
		Strength:       20, // TODO : poison
		Speed:          15,
		Luck:           10,
		VisionRange:    10,
		Probability:    100,
		Aggressiveness: 100,
	}
}
func TSnake() *MonsterType {
	return &MonsterType{
		Tile:           Snake,
		Name:           "Vipère",
		Health:         100,
		Energy:         200,
		Strength:       20, // TODO : poison
		Speed:          30,
		Luck:           0,
		VisionRange:    10,
		Probability:    80,
		Aggressiveness: 75,
	}
}
func TCat() *MonsterType {
	return &MonsterType{
		Tile:           Cat,
		Name:           "Chat sauvage",
		Health:         100,
		Energy:         300,
		Strength:       25,
		Speed:          20,
		Luck:           20,
		VisionRange:    10,
		Probability:    70,
		Aggressiveness: 50,
	}
}
func TEagle() *MonsterType {
	return &MonsterType{
		Tile:           Eagle,
		Name:           "Aigle",
		Health:         150,
		Energy:         200,
		Strength:       40,
		Speed:          30,
		Luck:           25,
		VisionRange:    12,
		Probability:    50,
		Aggressiveness: 80,
	}
}
func TWolf() *MonsterType {
	return &MonsterType{
		Tile:           Wolf,
		Name:           "Loup",
		Health:         300,
		Energy:         300,
		Strength:       50,
		Speed:          20,
		Luck:           10,
		VisionRange:    10,
		Probability:    30,
		Aggressiveness: 400,
	}
}
func TBear() *MonsterType {
	return &MonsterType{
		Tile:           Bear,
		Name:           "Ours brun",
		Health:         500,
		Energy:         200,
		Strength:       80,
		Speed:          20,
		Luck:           20,
		VisionRange:    10,
		Probability:    20,
		Aggressiveness: 200,
	}
}
func TElephant() *MonsterType {
	return &MonsterType{
		Tile:           Elephant,
		Name:           "Elephant",
		Health:         1000,
		Energy:         200,
		Strength:       100,
		Speed:          20,
		Luck:           20,
		VisionRange:    8,
		Probability:    20,
		Aggressiveness: 100,
	}
}
func TBat() *MonsterType {
	return &MonsterType{
		Tile:           Bat,
		Name:           "Chauve-souris vampire",
		Health:         50,
		Energy:         50,
		Strength:       10, // TODO : magies
		Speed:          30,
		Luck:           10,
		VisionRange:    10,
		Probability:    80,
		Aggressiveness: 300,
	}
}
func TScorpion() *MonsterType {
	return &MonsterType{
		Tile:           Scorpion,
		Name:           "Scorpion",
		Health:         100,
		Energy:         50,
		Strength:       50, // TODO : poison
		Speed:          5,
		Luck:           10,
		VisionRange:    5,
		Probability:    50,
		Aggressiveness: 75,
	}
}

func (g *Game) MeetMonsters() {
	p := g.Level.Player
	if p.Shadow {
		return
	}
	l := g.Level
	r := rand.Intn(100000) % 100
	cc := l.Map[p.Z][p.Y][p.X]
	if r >= cc.MonstersProbability {
		return
	}
	switch cc.T {
	case MonsterFloor:
		bestiary := g.Level.Bestiary
		g.FightMonsters(bestiary)
	}
}

func (g *Game) FightMonsters(bestiary []*MonsterType) {
	var enemies []FighterInterface
	nb := rand.Intn(2) + 1
	for i := 0; i < nb; i++ {
		m := rand.Intn(len(bestiary))
		proba := rand.Intn(100)
		mt := bestiary[m]
		for proba > mt.Probability {
			m := rand.Intn(len(bestiary))
			proba = rand.Intn(100)
			mt = bestiary[m]
		}
		mo := NewMonster(mt)
		enemies = append(enemies, mo)
	}
	g.Fight(enemies)
}

func (g *Game) LoadMonsters() {
	defaultMonsters := []*MonsterType{
		TRat(),
		TSpider(),
		TScorpion(),
		TSnake(),
		TCat(),
		TEagle(),
		TBear(),
		TElephant(),
	}
	specificMonsters := make(map[string][]*MonsterType)
	specificMonsters["abigail_crypt"] = []*MonsterType{
		TRat(),
		TSpider(),
		TScorpion(),
	}
	for name, level := range g.Levels {
		log.Printf("Loading monsters of level %s", name)
		level.Bestiary = defaultMonsters
		b, e := specificMonsters[name]
		if e {
			level.Bestiary = b
		}
	}
}
