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
	Friend            *Friend
	currentAttack     *Attack
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
		p.PowerUse(g)
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

func (p *Player) TakeDamages(damage int) {
	p.Character.TakeDamages(damage)
	p.Health.RaiseXp(damage)
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

func (p *Player) ChooseAction(ring *FightingRing) int {
	switch ring.SelectedPlayerAction {
	case "run":
		return p.Speed.Current
	case "attack":
		att := ring.PossibleAttacks.List[ring.PossibleAttacks.Selected]
		p.currentAttack = att
		return att.Speed
	}
	return 0
}

func (p *Player) Fight(ring *FightingRing) {
	switch ring.SelectedPlayerAction {
	case "run":
		ring.End()
	case "attack":
		att := p.currentAttack
		var to []FighterInterface
		idx := 0
		for i := ring.TargetSelected; idx < att.Range && i < len(ring.Enemies); i++ {
			f := ring.Enemies[i]
			if !f.IsDead() {
				to = append(to, f)
				idx++
			}
		}

		p.isAttacking = true
		for p.AttackPos = 0; p.AttackPos < CaseLen; p.AttackPos++ {
			att.adaptSpeed()
		}
		p.isAttacking = false

		switch att.Type {
		case AttackTypePhysical:
			damages := att.Damages * p.CalculateAttackScore() / 10
			for _, f := range to {
				f.TakeDamages(damages)
			}
			p.Strength.RaiseXp(damages * len(to) / 10)
			p.Dexterity.RaiseXp(1)
		case AttackTypeMagick:
			switch att.MagickType {
			case PowerHealing:
				p.Health.Add(att.Damages * p.CalculatePowerAttackScore() / 10)
			case PowerInvocation:
				// TODO
			default:
				damages := att.Damages * p.CalculatePowerAttackScore() / 10
				for _, f := range to {
					f.TakeDamages(damages)
				}
				p.Energy.RaiseXp(damages)
				p.Will.RaiseXp(damages * len(to) / 10)
			}
			p.LooseEnergy(att.EnergyCost)
			p.Intelligence.RaiseXp(1)
		}
	}
}
