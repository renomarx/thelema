package uipixel

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"thelema/game"

	"github.com/faiface/pixel"
)

func (ui *UI) loadTextureIndex(filename string) map[game.Tile][]pixel.Rect {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	indexMap := make(map[game.Tile][]pixel.Rect)
	for scanner.Scan() {
		sline := scanner.Text()
		if sline == "" || sline[0:2] == "//" {
			continue
		}
		line := strings.Split(sline, " ")
		if len(line) < 2 {
			log.Fatalf("Bad format for line %s", sline)
		}
		tile := line[0]
		xy := line[1]
		splitXY := strings.Split(xy, ",")
		x, _ := strconv.ParseInt(splitXY[0], 10, 64)
		y, _ := strconv.ParseInt(splitXY[1], 10, 64)
		c, _ := strconv.ParseInt(splitXY[2], 10, 64)
		var rects []pixel.Rect
		for i := 0; i < int(c); i++ {
			rects = append(rects, pixel.R(float64(x*32), float64(y*32), float64((x+1)*32), float64((y+1)*32)))
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		indexMap[game.Tile(tile)] = rects
	}

	return indexMap
}
