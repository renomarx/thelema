package game

type BiteAttack struct {
	From    FighterInterface
	To      FighterInterface
	Damages int
}

func (att *BiteAttack) GetSpeed() int {
	return 10
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

func (att *SwordAttack) GetSpeed() int {
	return att.Speed
}

func (att *SwordAttack) Play(ring *FightingRing) {
	att.To.TakeDamages(att.Damages)
}
