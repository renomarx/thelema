package game

import (
	"log"
	"strings"
	"time"
)

const DialogDeltaTime = 200

type Dialog struct {
	CurrentNode string                `json:"current_node"`
	Nodes       map[string]*StoryNode `json:"nodes"`
}

type StoryNode struct {
	Initial    bool           `json:"initial"`
	Messages   []string       `json:"messages"`
	AllChoices []*StoryChoice `json:"choices"`
	Choices    []*StoryChoice
}

type StoryChoice struct {
	Cmd         string `json:"cmd"`
	NodeId      string `json:"node"`
	Highlighted bool
	Quest       QuestLink      `json:"quest"`
	Required    map[string]int `json:"required"`
	Actions     []string       `json:"actions"`
	BooksGiven  []string       `json:"books_given"`
}

func adaptDialogSpeed() {
	time.Sleep(time.Duration(DialogDeltaTime) * time.Millisecond)
}

func (d *Dialog) Init(p *Player) {
	for _, n := range d.Nodes {
		n.filterPossibleChoices(p)
	}
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

func (d *Dialog) SetInitialNode(key string) *StoryNode {
	node, exists := d.Nodes[key]
	if !exists {
		log.Printf("Node %s does not exist", key)
		return nil
	}
	for _, n := range d.Nodes {
		n.Initial = false
	}
	node.Initial = true
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
	len := len(n.Choices)
	if i < 0 {
		i = 0
	}
	if i >= len {
		i = len - 1
	}
	idx := i
	n.Choices[idx].Highlighted = true
}

func (n *StoryNode) GetCurrentChoice() *StoryChoice {
	idx := n.GetHighlightedIndex()
	return n.Choices[idx]
}

func (n *StoryNode) filterPossibleChoices(p *Player) {
	var res []*StoryChoice
	for _, choice := range n.AllChoices {
		isPossible := true
		for _, questStep := range choice.Quest.StepsMandatory {
			arr := strings.Split(questStep, ":")
			if len(arr) > 1 {
				questID := arr[0]
				stepID := arr[1]
				if !p.IsQuestOpen(questID) {
					isPossible = false
				}
				if !p.IsQuestOpenStepFinished(questID, stepID) {
					isPossible = false
				}
			}
		}
		for _, questStep := range choice.Quest.StepsFullfilling {
			arr := strings.Split(questStep, ":")
			if len(arr) > 1 {
				questID := arr[0]
				stepID := arr[1]
				if p.IsStepFinished(questID, stepID) {
					isPossible = false
				}
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
