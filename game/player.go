package game

import "math/rand"
import "time"

type Player struct {
	Character
	Talker
	TalkingTo         *Pnj
	Quests            map[string]*Quest
	Inventory         *Inventory
	Library           *Library
	IsTaking          bool
	Menu              *Menu
	QuestMenuOpen     bool
	CharacterMenuOpen bool
	MapMenuOpen       bool
	Weapons           []*Weapon
	Friend            *Friend
}

func (p *Player) Update(g *Game) {
	if p.isDead {
		return
	}
	input := g.GetInput()
	p.regenerate()
	if p.IsTalking {
		p.Discuss(g)
	} else {
		if input.Typ == Up || input.Typ == Down || input.Typ == Left || input.Typ == Right {
			p.LookAt = input.Typ
		}
		if g.GetInput2().Typ == SpeedUp {
			p.Speed.Current = p.Speed.Initial * 2.0
		} else if g.GetInput2().Typ == None {
			p.Speed.Current = p.Speed.Initial
		}
		p.Move(g)
	}
}

func (p *Player) Move(g *Game) {
	input := g.GetInput()
	level := g.Level
	switch input.Typ {
	case Up:
		openDoor(g, Pos{p.X, p.Y - 1})
		if canGo(level, Pos{p.X, p.Y - 1}) {
			p.beforeMovingActions(g)
			p.WalkUp()
			p.afterMovingActions(g)
		}
	case Down:
		openDoor(g, Pos{p.X, p.Y + 1})
		if canGo(level, Pos{p.X, p.Y + 1}) {
			p.beforeMovingActions(g)
			p.WalkDown()
			p.afterMovingActions(g)
		}
	case Left:
		openDoor(g, Pos{p.X - 1, p.Y})
		if canGo(level, Pos{p.X - 1, p.Y}) {
			p.beforeMovingActions(g)
			p.WalkLeft()
			p.afterMovingActions(g)
		}
	case Right:
		openDoor(g, Pos{p.X + 1, p.Y})
		if canGo(level, Pos{p.X + 1, p.Y}) {
			p.beforeMovingActions(g)
			p.WalkRight()
			p.afterMovingActions(g)
		}
	case Action:
		posTo := Pos{p.X, p.Y + 1}
		switch p.LookAt {
		case Up:
			posTo = Pos{p.X, p.Y - 1}
		case Down:
			posTo = Pos{p.X, p.Y + 1}
		case Left:
			posTo = Pos{p.X - 1, p.Y}
		case Right:
			posTo = Pos{p.X + 1, p.Y}
		}
		p.Talk(g, posTo)
		p.Take(g, posTo)
	case Power:
		// TODO
	default:
	}
}

func (p *Player) beforeMovingActions(g *Game) {
	p.DispatchWalkingEvent(g)
}

func (p *Player) afterMovingActions(g *Game) {
	p.DispatchWalkingEvent(g)
	p.openPortal(g, p.Pos)
	p.MeetMonsters(g)
}

func (p *Player) DispatchWalkingEvent(g *Game) {
	g.GetEventManager().Dispatch(&Event{Action: ActionWalk})
}

func (p *Player) WalkDown() {
	p.Y++
	p.Yb = CaseLen
	p.moveDown()
}

func (p *Player) WalkUp() {
	p.Y--
	p.Yb = -1 * CaseLen
	p.moveUp()
}

func (p *Player) WalkLeft() {
	p.X--
	p.Xb = -1 * CaseLen
	p.moveLeft()
}

func (p *Player) WalkRight() {
	p.X++
	p.Xb = CaseLen
	p.moveRight()
}

func (p *Player) TakeDamages(damage int) {
	if p.Health.Current <= 0 {
		p.isDead = true
		return
	}
	p.Health.Current -= damage
	p.Health.RaiseXp(damage)
}

func (p *Player) openPortal(g *Game, pos Pos) {
	g.Level.OpenPortal(g, pos)
}

func (p *Player) Talk(g *Game, posTo Pos) {
	level := g.Level
	pnj := level.Map[posTo.Y][posTo.X].Pnj
	if pnj != nil && pnj.Talkable {
		g.GetEventManager().Dispatch(&Event{Action: ActionTalk})
		p.IsTalking = true
		p.TalkingTo = pnj
		pnj.Talk(p, g)
		adaptDialogSpeed()
	}
}

