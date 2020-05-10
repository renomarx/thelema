package game

import (
	"log"
	"time"
)

const DialogDeltaTime = 200

type Dialog struct {
	CurrentNode string                `yaml:"current_node"`
	Nodes       map[string]*StoryNode `yaml:"nodes"`
	initialNode string
}

type StoryNode struct {
	Messages   []string      `yaml:"messages"`
	AllChoices []StoryChoice `yaml:"choices"`
	Choices    []StoryChoice `yaml:"-"`
}

type StoryChoice struct {
	Cmd            string `yaml:"cmd"`
	NodeId         string `yaml:"node"`
	Highlighted    bool
	StepsFinishing []string       `yaml:"steps_finishing"`
	StepsCanceling []string       `yaml:"steps_canceling"`
	StepsBeginning []string       `yaml:"steps_beginning"`
	Required       map[string]int `yaml:"required"`
	Actions        []string       `yaml:"actions"`
}

func adaptDialogSpeed() {
	time.Sleep(time.Duration(DialogDeltaTime) * time.Millisecond)
}

func (d *Dialog) Init(g *Game) {
	d.initialNode = d.CurrentNode
}

func (d *Dialog) Close() {
	d.CurrentNode = d.initialNode
}

func (d *Dialog) GetCurrentNode() *StoryNode {
	node, exists := d.Nodes[d.CurrentNode]
	if !exists {
		panic("Dialog node " + d.CurrentNode + " does not exist")
	}
	return node
}

func (d *Dialog) GetNode(key string) *StoryNode {
	node, exists := d.Nodes[key]
	if !exists {
		panic("Dialog node " + key + " does not exist")
	}
	return node
}

func (d *Dialog) SetCurrentNode(key string) *StoryNode {
	node, exists := d.Nodes[key]
	if !exists {
		log.Printf("Node %s does not exist", key)
		return nil
	}
	d.CurrentNode = key
	d.initialNode = d.CurrentNode
	return node
}

func (n *StoryNode) GetHighlightedIndex() int {
	for i := 0; i < len(n.Choices); i++ {
		if n.Choices[i].Highlighted {
			return i
		}
	}
	return 0
}

func (n *StoryNode) ClearHighlight() {
	for j := 0; j < len(n.Choices); j++ {
		n.Choices[j].Highlighted = false
	}
}

func (n *StoryNode) SetHighlightedIndex(i int) {
	n.ClearHighlight()
	length := len(n.Choices)
	if i >= length {
		i = length - 1
	}
	if i < 0 {
		i = 0
	}
	idx := i
	if len(n.Choices) > idx {
		n.Choices[idx].Highlighted = true
	}
}

func (n *StoryNode) GetCurrentChoice() StoryChoice {
	idx := n.GetHighlightedIndex()
	if idx > len(n.Choices) {
		log.Printf("Error: No choice possible for dialog")
		return StoryChoice{}
	}
	return n.Choices[idx]
}

func (n *StoryNode) filterPossibleChoices(g *Game) {
	p := g.Level.Player
	var res []StoryChoice
	for _, choice := range n.AllChoices {
		isPossible := true
		for _, stID := range choice.StepsFinishing {
			st, e := g.Steps[stID]
			if !e {
				log.Fatalf("Step %s does not exist", stID)
			}
			if st.state != StepStateTODO {
				isPossible = false
			}
		}
		if len(choice.Required) > 0 {
			for ch, val := range choice.Required {
				switch ch {
				case "intelligence":
					if p.Intelligence.Current < val {
						isPossible = false
					}
				case "charisma":
					if p.Charisma.Current < val {
						isPossible = false
					}
				case "will":
					if p.Will.Current < val {
						isPossible = false
					}
				case "beauty":
					if p.Beauty.Current < val {
						isPossible = false
					}
				case "gold":
					if p.Inventory.Gold < val {
						isPossible = false
					}
				}
			}
		}
		if isPossible {
			res = append(res, choice)
		}
	}
	n.Choices = res
}

func (g *Game) UpdateNpcDialog(fromName, npcName, node string) {
	from, exists := g.Levels[fromName]
	if !exists {
		panic("Level " + fromName + " does not exist")
	}
	npc := from.SearchNpc(npcName)
	if npc == nil {
		panic("Npc " + npcName + " on level " + fromName + " does not exist")
	}
	log.Printf("Updating dialog node of %s to %s", npcName, node)
	npc.Dialog.SetCurrentNode(node)
}
