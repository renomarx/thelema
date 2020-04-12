package game

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
	if p.Dead {
		return
	}
	input := g.GetInput()
	p.regenerate()
	if p.TalkingTo != nil {
		p.Discuss(g)
	} else {
		if input.Typ == Up || input.Typ == Down || input.Typ == Left || input.Typ == Right {
			p.LookAt = input.Typ
		}
		switch g.GetInput2().Typ {
		case SpeedUp:
			p.Speed.Current = p.Speed.Initial * 2
		case Shadow:
			p.Shadow = true
			p.Speed.Current = p.Speed.Initial / 2
		case Meditate:
			p.Meditating = true
			p.RegenerationSpeed.Current = p.RegenerationSpeed.Initial * 5
		case None:
			p.Shadow = false
			p.Meditating = false
			p.Speed.Current = p.Speed.Initial
			p.RegenerationSpeed.Current = p.RegenerationSpeed.Initial
		}
		if !p.Meditating {
			p.Move(g)
		}
	}
}

func (p *Player) Move(g *Game) {
	input := g.GetInput()
	level := g.Level
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
	switch input.Typ {
	case Up:
		if canGo(level, Pos{p.X, p.Y - 1}) {
			p.beforeMovingActions(g)
			p.WalkUp()
			p.afterMovingActions(g)
		}
	case Down:
		if canGo(level, Pos{p.X, p.Y + 1}) {
			p.beforeMovingActions(g)
			p.WalkDown()
			p.afterMovingActions(g)
		}
	case Left:
		if canGo(level, Pos{p.X - 1, p.Y}) {
			p.beforeMovingActions(g)
			p.WalkLeft()
			p.afterMovingActions(g)
		}
	case Right:
		if canGo(level, Pos{p.X + 1, p.Y}) {
			p.beforeMovingActions(g)
			p.WalkRight()
			p.afterMovingActions(g)
		}
	case Action:
		p.Talk(g, posTo)
		p.Take(g, posTo)
	case Power:
		p.PowerUse(g, posTo)
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
	EM.Dispatch(&Event{Action: ActionWalk})
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
	if pnj != nil && pnj.Talkable && !pnj.IsDead() {
		EM.Dispatch(&Event{Action: ActionTalk})
		p.TalkingTo = pnj
		pnj.Talk(p, g)
		adaptDialogSpeed()
	}
}

func (p *Player) Discuss(g *Game) {
	input := g.GetInput()
	pnj := p.TalkingTo
	if pnj == nil {
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
		EM.Dispatch(&Event{Action: ActionTalk})
		pnj.TalkConfirmChoice(g)
		adaptDialogSpeed()
		p.Intelligence.RaiseXp(1)
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
			for i := 32; i > 0; i = i - 1 {
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
	EM.Dispatch(&Event{
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
		EM.Dispatch(&Event{Action: ActionTake})
		g.Level.Map[o.Y][o.X].Object = nil
	}

	return taken
}

func (p *Player) TakeBook(o *Object, g *Game) bool {
	taken := p.AddBook(o, g)
	if taken {
		EM.Dispatch(&Event{
			Action:  ActionTake,
			Message: "You got a new book!",
		})
		g.Level.Map[o.Y][o.X].Object = nil
		p.Intelligence.RaiseXp(10)
	}

	return taken
}

func (p *Player) Recruit(pnj *Pnj, g *Game) {
	if p.Friend != nil {
		EM.Dispatch(&Event{
			Message: "You already have a friend, you can't recruit.",
		})
		return
	}
	f := g.Level.MakeFriend(pnj)
	p.Friend = f
}

func (c *Player) PowerUse(g *Game, posTo Pos) {
	if c.Energy.Current > 0 {
		c.IsPowerUsing = true
		for c.PowerPos = 0; c.PowerPos < 32; c.PowerPos++ {
			c.CurrentPower.adaptSpeed()
		}
		switch c.CurrentPower.Type {
		case PowerInvocation:
			// TODO : make input control invocation for lifetime
		case PowerCalm:
			g.Level.MakeEffect(posTo, string(Calm), 200)
			// TODO
		case PowerDeadSpeaking:
			g.Level.MakeEffect(posTo, string(Necromancy), 200)
			c.TalkToDead(g, posTo)
		case PowerStorm:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerStorm}})
			g.Level.MakeRangeStorm(c.Pos, c.CalculatePowerAttackScore(), c.LookAt, 1, 10)
			c.LooseEnergy(c.CurrentPower.Energy)
		case PowerFlames:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerFlames}})
			g.Level.MakeFlames(c.Pos, c.CalculatePowerAttackScore(), 1, 5)
			c.LooseEnergy(c.CurrentPower.Energy)
		case PowerHealing:
			EM.Dispatch(&Event{Action: ActionPower, Payload: map[string]string{"type": PowerHealing}})
			g.Level.MakeEffect(c.Pos, string(Healing), 200)
			c.Health.Add(c.CalculatePowerAttackScore())
			c.LooseEnergy(c.CurrentPower.Energy)
		default:
		}
		c.IsPowerUsing = false
	}
}

func (p *Player) TalkToDead(g *Game, posTo Pos) {
	level := g.Level
	pnj := level.Map[posTo.Y][posTo.X].Pnj
	if pnj != nil && pnj.Talkable && pnj.IsDead() {
		EM.Dispatch(&Event{Action: ActionTalk})
		p.TalkingTo = pnj
		pnj.Talk(p, g)
		adaptDialogSpeed()
	}
}

func (p *Player) LearnAttack(attackName string) {
	for _, att := range p.Attacks {
		if att.Name == attackName {
			return
		}
	}
	for _, att := range Attacks() {
		if att.Name == attackName {
			p.Attacks = append(p.Attacks, att)
			EM.Dispatch(&Event{
				Message: "You learned a new attack: " + att.Name + " !",
			})
		}
	}
}

func (p *Player) AddKey(key string) {
	new := p.Inventory.AddKey(key)
	if new {
		EM.Dispatch(&Event{
			Message: "You got a new key: " + key + " !",
		})
	}
}

func (p *Player) LooseGold(value int) {
	p.Inventory.Gold -= value
}
