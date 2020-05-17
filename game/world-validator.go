package game

import (
	"log"
	"strings"
)

func (g *Game) validateStoryGraph() {
	initialSteps := make(map[string]Step)
	remainingSteps := make(map[string]Step)
	for stID, st := range g.Steps {
		initialSteps[stID] = st
		remainingSteps[stID] = st
	}
	ng := Game{}
	ng.Steps = initialSteps

	numberOfValidatedStepsThisTurn := 1
	for len(remainingSteps) > 0 && numberOfValidatedStepsThisTurn > 0 {
		numberOfValidatedStepsThisTurn := 0
		for stID := range remainingSteps {
			if ng.isStepAccessible(stID) {
				st := ng.Steps[stID]
				st.State = StepStateDONE
				ng.Steps[stID] = st
				delete(remainingSteps, stID)
				numberOfValidatedStepsThisTurn++
			}
		}
	}

	if len(remainingSteps) > 0 {
		var keys []string
		for stID := range remainingSteps {
			keys = append(keys, stID)
		}
		log.Fatalf("Steps never reached: %s", strings.Join(keys, ","))
	}
}

func (book *OBook) validate(key string, g *Game) {
	for _, st := range book.StepsBeginning {
		_, e := g.Steps[st]
		if !e {
			log.Fatalf("Book %s: steps_beginning: step %s does not exist", key, st)
		}
	}
	for _, st := range book.StepsFinishing {
		_, e := g.Steps[st]
		if !e {
			log.Fatalf("Book %s: steps_finishing: step %s does not exist", key, st)
		}
	}
}

func (npc *Npc) validate(g *Game) {
	if npc.Dialog == nil {
		return
	}
	for key, node := range npc.Dialog.Nodes {
		for i, ch := range node.AllChoices {
			for _, st := range ch.StepsBeginning {
				_, e := g.Steps[st]
				if !e {
					log.Fatalf("NPC %s, node %s, choice %d: steps_beginning: step %s does not exist", npc.Name, key, i, st)
				}
			}
			for _, st := range ch.StepsFinishing {
				_, e := g.Steps[st]
				if !e {
					log.Fatalf("NPC %s, node %s, choice %d: steps_finishing: step %s does not exist", npc.Name, key, i, st)
				}
			}
		}
	}
}

func (qo *SpecialObject) validate(key string, g *Game) {
	for _, st := range qo.StepsBeginning {
		_, e := g.Steps[st]
		if !e {
			log.Fatalf("Quest object %s: steps_beginning: step %s does not exist", key, st)
		}
	}
	for _, st := range qo.StepsFinishing {
		_, e := g.Steps[st]
		if !e {
			log.Fatalf("Quest object %s: steps_finishing: step %s does not exist", key, st)
		}
	}
}
