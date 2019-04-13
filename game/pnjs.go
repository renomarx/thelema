package game

import "time"
import "math/rand"

type Pnj struct {
	Character
	Talker
	TalkingTo *Player
}

func NewPnj(p Pos, r rune, name string) *Pnj {
	pnj := &Pnj{}
	pnj.Rune = r
	pnj.Name = name
	pnj.Hitpoints.Init(50)
	pnj.Strength.Init(5)
	pnj.Speed.Init(2)
	pnj.ActionPoints = 0.0
	pnj.Pos = p
	pnj.Xb = 0
	pnj.Yb = 0
	pnj.LastActionTime = time.Now()
	pnj.IsMoving = false
	pnj.LookAt = Left
	pnj.IsTalking = false

	return pnj
}

func (pnj *Pnj) Talk(p *Player) {
	pnj.Dialog.Init(p)
	node := pnj.Dialog.GetCurrentNode()
	node.ClearHighlight()
	node.SetHighlightedIndex(0)
	pnj.IsTalking = true
	pnj.TalkingTo = p
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

func (pnj *Pnj) TalkConfirmChoice() {
	node := pnj.Dialog.GetCurrentNode()
	choice := node.GetCurrentChoice()
	pnj.ChooseTalkOption(choice.Cmd)
}

func (pnj *Pnj) ChooseTalkOption(cmd string) {
	node := pnj.Dialog.GetCurrentNode()
	nodeTo := pnj.Dialog.CurrentNode
	p := pnj.TalkingTo
	for _, choice := range node.Choices {
		if choice.Cmd == cmd {
			for _, stepID := range choice.Quest.StepsFullfilling {
				p.finishQuestStep(choice.Quest.ID, stepID)
			}
			if choice.NodeId == "" {
				pnj.StopTalking()
				return
			}
			nodeTo = choice.NodeId
		}
	}
	pnj.Dialog.CurrentNode = nodeTo
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

func (pnj *Pnj) Update(game *Game) {
	if pnj.IsMoving || pnj.IsTalking {
		return
	}
	level := game.Level
	t := time.Now()
	deltaD := t.Sub(pnj.LastActionTime)
	delta := 0.001 * float64(deltaD.Nanoseconds())
	pnj.ActionPoints += float64(pnj.Speed.Current) * delta
	pos := pnj.getWantedPosition()
	if pnj.ActionPoints >= 1000000 { // 1 second
		if pnj.canMove(pos, level) {
			pnj.Move(pos, level)
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
	if level.Map[to.Y][to.X] == DoorClosed {
		return false
	}
	return true
}

func (pnj *Pnj) Move(to Pos, level *Level) {
	pnj.IsMoving = true
	lastPos := Pos{X: pnj.Pos.X, Y: pnj.Pos.Y}
	Mux.Lock()
	delete(level.Pnjs, pnj.Pos)
	level.Pnjs[to] = pnj
	Mux.Unlock()
	pnj.moveFromTo(lastPos, to)
}
