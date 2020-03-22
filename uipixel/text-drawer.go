package uipixel

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const TextSizeXXL float64 = 3
const TextSizeXL float64 = 2
const TextSizeL float64 = 1.5
const TextSizeM float64 = 1
const TextSizeS float64 = 0.8
const TextSizeXS float64 = 0.5

var ColorActive = colornames.Green
var ColorDisabled = colornames.Gray
var ColorStandard = colornames.White

func (ui *UI) DrawText(str string, size float64, color color.RGBA, x, y float64) (float64, float64) {
	y = ui.WindowHeight - y - 10*size
	fontTxt := text.New(pixel.V(x, y), ui.fontAtlas)
	fontTxt.Color = color
	fmt.Fprintln(fontTxt, str)
	fontTxt.Draw(ui.win, pixel.IM.Scaled(fontTxt.Orig, float64(size)))

	bounds := fontTxt.Bounds()
	//log.Printf("Bounds of %s: %+v", str, bounds)
	return (bounds.Max.X - bounds.Min.X) * size, (bounds.Max.Y - bounds.Min.Y) * size
}
