package game

import (
	"sort"
	"time"
)

const StepStateTODO = "TODO"
const StepStateDONE = "DONE"
const StepStateCANCELED = "CANCELED"

type Step struct {
	Name         string         `yaml:"name"`
	ObjectsTaken []string       `yaml:"objects_taken"`
	GoldGiven    int            `yaml:"gold_given"`
	ObjectsGiven []string       `yaml:"objects_given"`
	Raising      map[string]int `yaml:"raising"`
	state        string
	updatedAt    time.Time
}

type StepsByDate []*Step

func (a StepsByDate) Len() int           { return len(a) }
func (a StepsByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StepsByDate) Less(i, j int) bool { return a[i].updatedAt.After(a[j].updatedAt) }

func (g *Game) GetOrderedSteps(state string) []*Step {
	var res []*Step
	for _, st := range g.Steps {
		if st.state == state {
			res = append(res, st)
		}
	}
	sort.Sort(StepsByDate(res))
	return res
}

func (g *Game) beginStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.state = StepStateTODO
	st.updatedAt = time.Now()
}

func (g *Game) cancelStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.state = StepStateCANCELED
	st.updatedAt = time.Now()
}

func (g *Game) finishStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.state = StepStateDONE
	st.updatedAt = time.Now()

	p := g.Level.Player
	EM.Dispatch(&Event{
		Action: ActionStepFinished,
	})
	for _, s := range st.ObjectsTaken {
		EM.Dispatch(&Event{
			Message: "Objet spécial pris.",
			Action:  ActionTake})
		_, exists := p.Inventory.QuestObjects[s]
		if exists {
			delete(p.Inventory.QuestObjects, s)
		}
	}
	for _, s := range st.ObjectsGiven {
		EM.Dispatch(&Event{
			Message: "Objet spécial récupéré!",
			Action:  ActionTake})
		_, exists := g.SpecialObjects[s]
		if !exists {
			panic("Quest object " + s + " does not exist")
		}
		p.Inventory.QuestObjects[s] = &Object{Rune: s}
	}
	if st.GoldGiven > 0 {
		EM.Dispatch(&Event{Action: ActionTakeGold})
		p.Inventory.Gold += st.GoldGiven
	}
	for ch, val := range st.Raising {
		switch ch {
		case "strength":
			p.Strength.RaisePermanently(val)
		case "dexterity":
			p.Dexterity.RaisePermanently(val)
		case "beauty":
			p.Beauty.RaisePermanently(val)
		case "intelligence":
			p.Intelligence.RaisePermanently(val)
		case "will":
			p.Will.RaisePermanently(val)
		case "charisma":
			p.Charisma.RaisePermanently(val)
		}
	}
}
