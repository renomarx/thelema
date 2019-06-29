package game

import "fmt"

type AttackInterface interface {
	GetSpeed() int
	GetName() string
	Play(ring *FightingRing)
	SetTo(f FighterInterface)
}

type BiteAttack struct {
	From    FighterInterface
	To      FighterInterface
	Damages int
}

func (att *BiteAttack) GetSpeed() int {
	return 10
}

func (att *BiteAttack) GetName() string {
	return "Bite (10)"
}

func (att *BiteAttack) SetTo(f FighterInterface) {
	att.To = f
}

func (att *BiteAttack) Play(ring *FightingRing) {
	att.To.TakeDamages(att.Damages)
}

type SwordAttack struct {
	From    FighterInterface
	To      FighterInterface
	Damages int
	Speed   int
}

func (att *SwordAttack) GetName() string {
	return fmt.Sprintf("Sword attack : damages %d, energy cost 0", att.Damages)
}

func (att *SwordAttack) GetSpeed() int {
	return att.Speed
}

func (att *SwordAttack) SetTo(f FighterInterface) {
	att.To = f
}

func (att *SwordAttack) Play(ring *FightingRing) {
	att.To.TakeDamages(att.Damages)
}
