package game

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Keymap *Keymap `json:"keymap"`
}

type Keymap struct {
	Up       string `json:"up"`
	Down     string `json:"down"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	Action   string `json:"action"`
	Action2  string `json:"action2"`
	Speed    string `json:"speed"`
	Shadow   string `json:"shadow"`
	Meditate string `json:"meditate"`
	Select   string `json:"select"`
	Escape   string `json:"escape"`
}

func (g *Game) LoadConfig() {

	config := &Config{}
	config.Keymap = &Keymap{
		Up:       "up",
		Down:     "down",
		Left:     "left",
		Right:    "right",
		Action:   "a",
		Action2:  "z",
		Speed:    "e",
		Shadow:   "s",
		Meditate: "d",
		Select:   "select",
		Escape:   "escape",
	}

	dirpath := g.GameDir
	filename := dirpath + "/../config.json"
	jsonFile, err := os.Open(filename)
	if err == nil {
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, config)
	} else {
		log.Println("No json config file found. Default values used.")
	}

	g.Config = config

}
