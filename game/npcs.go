package game

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Npc struct {
	Character
	Talker
	TalkingTo      *Player
	IsPlayerFriend bool
}

type NpcConf struct {
	Level             string                `yaml:"level"`
	PosX              int                   `yaml:"posX"`
	PosY              int                   `yaml:"posY"`
	PosZ              int                   `yaml:"posZ"`
	Dead              bool                  `yaml:"dead"`
	Voice             string                `yaml:"voice"`
	CurrentNode       string                `yaml:"current_node"`
	Nodes             map[string]*StoryNode `yaml:"nodes"`
	Health            int                   `yaml:"health"`
	Energy            int                   `yaml:"energy"`
	Strength          int                   `yaml:"strength"`
	Dexterity         int                   `yaml:"dexterity"`
	Beauty            int                   `yaml:"beauty"`
	Will              int                   `yaml:"will"`
	Intelligence      int                   `yaml:"intelligence"`
	Charisma          int                   `yaml:"charisma"`
	Luck              int                   `yaml:"luck"`
	Aggressiveness    int                   `yaml:"aggressiveness"`
	RegenerationSpeed int                   `yaml:"regeneration_speed"`
	Powers            []string              `yaml:"powers"`
}

func NewNpc(p Pos, name string) *Npc {
	npc := &Npc{}
	npc.Character = NewCharacter()
	npc.Name = name
	npc.Speed.Init(5)
	npc.Health.Init(200)
	npc.Energy.Init(200)
	npc.Strength.Init(20)
	npc.Dexterity.Init(20)
	npc.Beauty.Init(20)
	npc.Will.Init(20)
	npc.Intelligence.Init(20)
	npc.Charisma.Init(20)
	npc.RegenerationSpeed.Init(1)
	npc.Luck.Init(10)
	npc.Aggressiveness.Init(200)
	npc.ActionPoints = 0.0
	npc.Pos = p
	npc.Xb = 0
	npc.Yb = 0
	npc.LastActionTime = time.Now()
	npc.LookAt = Left
	npc.Talkable = true

	return npc
}

func (c *Npc) LoadNpc(filename string) (string, Pos) {
	yamlFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()

	byteValue, _ := ioutil.ReadAll(yamlFile)

	conf := NpcConf{}
	yaml.Unmarshal(byteValue, &conf)

	c.Dead = conf.Dead
	if conf.Health != 0 {
		c.Health.Init(conf.Health)
	}
	if conf.Energy != 0 {
		c.Energy.Init(conf.Energy)
	}
	if conf.Strength != 0 {
		c.Strength.Init(conf.Strength)
	}
	if conf.Dexterity != 0 {
		c.Dexterity.Init(conf.Dexterity)
	}
	if conf.Beauty != 0 {
		c.Beauty.Init(conf.Beauty)
	}
	if conf.Will != 0 {
		c.Will.Init(conf.Will)
	}
	if conf.Intelligence != 0 {
		c.Intelligence.Init(conf.Intelligence)
	}
	if conf.Charisma != 0 {
		c.Charisma.Init(conf.Charisma)
	}
	if conf.Luck != 0 {
		c.Luck.Init(conf.Luck)
	}
	if conf.Aggressiveness != 0 {
		c.Aggressiveness.Init(conf.Aggressiveness)
	}
	if conf.RegenerationSpeed != 0 {
		c.RegenerationSpeed.Init(conf.RegenerationSpeed)
	}
	powers := Powers()
	for _, pname := range conf.Powers {
		pow := powers.GetPower(pname)
		if pow == nil {
			log.Printf("Error: power %s does not exist", pname)
		} else {
			c.Powers[pname] = pow
		}
	}

	c.Dialog = &Dialog{
		CurrentNode: conf.CurrentNode,
		Nodes:       conf.Nodes,
	}

	return conf.Level, Pos{X: conf.PosX, Y: conf.PosY, Z: conf.PosZ}
}

func (npc *Npc) Talk(p *Player, g *Game) {
	EM.Dispatch(&Event{Action: ActionTalk, Payload: map[string]string{"voice": npc.Voice}})
	npc.Dialog.Init(g)
	npc.TalkingTo = p
	if p.X == npc.X && p.Y < npc.Y {
		npc.LookAt = Up
	}
	if p.X == npc.X && p.Y > npc.Y {
		npc.LookAt = Down
	}
	if p.Y == npc.Y && p.X < npc.X {
		npc.LookAt = Left
	}
	if p.Y == npc.Y && p.X > npc.X {
		npc.LookAt = Right
	}
	npc.doTalk(g)
}

func (npc *Npc) doTalk(g *Game) {
	EM.Dispatch(&Event{Action: ActionTalk, Payload: map[string]string{"voice": npc.Voice}})
	node := npc.Dialog.GetCurrentNode()
	node.filterPossibleChoices(g)
	node.ClearHighlight()
	node.SetHighlightedIndex(0)
}

