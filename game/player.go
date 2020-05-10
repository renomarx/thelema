package game

type Player struct {
	Character
	Talker
	TalkingTo         *Npc
	Inventory         *Inventory
	Library           *Library
	IsTaking          bool
	Menu              *Menu
	QuestMenuOpen     bool
	CharacterMenuOpen bool
	MapMenuOpen       bool
	Friend            *Friend
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
	posTo := Pos{X: p.X, Y: p.Y + 1, Z: p.Z}
	switch p.LookAt {
	case Up:
		posTo = Pos{X: p.X, Y: p.Y - 1, Z: p.Z}
	case Down:
		posTo = Pos{X: p.X, Y: p.Y + 1, Z: p.Z}
	case Left:
		posTo = Pos{X: p.X - 1, Y: p.Y, Z: p.Z}
	case Right:
		posTo = Pos{X: p.X + 1, Y: p.Y, Z: p.Z}
	}
	switch input.Typ {
	case Up:
		if canGo(level, Pos{X: p.X, Y: p.Y - 1, Z: p.Z}) {
			p.beforeMovingActions(g)
			p.WalkUp()
			p.afterMovingActions(g)
		}
	case Down:
		if canGo(level, Pos{X: p.X, Y: p.Y + 1, Z: p.Z}) {
			p.beforeMovingActions(g)
			p.WalkDown()
			p.afterMovingActions(g)
		}
	case Left:
		if canGo(level, Pos{X: p.X - 1, Y: p.Y, Z: p.Z}) {
			p.beforeMovingActions(g)
			p.WalkLeft()
			p.afterMovingActions(g)
		}
	case Right:
		if canGo(level, Pos{X: p.X + 1, Y: p.Y, Z: p.Z}) {
			p.beforeMovingActions(g)
			p.WalkRight()
			p.afterMovingActions(g)
		}
	case Action:
		p.Talk(g, posTo)
		p.Take(g, posTo)
	case Action2:
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
	g.MeetMonsters()
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
	npc := level.Map[posTo.Z][posTo.Y][posTo.X].Npc
	if npc != nil && npc.Talkable && !npc.IsDead() {
		EM.Dispatch(&Event{Action: ActionTalk})
		p.TalkingTo = npc
		npc.Talk(p, g)
		adaptDialogSpeed()
	}
}

func (p *Player) Discuss(g *Game) {
	input := g.GetInput()
	npc := p.TalkingTo
	if npc == nil {
		return
	}
	switch input.Typ {
	case Up:
		npc.TalkChoiceUp()
		adaptDialogSpeed()
	case Down:
		npc.TalkChoiceDown()
		adaptDialogSpeed()
	case Action:
		EM.Dispatch(&Event{Action: ActionTalk})
		npc.TalkConfirmChoice(g)
		adaptDialogSpeed()
		p.Intelligence.RaiseXp(1)
	default:
	}
}

func (p *Player) Take(g *Game, posTo Pos) bool {
	level := g.Level
	o := level.Map[posTo.Z][posTo.Y][posTo.X].Object
	if o != nil {
		p.IsTaking = true
		ut := p.TakeUsable(o, g)
		bt := p.TakeBook(o, g)
		qot := p.TakeSpecialObject(o, g)
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

func (p *Player) TakeSpecialObject(o *Object, g *Game) bool {
	qo, isQuestObject := g.SpecialObjects[o.Rune]
	if !isQuestObject {
		return false
	}
	EM.Dispatch(&Event{
		Action:  ActionTake,
		Message: "Vous avez un nouvel objet spécial!",
	})
	p.Inventory.QuestObjects[o.Rune] = o
	g.Level.Map[o.Z][o.Y][o.X].Object = nil

	for _, stID := range qo.StepsBeginning {
		g.beginStep(stID)
	}
	for _, stID := range qo.StepsFinishing {
		g.finishStep(stID)
	}

	return true
}

func (p *Player) TakeUsable(o *Object, g *Game) bool {
	taken := p.Inventory.TakeUsable(o)
	if taken {
		EM.Dispatch(&Event{Action: ActionTake})
		g.Level.Map[o.Z][o.Y][o.X].Object = nil
	}

	return taken
}

func (p *Player) TakeBook(o *Object, g *Game) bool {
	taken := p.AddBook(o, g)
	if taken {
		EM.Dispatch(&Event{
			Action:  ActionTake,
			Message: "Vous avez un nouveau livre!",
		})
		g.Level.Map[o.Z][o.Y][o.X].Object = nil
		p.Intelligence.RaiseXp(10)
	}

	return taken
}

func (p *Player) Recruit(npc *Npc, g *Game) {
	if p.Friend != nil {
		EM.Dispatch(&Event{
			Message: "Vous avez déjà un compagnon, impossible de recruter plus.",
		})
		return
	}
	f := g.Level.MakeFriend(npc)
	p.Friend = f
}

func (c *Player) PowerUse(g *Game, posTo Pos) {
	if c.Energy.Current > 0 {
		c.IsPowerUsing = true
		for c.PowerUsingStage = 0; c.PowerUsingStage < 32; c.PowerUsingStage++ {
			c.CurrentPower.adaptSpeed()
		}
		switch c.CurrentPower.UID {
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
			c.Health.Restore(c.CalculatePowerAttackScore())
			c.LooseEnergy(c.CurrentPower.Energy)
		default:
		}
		c.IsPowerUsing = false
	}
}

func (p *Player) TalkToDead(g *Game, posTo Pos) {
	level := g.Level
	npc := level.Map[posTo.Z][posTo.Y][posTo.X].Npc
	if npc != nil && npc.Talkable && npc.IsDead() {
		EM.Dispatch(&Event{Action: ActionTalk})
		p.TalkingTo = npc
		npc.Talk(p, g)
		adaptDialogSpeed()
	}
}

func (p *Player) LearnAttack(attackUID string) {
	for _, att := range p.Attacks {
		if att.Name == attackUID {
			return
		}
	}
	for _, att := range Attacks() {
		if att.UID == attackUID {
			p.Attacks = append(p.Attacks, att)
			EM.Dispatch(&Event{
				Message: "Vous avez appris une nouvelle attaque: " + att.Name + " !",
			})
		}
	}
}

func (p *Player) AddKey(key string) {
	new := p.Inventory.AddKey(key)
	if new {
		EM.Dispatch(&Event{
			Message: "Vous avez une nouvelle clé: " + key + " !",
		})
	}
}

func (p *Player) LooseGold(value int) {
	p.Inventory.Gold -= value
}
