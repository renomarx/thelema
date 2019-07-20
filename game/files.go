package game

import (
	"io/ioutil"
	"log"
)

func LoadFilenames(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, f := range files {
		if !f.IsDir() {
			res = append(res, f.Name())
		}
	}
	return res
}