func (npc *Npc) TalkChoiceUp() {
	node := npc.Dialog.GetCurrentNode()
	choiceIdx := node.GetHighlightedIndex()
	node.SetHighlightedIndex(choiceIdx - 1)
}

func (npc *Npc) TalkChoiceDown() {
	node := npc.Dialog.GetCurrentNode()
	choiceIdx := node.GetHighlightedIndex()
	node.SetHighlightedIndex(choiceIdx + 1)
}

func (npc *Npc) TalkConfirmChoice(g *Game) {
	node := npc.Dialog.GetCurrentNode()
	choice := node.GetCurrentChoice()
	npc.ChooseTalkOption(choice.Cmd, g)
}

func (npc *Npc) ChooseTalkOption(cmd string, g *Game) {
	node := npc.Dialog.GetCurrentNode()
	nodeTo := npc.Dialog.CurrentNode
	p := npc.TalkingTo
	for _, choice := range node.Choices {
		if choice.Cmd == cmd {
			for _, stID := range choice.StepsBeginning {
				g.beginStep(stID)
			}
			for _, stID := range choice.StepsFinishing {
				g.finishStep(stID)
			}
			for _, stID := range choice.StepsCanceling {
				g.cancelStep(stID)
			}
			for _, action := range choice.Actions {
				act := strings.Split(action, ":")
				switch act[0] {
				case "recruit":
					p.Recruit(npc, g)
				case "teleport_to":
					npc.Teleport(act[1], g)
				case "become_enemy":
					npc.BecomeEnemy(g)
				case "set_current_node":
					npc.Dialog.SetCurrentNode(act[1])
				case "send_to_level":
					levelNpc := strings.Split(act[1], "|")
					g.SendToLevel(levelNpc[0], levelNpc[1], levelNpc[2])
				case "update_dialog":
					levelNpcDialog := strings.Split(act[1], "|")
					g.UpdateNpcDialog(levelNpcDialog[0], levelNpcDialog[1], levelNpcDialog[2])
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
							Message: "Vous avez un nouveau livre!",
						})
					} else {
						log.Printf("Book %s does not exist.", act[1])
					}
				}
			}
			if choice.NodeId == "" {
				npc.StopTalking()
				return
			}
			nodeTo = choice.NodeId
		}
	}
	npc.Dialog.CurrentNode = nodeTo
	npc.doTalk(g)
}

func (npc *Npc) StopTalking() {
	npc.Dialog.Close()
	p := npc.TalkingTo
	p.TalkingTo = nil
	npc.TalkingTo = nil
}

func (npc *Npc) Update(l *Level) {
	if npc.Dead {
		return
	}
	if npc.TalkingTo != nil {
		return
	}
	t := time.Now()
	deltaD := t.Sub(npc.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	npc.ActionPoints += float64(npc.Speed.Current) * delta
	pos := npc.getWantedPosition()
	if npc.ActionPoints >= 1000000 { // 1 second
		if npc.canMove(pos, l) {
			npc.Move(pos, l)
		}
		npc.ActionPoints = 0.0
	}
	npc.LastActionTime = time.Now()
}

func (npc *Npc) getWantedPosition() Pos {
	pos := npc.Pos

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

func (npc *Npc) canMove(to Pos, level *Level) bool {
	if npc.TalkingTo != nil {
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

func (npc *Npc) Move(to Pos, l *Level) {
	lastPos := npc.Pos
	l.Map[npc.Z][npc.Y][npc.X].Npc = nil
	l.Map[to.Z][to.Y][to.X].Npc = npc
	npc.Pos = to
	npc.moveFromTo(lastPos, to)
}

func (npc *Npc) Teleport(levelName string, g *Game) {
	g.Level.MakeEffect(npc.Pos, string(Teleport), 200)
	level := g.Levels[levelName]
	npc.Talkable = false
	npc.IsPowerUsing = true
	for npc.PowerUsingStage = 0; npc.PowerUsingStage < CaseLen; npc.PowerUsingStage++ {
		npc.adaptSpeed()
	}
	npc.ChangeLevel(g.Level, level)
	npc.IsPowerUsing = false
	npc.Talkable = true
}

func (npc *Npc) ChangeLevel(from *Level, to *Level) {
	from.Map[npc.Z][npc.Y][npc.X].Npc = nil
	if to != nil {
		pos := to.GetRandomFreePos(0) // FIXME
		npc.Pos = *pos
		to.Map[pos.Z][pos.Y][pos.X].Npc = npc
	}
}

func (npc *Npc) BecomeEnemy(g *Game) {
	e := g.Level.MakeEnemy(npc)
	g.Fight([]FighterInterface{e})
	if e.IsDead() {
		npc.Die(g)
	}
}

func (npc *Npc) Die(g *Game) {
	npc.Dead = true
	npc.Dialog.SetCurrentNode("dead_greetings")
}
