package game

import "fmt"

type AttackInterface interface {
	GetSpeed() int
	GetName() string
	Play(ring *FightingRing)
	SetFrom(f FighterInterface)
	SetTo(fs []FighterInterface)
}

type Attack struct {
	From  FighterInterface
	To    []FighterInterface
	Speed int
}

func (att *Attack) SetFrom(f FighterInterface) {
	att.From = f
}

func (att *Attack) SetTo(fs []FighterInterface) {
	att.To = fs
}

func (att *Attack) GetSpeed() int {
	return att.Speed
}

type BiteAttack struct {
	Attack
	Damages int
}

func (att *BiteAttack) GetName() string {
	return "Bite (10)"
}

func (att *BiteAttack) Play(ring *FightingRing) {
	for _, f := range att.To {
		f.TakeDamages(att.Damages)
	}
}

type SwordAttack struct {
	Attack
	Damages int
	Speed   int
}

func (att *SwordAttack) GetName() string {
	return fmt.Sprintf("Sword attack : damages %d, energy cost 0", att.Damages)
}

func (att *SwordAttack) GetSpeed() int {
	return att.Speed
}

func (att *SwordAttack) Play(ring *FightingRing) {
	for _, f := range att.To {
		f.TakeDamages(att.Damages)
	}
}

type PowerAttack struct {
	Attack
	Damages    int
	Name       string
	EnergyCost int
}

func (att *PowerAttack) GetName() string {
	return fmt.Sprintf(att.Name+" : damages %d, energy cost %d", att.Damages, att.EnergyCost)
}

func (att *PowerAttack) Play(ring *FightingRing) {
	if !att.From.HasEnoughEnergy(att.EnergyCost) {
		return
	}
	for _, f := range att.To {
		f.TakeDamages(att.Damages)
	}
	att.From.LooseEnergy(att.EnergyCost)
}
