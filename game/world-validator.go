package game

import "log"

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
