package game

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Quest struct {
	CurrentStepId string           `json:"current_step"`
	Name          string           `json:"name"`
	Steps         map[string]*Step `json:"steps"`
	IsFinished    bool
}

type Step struct {
	Order        int    `json:"order"`
	Description  string `json:"description"`
	IsFinished   bool
	ObjectsTaken []string       `json:"objects_taken"`
	GoldGiven    int            `json:"gold_given"`
	ObjectsGiven []string       `json:"objects_given"`
	Raising      map[string]int `json:"raising"`
	Final        bool           `json:"final"`
}

type QuestObject struct {
	Level string `json:"level"`
	Quest struct {
		ID               string   `json:"id"`
		StepsFullfilling []string `json:"steps_fullfilling"`
	} `json:"quest"`
}

type QuestLink struct {
	StepsMandatory   []string `json:"steps_mandatory"`
	StepsFullfilling []string `json:"steps_fullfilling"`
}

func (p *Player) LoadQuests(dirpath string) {
	filename := dirpath + "/quests/quests.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	quests := make(map[string]*Quest)

	json.Unmarshal(byteValue, &quests)

	p.Quests = quests
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
			p.Strength.Raise(val)
		case "dexterity":
			p.Dexterity.Raise(val)
		case "beauty":
			p.Beauty.Raise(val)
		case "intelligence":
			p.Intelligence.Raise(val)
		case "will":
			p.Will.Raise(val)
		case "charisma":
			p.Charisma.Raise(val)
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