func (p *Player) Discuss(g *Game) {
	input := g.GetInput()
	pnj := p.TalkingTo
	if pnj == nil {
		p.IsTalking = false
		return
	}
	switch input.Typ {
	case Up:
		pnj.TalkChoiceUp()
		adaptDialogSpeed()
	case Down:
		pnj.TalkChoiceDown()
		adaptDialogSpeed()
	case Action:
		g.GetEventManager().Dispatch(&Event{Action: ActionTalk})
		pnj.TalkConfirmChoice(g)
		adaptDialogSpeed()
	default:
	}
}

func (p *Player) IsQuestOpen(questID string) bool {
	return !p.Quests[questID].IsFinished
}

func (p *Player) IsStepFinished(questID string, stepID string) bool {
	if p.Quests[questID].IsFinished {
		return true
	}
	return p.Quests[questID].Steps[stepID].IsFinished
}

func (p *Player) IsQuestOpenStepFinished(questID string, stepID string) bool {
	if p.Quests[questID].IsFinished {
		return false
	}
	return p.Quests[questID].Steps[stepID].IsFinished
}

func (p *Player) Take(g *Game, posTo Pos) bool {
	level := g.Level
	o := level.Map[posTo.Y][posTo.X].Object
	if o != nil {
		p.IsTaking = true
		ut := p.TakeUsable(o, g)
		bt := p.TakeBook(o, g)
		qot := p.TakeQuestObject(o, g)
		taken := ut || bt || qot
		if taken {
			for i := CaseLen; i > 0; i = i - 1 {
				p.adaptSpeed()
			}
		}
		p.IsTaking = false
		return taken
	}
	return false
}

func (p *Player) TakeQuestObject(o *Object, g *Game) bool {
	qo, isQuestObject := g.QuestsObjects[o.Rune]
	if !isQuestObject {
		return false
	}
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionTake,
		Message: "You got a special object!",
	})
	p.Inventory.QuestObjects[o.Rune] = o
	g.Level.Map[o.Y][o.X].Object = nil

	for _, stepID := range qo.Quest.StepsFullfilling {
		p.finishQuestStep(qo.Quest.ID, stepID, g)
	}

	return true
}

func (p *Player) TakeUsable(o *Object, g *Game) bool {
	taken := p.Inventory.TakeUsable(o)
	if taken {
		g.GetEventManager().Dispatch(&Event{Action: ActionTake})
		g.Level.Map[o.Y][o.X].Object = nil
	}

	return taken
}

func (p *Player) TakeBook(o *Object, g *Game) bool {
	taken := p.AddBook(o, g)
	if taken {
		g.GetEventManager().Dispatch(&Event{
			Action:  ActionTake,
			Message: "You got a new book!",
		})
		g.Level.Map[o.Y][o.X].Object = nil
	}

	return taken
}

func (p *Player) Recruit(pnj *Pnj, g *Game) {
	if p.Friend != nil {
		g.GetEventManager().Dispatch(&Event{
			Message: "You already have a friend, you can't recruit.",
		})
		return
	}
	f := g.Level.MakeFriend(pnj)
	p.Friend = f
}

func (p *Player) MeetMonsters(g *Game) {
	// TODO : handle monsters types by level case type
	l := g.Level
	r := rand.Intn(100)
	if r > l.MonstersProbability {
		return
	}
	switch l.Type {
	case LevelTypeOutdoor:
		bestiary := Bestiary()
		g.FightMonsters(bestiary)
	case LevelTypeGrotto:
		bestiary := BestiaryUnderworld()
		g.FightMonsters(bestiary)
	}
}

func (p *Player) Fight(ring *FightingRing) AttackInterface {
	// TODO
	time.Sleep(1 * time.Second)
	to := ring.Enemies[0]
	i := 0
	for to.IsDead() && i < len(ring.Enemies) {
		to = ring.Enemies[i]
		i++
	}
	if to.IsDead() {
		return nil
	}
	att := &SwordAttack{
		From:    p,
		To:      to,
		Speed:   p.Weapon.Speed,
		Damages: p.CalculateAttackScore(),
	}
	return att
}
