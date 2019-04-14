package game

import "strconv"
import "time"

type Player struct {
	Character
	Talker
	TalkingTo            *Pnj
	Quests               map[string]*Quest
	QuestsObjects        map[rune]*QuestObject
	Inventory            *Inventory
	Library              *Library
	IsTaking             bool
	Menu                 *Menu
	QuestMenuOpen        bool
	CharacterMenuOpen    bool
	Powers               map[string]*PlayerPower
	CurrentPower         *PlayerPower
	LastRegenerationTime time.Time
}

func (p *Player) Update(g *Game) {
	if p.isDead {
		return
	}
	input := g.GetInput()
	p.regenerate()
	if !p.IsMoving {
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
}

func (p *Player) regenerate() {
	t := time.Now()
	deltaD := t.Sub(p.LastRegenerationTime)
	if deltaD > time.Duration(1000/p.RegenerationSpeed.Current)*time.Millisecond {
		if p.Energy.Current < p.Energy.Initial {
			p.Energy.Current += 5
		}
		if p.Hitpoints.Current < p.Hitpoints.Initial {
			p.Hitpoints.Current += 1
		}
		p.LastRegenerationTime = time.Now()
	}
}

func (p *Player) Move(g *Game) {
	input := g.GetInput()
	level := g.Level
	switch input.Typ {
	case Up:
		if canGo(level, Pos{p.X, p.Y - 1}) {
			openDoor(g, Pos{p.X, p.Y - 1})
			closeDoor(g, Pos{p.X, p.Y})
			p.DispatchWalkingEvent(g)
			p.WalkUp()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Down:
		if canGo(level, Pos{p.X, p.Y + 1}) {
			openDoor(g, Pos{p.X, p.Y + 1})
			closeDoor(g, Pos{p.X, p.Y})
			p.DispatchWalkingEvent(g)
			p.WalkDown()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Left:
		if canGo(level, Pos{p.X - 1, p.Y}) {
			openDoor(g, Pos{p.X - 1, p.Y})
			closeDoor(g, Pos{p.X, p.Y})
			p.DispatchWalkingEvent(g)
			p.WalkLeft()
			p.openPortal(g, Pos{p.X, p.Y})
		}
	case Right:
		if canGo(level, Pos{p.X + 1, p.Y}) {
			openDoor(g, Pos{p.X + 1, p.Y})
			closeDoor(g, Pos{p.X, p.Y})
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
		p.Attack(g, posTo)
	case Power:
		p.PowerAttack(g)
	default:
	}
}

func (p *Player) DispatchWalkingEvent(g *Game) {
	g.GetEventManager().Dispatch(&Event{
		Type:   PlayerEventsType,
		Action: ActionWalk})
}

func (p *Player) WalkDown() {
	p.IsMoving = true
	p.Y++
	p.Yb = CaseLen
	go p.moveDown()
}

func (p *Player) WalkUp() {
	p.IsMoving = true
	p.Y--
	p.Yb = -1 * CaseLen
	go p.moveUp()
}

func (p *Player) WalkLeft() {
	p.IsMoving = true
	p.X--
	p.Xb = -1 * CaseLen
	go p.moveLeft()
}

func (p *Player) WalkRight() {
	p.IsMoving = true
	p.X++
	p.Xb = CaseLen
	go p.moveRight()
}

func (p *Player) Attack(g *Game, posToAttack Pos) {
	if !p.IsTalking && !p.IsTaking {
		level := g.Level
		p.IsMoving = true
		p.IsAttacking = true
		go func(p *Player) {
			for i := CaseLen; i > 0; i = i - 2 {
				p.adaptSpeed()
			}
			p.IsMoving = false
			p.IsAttacking = false
		}(p)
		if isThereAMonster(level, posToAttack) {
			m := level.Monsters[posToAttack]
			m.TakeDamage(level, p.Strength.Current)
			p.Strength.RaiseXp(2)
		}
	}
}

func (p *Player) PowerAttack(g *Game) {
	if p.Energy.Current > 0 {
		p.IsMoving = true
		p.IsPowerAttacking = true
		go func(p *Player) {
			for i := CaseLen; i > 0; i = i - 2 {
				p.adaptSpeed()
			}
			p.IsMoving = false
			p.IsPowerAttacking = false
			p.Energy.RaiseXp(10)
		}(p)
		switch p.CurrentPower.Type {
		case PowerEnergyBall:
			g.Level.MakeEnergyball(p.Pos, p.LookAt, p.CurrentPower.Strength, p.CurrentPower.Speed)
			p.Energy.Current -= p.CurrentPower.Energy
		case PowerInvocation:
			if g.Level.MakeInvocation(p.Pos, p.LookAt, p.CurrentPower) {
				p.Energy.Current -= p.CurrentPower.Energy
			}
		default:
		}
	}
}

func (p *Player) TakeDamage(game *Game, damage int) {
	if p.Hitpoints.Current <= 0 {
		p.Die(game)
		return
	}
	p.Hitpoints.Current -= damage
	game.Level.MakeExplosion(p.Pos, damage, 50)
	p.Hitpoints.RaiseXp(damage)

	game.GetEventManager().Dispatch(&Event{
		Type:    PlayerEventsType,
		Action:  ActionHurt,
		Message: "Health left :" + strconv.Itoa(p.Hitpoints.Current)})
}

func (p *Player) Die(g *Game) {
	if p.isDead {
		return
	}
	p.isDead = true
	g.GetMenu().Choices[1].Disabled = true
	g.GetEventManager().Dispatch(&Event{
		Type:    PlayerEventsType,
		Action:  ActionDie,
		Message: "You're dead !"})
}

func (p *Player) openPortal(g *Game, pos Pos) {
	g.Level.OpenPortal(g, pos)
}

func (p *Player) Talk(g *Game, posTo Pos) {
	level := g.Level
	pnj, exists := level.Pnjs[posTo]
	if exists {
		p.IsTalking = true
		p.TalkingTo = pnj
		pnj.Talk(p)
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
		pnj.TalkConfirmChoice()
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
	o, exists := level.Objects[posTo]
	if exists {
		p.IsMoving = true
		p.IsTaking = true
		p.TakeUsable(o, g)
		p.TakeBook(o, g)
		p.TakeQuestObject(o, g)
		go func(p *Player) {
			for i := CaseLen; i > 0; i = i - 2 {
				p.adaptSpeed()
			}
			p.IsTaking = false
			p.IsMoving = false
		}(p)
	}
}

func (p *Player) TakeQuestObject(o *Object, g *Game) bool {
	qo, isQuestObject := p.QuestsObjects[o.Rune]
	if !isQuestObject {
		return false
	}
	Mux.Lock()
	p.Inventory.QuestObjects[o.Rune] = o
	delete(g.Level.Objects, o.Pos)
	Mux.Unlock()

	for _, stepID := range qo.Quest.StepsFullfilling {
		p.finishQuestStep(qo.Quest.ID, stepID)
	}

	return true
}

func (p *Player) TakeUsable(o *Object, g *Game) bool {
	taken := p.Inventory.TakeUsable(o)
	if taken {
		Mux.Lock()
		delete(g.Level.Objects, o.Pos)
		Mux.Unlock()
	}

	return true
}

func (p *Player) TakeBook(o *Object, g *Game) bool {
	taken := p.AddBook(o, g)
	if taken {
		Mux.Lock()
		delete(g.Level.Objects, o.Pos)
		Mux.Unlock()
	}

	return true
}
