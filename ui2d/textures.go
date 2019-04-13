package ui2d

import (
	"bufio"
	"thelema/game"
	"image/png"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *UI) loadTextureIndex(filename string) map[game.Tile][]sdl.Rect {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	indexMap := make(map[game.Tile][]sdl.Rect)
	for scanner.Scan() {
		sline := scanner.Text()
		line := strings.Split(sline, " ")
		tileRune, _ := utf8.DecodeRune([]byte(line[0]))
		xy := line[1]
		splitXY := strings.Split(xy, ",")
		x, _ := strconv.ParseInt(splitXY[0], 10, 64)
		y, _ := strconv.ParseInt(splitXY[1], 10, 64)
		c, _ := strconv.ParseInt(splitXY[2], 10, 64)
		var rects []sdl.Rect
		for i := 0; i < int(c); i++ {
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		indexMap[game.Tile(tileRune)] = rects
	}

	return indexMap
}

func (ui *UI) pixelsToTexture(pixels []byte, w int, h int) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	return tex
}

func (ui *UI) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	pixels := make([]byte, w*h*4)

	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	tex := ui.pixelsToTexture(pixels, w, h)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}
