package game

import "time"
import "math/rand"
import "strings"

type Pnj struct {
	Character
	Talker
	TalkingTo      *Player
	IsPlayerFriend bool
}

func NewPnj(p Pos, name string, voice string) *Pnj {
	pnj := &Pnj{}
	pnj.Name = name
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
	pnj.Speed.Init(4)
	pnj.ActionPoints = 0.0
	pnj.Pos = p
	pnj.Xb = 0
	pnj.Yb = 0
	pnj.LastActionTime = time.Now()
	pnj.LookAt = Left
	pnj.Talkable = true
	pnj.IsTalking = false
	pnj.Voice = voice

	return pnj
}

func (pnj *Pnj) Talk(p *Player, g *Game) {
	EM.Dispatch(&Event{Action: ActionTalk, Payload: map[string]string{"voice": pnj.Voice}})
	pnj.Dialog.Init(p)
	node := pnj.Dialog.GetCurrentNode()
	node.ClearHighlight()
	node.SetHighlightedIndex(0)
	pnj.IsTalking = true
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
				case "set_initial_node":
					pnj.Dialog.SetInitialNode(act[1])
				case "send_to_level":
					levelPnj := strings.Split(act[1], "|")
					g.SendToLevel(levelPnj[0], levelPnj[1], levelPnj[2])
				}
			}
			for _, book := range choice.BooksGiven {
				p.Library.AddBook(g.Books[book])
				EM.Dispatch(&Event{
					Action:  ActionTake,
					Message: "You got a new book!",
				})
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
	p := pnj.TalkingTo
	p.IsTalking = false
	p.TalkingTo = nil
	pnj.IsTalking = false
	pnj.TalkingTo = nil
	for k, node := range pnj.Dialog.Nodes {
		if node.Initial {
			pnj.Dialog.CurrentNode = k
			break
		}
	}
}

func (pnj *Pnj) Update(l *Level) {
	if pnj.Dead {
		return
	}
	if pnj.IsTalking {
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
	pos := Pos{}
	pos.X = pnj.X
	pos.Y = pnj.Y

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
	if !canGo(level, to) {
		return false
	}
	if to.X == level.Player.X && to.Y == level.Player.Y {
		return false
	}
	return true
}

func (pnj *Pnj) Move(to Pos, l *Level) {
	lastPos := Pos{X: pnj.Pos.X, Y: pnj.Pos.Y}
	l.Map[pnj.Y][pnj.X].Pnj = nil
	l.Map[to.Y][to.X].Pnj = pnj
	pnj.Pos = to
	pnj.moveFromTo(lastPos, to)
}

func (pnj *Pnj) Teleport(levelName string, g *Game) {
	g.Level.MakeEffect(pnj.Pos, rune(Teleport), 200)
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
	from.Map[pnj.Y][pnj.X].Pnj = nil
	pos := to.GetRandomFreePos()
	pnj.Pos = *pos
	to.Map[pos.Y][pos.X].Pnj = pnj
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
	pnj.Dialog.SetInitialNode("dead_greetings")
}
