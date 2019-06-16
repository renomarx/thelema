package game

import "math/rand"

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
			p.DispatchWalkingEvent(g)
			p.WalkUp()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Down:
		openDoor(g, Pos{p.X, p.Y + 1})
		if canGo(level, Pos{p.X, p.Y + 1}) {
			p.DispatchWalkingEvent(g)
			p.WalkDown()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Left:
		openDoor(g, Pos{p.X - 1, p.Y})
		if canGo(level, Pos{p.X - 1, p.Y}) {
			p.DispatchWalkingEvent(g)
			p.WalkLeft()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Right:
		openDoor(g, Pos{p.X + 1, p.Y})
		if canGo(level, Pos{p.X + 1, p.Y}) {
			p.DispatchWalkingEvent(g)
			p.WalkRight()
			p.openPortal(g, Pos{p.X, p.Y})
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
		if !p.IsTalking && !p.IsTaking {
			if p.Attack(g, posTo) {
				g.GetEventManager().Dispatch(&Event{Action: ActionAttack})
			}
		}
	case Power:
		p.PowerAttack(g)
	default:
	}
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

func (p *Player) TakeDamage(g *Game, damage int) {
	if p.Health.Current <= 0 {
		p.Die(g)
		return
	}
	p.Health.Current -= damage
	g.MakeExplosion(p.Pos, damage, 50)
	p.Health.RaiseXp(damage, g)

	g.GetEventManager().Dispatch(&Event{
		Action: ActionHurt})
	p.ParalyzedTime = rand.Intn(damage) * 10
}

func (p *Player) Die(g *Game) {
	if p.isDead {
		return
	}
	p.isDead = true
	g.GetMenu().Choices[1].Disabled = true
	g.GetEventManager().Dispatch(&Event{
		Action:  ActionDie,
		Message: "You're dead !"})
}

func (p *Player) openPortal(g *Game, pos Pos) {
	g.Level.OpenPortal(g, pos)
}

func (p *Player) Talk(g *Game, posTo Pos) {
	level := g.Level
	pnj := level.Pnjs[posTo.Y][posTo.X]
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

func (p *Player) IsQuestOpenStepFinished(questID string, stepID string) bool {
	if p.Quests[questID].IsFinished {
		return false
	}
	return p.Quests[questID].Steps[stepID].IsFinished
}

func (p *Player) Take(g *Game, posTo Pos) {
	level := g.Level
	o := level.Objects[posTo.Y][posTo.X]
	if o != nil {
		p.IsTaking = true
		p.TakeUsable(o, g)
		p.TakeBook(o, g)
		p.TakeQuestObject(o, g)
		for i := CaseLen; i > 0; i = i - 2 {
			p.adaptSpeed()
		}
		p.IsTaking = false
	}
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
	g.Mux.Lock()
	p.Inventory.QuestObjects[o.Rune] = o
	g.Mux.Unlock()
	g.Level.Objects[o.Y][o.X] = nil

	for _, stepID := range qo.Quest.StepsFullfilling {
		p.finishQuestStep(qo.Quest.ID, stepID, g)
	}

	return true
}

func (p *Player) TakeUsable(o *Object, g *Game) bool {
	taken := p.Inventory.TakeUsable(o)
	if taken {
		g.GetEventManager().Dispatch(&Event{Action: ActionTake})
		g.Level.Objects[o.Y][o.X] = nil
	}

	return true
}

func (p *Player) TakeBook(o *Object, g *Game) bool {
	taken := p.AddBook(o, g)
	if taken {
		g.GetEventManager().Dispatch(&Event{
			Action:  ActionTake,
			Message: "You got a new book!",
		})
		g.Level.Objects[o.Y][o.X] = nil
	}

	return true
}

func (p *Player) Recruit(pnj *Pnj, g *Game) {
	g.Level.Pnjs[pnj.Y][pnj.X] = nil
	g.Level.MakeFriend(pnj)
}
