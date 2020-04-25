package game

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Pnj struct {
	Character
	Talker
	TalkingTo      *Player
	IsPlayerFriend bool
}

type PnjConf struct {
	Level             string                `json:"level"`
	PosX              int                   `json:"posX"`
	PosY              int                   `json:"posY"`
	PosZ              int                   `json:"posZ"`
	Dead              bool                  `json:"dead"`
	Voice             string                `json:"voice"`
	CurrentNode       string                `json:"current_node"`
	Nodes             map[string]*StoryNode `json:"nodes"`
	Health            int                   `json:"health"`
	Energy            int                   `json:"energy"`
	Strength          int                   `json:"strength"`
	Dexterity         int                   `json:"dexterity"`
	Beauty            int                   `json:"beauty"`
	Will              int                   `json:"will"`
	Intelligence      int                   `json:"intelligence"`
	Charisma          int                   `json:"charisma"`
	Luck              int                   `json:"luck"`
	RegenerationSpeed int                   `json:"regeneration_speed"`
	Powers            []string              `json:"powers"`
}

func NewPnj(p Pos, name string) *Pnj {
	pnj := &Pnj{}
	pnj.Name = name
	pnj.Speed.Init(5)
	pnj.Health.Init(200)
	pnj.Energy.Init(200)
	pnj.Strength.Init(20)
	pnj.Dexterity.Init(20)
	pnj.Beauty.Init(20)
	pnj.Will.Init(20)
	pnj.Intelligence.Init(20)
	pnj.Charisma.Init(20)
	pnj.RegenerationSpeed.Init(1)
	pnj.Luck.Init(10)
	pnj.ActionPoints = 0.0
	pnj.Pos = p
	pnj.Xb = 0
	pnj.Yb = 0
	pnj.LastActionTime = time.Now()
	pnj.LookAt = Left
	pnj.Talkable = true
	pnj.Powers = make(map[string]*PlayerPower)

	return pnj
}

func (p *Pnj) LoadPnj(filename string) (string, Pos) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	conf := PnjConf{}
	json.Unmarshal(byteValue, &conf)

	p.Dead = conf.Dead
	if conf.Health != 0 {
		p.Health.Init(conf.Health)
	}
	if conf.Energy != 0 {
		p.Energy.Init(conf.Energy)
	}
	if conf.Strength != 0 {
		p.Strength.Init(conf.Strength)
	}
	if conf.Dexterity != 0 {
		p.Dexterity.Init(conf.Dexterity)
	}
	if conf.Beauty != 0 {
		p.Beauty.Init(conf.Beauty)
	}
	if conf.Will != 0 {
		p.Will.Init(conf.Will)
	}
	if conf.Intelligence != 0 {
		p.Intelligence.Init(conf.Intelligence)
	}
	if conf.Charisma != 0 {
		p.Charisma.Init(conf.Charisma)
	}
	if conf.Luck != 0 {
		p.Luck.Init(conf.Luck)
	}
	if conf.RegenerationSpeed != 0 {
		p.RegenerationSpeed.Init(conf.RegenerationSpeed)
	}
	powers := Powers()
	for _, pname := range conf.Powers {
		pow := powers.GetPower(pname)
		if pow == nil {
			log.Printf("Error: power %s does not exist", pname)
		} else {
			p.Powers[pname] = pow
		}
	}

	p.Dialog = &Dialog{
		CurrentNode: conf.CurrentNode,
		Nodes:       conf.Nodes,
	}

	return conf.Level, Pos{X: conf.PosX, Y: conf.PosY, Z: conf.PosZ}
}

func (pnj *Pnj) Talk(p *Player, g *Game) {
	EM.Dispatch(&Event{Action: ActionTalk, Payload: map[string]string{"voice": pnj.Voice}})
	pnj.Dialog.Init(p)
	node := pnj.Dialog.GetCurrentNode()
	node.ClearHighlight()
	node.SetHighlightedIndex(0)
	pnj.TalkingTo = p
	if p.X == pnj.X && p.Y < pnj.Y {
		pnj.LookAt = Up
	}
	if p.X == pnj.X && p.Y > pnj.Y {
		pnj.LookAt = Down
	}
	if p.Y == pnj.Y && p.X < pnj.X {
		pnj.LookAt = Left
	}
	if p.Y == pnj.Y && p.X > pnj.X {
		pnj.LookAt = Right
	}
}

func (pnj *Pnj) TalkChoiceUp() {
	node := pnj.Dialog.GetCurrentNode()
	choiceIdx := node.GetHighlightedIndex()
	node.SetHighlightedIndex(choiceIdx - 1)
}

func (pnj *Pnj) TalkChoiceDown() {
	node := pnj.Dialog.GetCurrentNode()
	choiceIdx := node.GetHighlightedIndex()
	node.SetHighlightedIndex(choiceIdx + 1)
}

func (pnj *Pnj) TalkConfirmChoice(g *Game) {
	node := pnj.Dialog.GetCurrentNode()
	choice := node.GetCurrentChoice()
	pnj.ChooseTalkOption(choice.Cmd, g)
}

