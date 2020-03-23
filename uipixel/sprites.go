package uipixel

import (
	"bufio"
	"image"
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

func (ui *UI) drawObject(pos game.Pos, tile game.Tile) {
	if len(ui.textureIndex[tile]) > 0 {
		sprite := pixel.NewSprite(ui.textureAtlas, ui.textureIndex[tile][(pos.X+pos.Y)%len(ui.textureIndex[tile])])
		ui.DrawSprite(sprite, float64(pos.X)*Res+ui.Cam.X, float64(pos.Y)*Res+ui.Cam.Y)
	}
}

func (ui *UI) NewSprite(pic pixel.Picture, rect pixel.Rect) *pixel.Sprite {
	bounds := pic.Bounds()
	rect.Min.Y = bounds.Max.Y - rect.Min.Y
	rect.Max.Y = bounds.Max.Y - rect.Max.Y
	return pixel.NewSprite(pic, rect)
}

func (ui *UI) DrawSprite(sprite *pixel.Sprite, X, Y float64) {
	Y = ui.WindowHeight - Y - Res
	mat := pixel.IM
	mat = mat.Moved(pixel.V(X, Y))
	sprite.Draw(ui.win, mat)
}

func (ui *UI) DrawSpriteScaled(sprite *pixel.Sprite, X, Y, w, h float64) {
	Y = ui.WindowHeight - Y - Res
	mat := pixel.IM
	mat = mat.ScaledXY(pixel.V(0, 0), pixel.V(w, h))
	mat = mat.Moved(pixel.V(X, Y))
	sprite.Draw(ui.win, mat)
}

func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}
