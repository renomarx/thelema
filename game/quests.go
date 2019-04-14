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
	ObjectsTaken []string `json:"objects_taken"`
	GoldGiven    int      `json:"gold_given"`
	ObjectsGiven []string `json:"objects_given"`
	Final        bool     `json:"final"`
}

type QuestObject struct {
	Quest struct {
		ID               string   `json:"id"`
		StepsFullfilling []string `json:"steps_fullfilling"`
	} `json:"quest"`
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

func (p *Player) LoadQuestsObjects(dirpath string) {
	filename := dirpath + "/quests/objects.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	objects := make(map[string]*QuestObject)

	err = json.Unmarshal(byteValue, &objects)
	if err != nil {
		log.Fatal(err)
	}

	objectsByRune := make(map[rune]*QuestObject)
	for key, obj := range objects {
		objectsByRune[rune(key[0])] = obj
	}

	p.QuestsObjects = objectsByRune
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
		g.GetEventManager().Dispatch(&Event{Action: ActionQuestFinished})
	}
	for _, s := range st.ObjectsTaken {
		g.GetEventManager().Dispatch(&Event{
			Message: "Quest object taken",
			Action:  ActionTake})
		r := rune(s[0])
		_, exists := p.Inventory.QuestObjects[r]
		if exists {
			delete(p.Inventory.QuestObjects, r)
		}
	}
	for _, s := range st.ObjectsGiven {
		g.GetEventManager().Dispatch(&Event{
			Message: "Quest object given!",
			Action:  ActionTake})
		r := rune(s[0])
		_, exists := p.QuestsObjects[r]
		if !exists {
			panic("Quest object " + s + " does not exist")
		}
		p.Inventory.QuestObjects[r] = &Object{Rune: r}
	}
	if st.GoldGiven > 0 {
		g.GetEventManager().Dispatch(&Event{Action: ActionTakeGold})
		p.Inventory.Gold += st.GoldGiven
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
