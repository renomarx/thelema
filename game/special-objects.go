package game

type SpecialObject struct {
	Level          string   `yaml:"level"`
	StepsFinishing []string `yaml:"steps_finishing"`
	StepsBeginning []string `yaml:"steps_beginning"`
}
