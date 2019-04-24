package game

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const DialogDeltaTime = 200

type Dialog struct {
	CurrentNode string                `json:"current_node"`
	Nodes       map[string]*StoryNode `json:"nodes"`
}

type StoryNode struct {
	Initial    bool           `json:"initial"`
	Message    string         `json:"message"`
	AllChoices []*StoryChoice `json:"choices"`
	Choices    []*StoryChoice
}

type StoryChoice struct {
	Cmd         string `json:"cmd"`
	NodeId      string `json:"node"`
	Highlighted bool
	Quest       struct {
		ID               string   `json:"id"`
		StepsMandatory   []string `json:"steps_mandatory"`
		StepsFullfilling []string `json:"steps_fullfilling"`
	} `json:"quest"`
}

func (p *Pnj) LoadDialogs(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	dialog := &Dialog{}

	json.Unmarshal(byteValue, &dialog)

	p.Dialog = dialog
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
		if choice.Quest.ID != "" {
			if !p.IsQuestOpen(choice.Quest.ID) {
				isPossible = false
			}
			for _, stepID := range choice.Quest.StepsMandatory {
				if !p.IsQuestOpenStepFinished(choice.Quest.ID, stepID) {
					isPossible = false
				}
			}
		}
		if isPossible {
			res = append(res, choice)
		}
	}
	n.Choices = res
}
