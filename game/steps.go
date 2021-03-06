package game

import (
	"log"
	"sort"
	"time"
)

const StepStateTODO = "TODO"
const StepStateDONE = "DONE"
const StepStateCANCELED = "CANCELED"

const StepStrategyOneOf = "ONE_OF"
const StepStrategyAll = "ALL"
const StepStrategyDefault = StepStrategyOneOf

type Step struct {
	Name         string         `yaml:"name"`
	Parents      []string       `yaml:"parents"`
	Strategy     string         `yaml:"strategy"`
	ObjectsTaken []string       `yaml:"objects_taken"`
	GoldGiven    int            `yaml:"gold_given"`
	ObjectsGiven []string       `yaml:"objects_given"`
	Raising      map[string]int `yaml:"raising"`
	State        string
	UpdatedAt    time.Time
}

type StepsByDate []Step

func (a StepsByDate) Len() int           { return len(a) }
func (a StepsByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StepsByDate) Less(i, j int) bool { return a[i].UpdatedAt.After(a[j].UpdatedAt) }

func (g *Game) GetOrderedSteps(state string) []Step {
	var res []Step
	for _, st := range g.Steps {
		if st.State == state {
			res = append(res, st)
		}
	}
	sort.Sort(StepsByDate(res))
	return res
}

func (g *Game) isStepAccessible(stepID string) bool {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	if st.State == StepStateDONE || st.State == StepStateCANCELED {
		return false
	}
	if len(st.Parents) == 0 {
		return true
	}
	strategy := st.Strategy
	if strategy != StepStrategyAll && strategy != StepStrategyOneOf {
		strategy = StepStrategyDefault
	}
	switch strategy {
	case StepStrategyOneOf:
		for _, parentID := range st.Parents {
			parent, parentExists := g.Steps[parentID]
			if !parentExists {
				log.Fatalf("Step %s does not exist", parentID)
			}
			if parent.State == StepStateDONE {
				return true
			}
		}
		return false
	case StepStrategyAll:
		for _, parentID := range st.Parents {
			parent, parentExists := g.Steps[parentID]
			if !parentExists {
				log.Fatalf("Step %s does not exist", parentID)
			}
			if parent.State != StepStateDONE {
				return false
			}
		}
		return true
	}
	// Should never pass here
	return true
}

func (g *Game) beginStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.State = StepStateTODO
	st.UpdatedAt = time.Now()
	g.Steps[stepID] = st
}

func (g *Game) cancelStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.State = StepStateCANCELED
	st.UpdatedAt = time.Now()
	g.Steps[stepID] = st
}

func (g *Game) finishStep(stepID string) {
	st, stepExists := g.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + " does not exist")
	}
	st.State = StepStateDONE
	st.UpdatedAt = time.Now()
	g.Steps[stepID] = st

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
