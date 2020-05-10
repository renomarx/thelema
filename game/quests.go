package game

import (
	"sort"
)

type Quest struct {
	CurrentStepId string           `yaml:"current_step"`
	Name          string           `yaml:"name"`
	Steps         map[string]*Step `yaml:"steps"`
	IsFinished    bool
}

type Step struct {
	Order        int    `yaml:"order"`
	Description  string `yaml:"description"`
	IsFinished   bool
	ObjectsTaken []string       `yaml:"objects_taken"`
	GoldGiven    int            `yaml:"gold_given"`
	ObjectsGiven []string       `yaml:"objects_given"`
	Raising      map[string]int `yaml:"raising"`
	Final        bool           `yaml:"final"`
}

type QuestObject struct {
	Level string `yaml:"level"`
	Quest struct {
		ID               string   `yaml:"id"`
		StepsFullfilling []string `yaml:"steps_fullfilling"`
	} `yaml:"quest"`
}

type QuestLink struct {
	StepsMandatory   []string `yaml:"steps_mandatory"`
	StepsFullfilling []string `yaml:"steps_fullfilling"`
}

func (p *Player) finishQuestStep(questID string, stepID string, g *Game) {
	q, questExists := p.Quests[questID]
	if !questExists {
		panic("Quest " + questID + " does not exist")
	}
	st, stepExists := q.Steps[stepID]
	if !stepExists {
		panic("Step " + stepID + "in quest " + questID + " does not exist")
	}
	st.IsFinished = true
	if st.Final {
		q.IsFinished = true
		EM.Dispatch(&Event{
			Action:  ActionQuestFinished,
			Message: "Vous avez fini la quête " + q.Name + " !",
		})
	}
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
		_, exists := g.QuestsObjects[s]
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

func (q *Quest) IsRunning() bool {
	if q.IsFinished {
		return false
	}
	for _, st := range q.Steps {
		if st.IsFinished {
			return true
		}
	}

	return false
}

func (q *Quest) GetOrderedSteps() []*Step {
	steps := make([]*Step, 0, len(q.Steps))
	for _, step := range q.Steps {
		steps = append(steps, step)
	}

	sort.Slice(steps, func(i, j int) bool {
		if steps[i].Order < steps[j].Order {
			return true
		}
		if steps[i].Order == steps[j].Order {
			return steps[i].Description > steps[j].Description
		}
		return false
	})

	return steps
}
