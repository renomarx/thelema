package game

import "fmt"

type AttackInterface interface {
	GetSpeed() int
	GetName() string
	Play(ring *FightingRing)
	SetFrom(f FighterInterface)
	SetTo(f FighterInterface)
}

type Attack struct {
	From FighterInterface
	To   FighterInterface
}

func (att *Attack) SetFrom(f FighterInterface) {
	att.From = f
}

func (att *Attack) SetTo(f FighterInterface) {
	att.To = f
}

type BiteAttack struct {
	Attack
	Damages int
}

func (att *BiteAttack) GetSpeed() int {
	return 10
}

func (att *BiteAttack) GetName() string {
	return "Bite (10)"
}

func (att *BiteAttack) Play(ring *FightingRing) {
	att.To.TakeDamages(att.Damages)
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

func (att *SwordAttack) SetTo(f FighterInterface) {
	att.To = f
}

func (att *SwordAttack) Play(ring *FightingRing) {
	att.To.TakeDamages(att.Damages)
}