func (pnj *Pnj) ChooseTalkOption(cmd string, g *Game) {
	node := pnj.Dialog.GetCurrentNode()
	nodeTo := pnj.Dialog.CurrentNode
	p := pnj.TalkingTo
	for _, choice := range node.Choices {
		if choice.Cmd == cmd {
			for _, questStep := range choice.Quest.StepsFullfilling {
				arr := strings.Split(questStep, ":")
				if len(arr) > 1 {
					questID := arr[0]
					stepID := arr[1]
					p.finishQuestStep(questID, stepID, g)
				}
			}
			for _, action := range choice.Actions {
				act := strings.Split(action, ":")
				switch act[0] {
				case "recruit":
					p.Recruit(pnj, g)
				case "teleport_to":
					pnj.Teleport(act[1], g)
				case "become_enemy":
					pnj.BecomeEnemy(g)
				case "set_current_node":
					pnj.Dialog.SetCurrentNode(act[1])
				case "send_to_level":
					levelPnj := strings.Split(act[1], "|")
					g.SendToLevel(levelPnj[0], levelPnj[1], levelPnj[2])
				case "update_dialog":
					levelPnjDialog := strings.Split(act[1], "|")
					g.UpdatePnjDialog(levelPnjDialog[0], levelPnjDialog[1], levelPnjDialog[2])
				case "learn_attack":
					p.LearnAttack(act[1])
				case "add_key":
					p.AddKey(act[1])
				case "gold_taken":
					v, err := strconv.Atoi(act[1])
					if err == nil {
						p.LooseGold(v)
					}
				case "discover":
					g.DiscoverLevel(act[1])
				case "book_given":
					book, exists := g.Books[act[1]]
					if exists {
						p.Library.AddBook(book)
						EM.Dispatch(&Event{
							Action:  ActionTake,
							Message: "You got a new book!",
						})
					} else {
						log.Printf("Book %s does not exist.", act[1])
					}
				}
			}
			if choice.NodeId == "" {
				pnj.StopTalking()
				return
			}
			nodeTo = choice.NodeId
		}
	}
	pnj.Dialog.CurrentNode = nodeTo
	EM.Dispatch(&Event{Action: ActionTalk, Payload: map[string]string{"voice": pnj.Voice}})
}

func (pnj *Pnj) StopTalking() {
	pnj.Dialog.Close()
	p := pnj.TalkingTo
	p.TalkingTo = nil
	pnj.TalkingTo = nil
}

func (pnj *Pnj) Update(l *Level) {
	if pnj.Dead {
		return
	}
	if pnj.TalkingTo != nil {
		return
	}
	t := time.Now()
	deltaD := t.Sub(pnj.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	pnj.ActionPoints += float64(pnj.Speed.Current) * delta
	pos := pnj.getWantedPosition()
	if pnj.ActionPoints >= 1000000 { // 1 second
		if pnj.canMove(pos, l) {
			pnj.Move(pos, l)
		}
		pnj.ActionPoints = 0.0
	}
	pnj.LastActionTime = time.Now()
}

func (pnj *Pnj) getWantedPosition() Pos {
	pos := pnj.Pos

	r := rand.Intn(5)
	switch r {
	case 1:
		pos.X++
	case 2:
		pos.X--
	case 3:
		pos.Y++
	case 4:
		pos.Y--
	default:
	}

	return pos
}

func (pnj *Pnj) canMove(to Pos, level *Level) bool {
	if pnj.TalkingTo != nil {
		return false
	}
	if !canGo(level, to) {
		return false
	}
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	if isThereAPortalAround(level, to) {
		return false
	}
	return true
}

func (pnj *Pnj) Move(to Pos, l *Level) {
	lastPos := pnj.Pos
	l.Map[pnj.Z][pnj.Y][pnj.X].Pnj = nil
	l.Map[to.Z][to.Y][to.X].Pnj = pnj
	pnj.Pos = to
	pnj.moveFromTo(lastPos, to)
}

func (pnj *Pnj) Teleport(levelName string, g *Game) {
	g.Level.MakeEffect(pnj.Pos, string(Teleport), 200)
	level := g.Levels[levelName]
	pnj.Talkable = false
	pnj.IsPowerUsing = true
	for pnj.PowerPos = 0; pnj.PowerPos < CaseLen; pnj.PowerPos++ {
		pnj.adaptSpeed()
	}
	pnj.ChangeLevel(g.Level, level)
	pnj.IsPowerUsing = false
	pnj.Talkable = true
}

func (pnj *Pnj) ChangeLevel(from *Level, to *Level) {
	from.Map[pnj.Z][pnj.Y][pnj.X].Pnj = nil
	if to != nil {
		pos := to.GetRandomFreePos(0) // FIXME
		pnj.Pos = *pos
		to.Map[pos.Z][pos.Y][pos.X].Pnj = pnj
	}
}

func (pnj *Pnj) BecomeEnemy(g *Game) {
	e := g.Level.MakeEnemy(pnj)
	g.Fight([]FighterInterface{e})
	if e.IsDead() {
		pnj.Die(g)
	}
}

func (pnj *Pnj) Die(g *Game) {
	pnj.Dead = true
	pnj.Dialog.SetCurrentNode("dead_greetings")
}
