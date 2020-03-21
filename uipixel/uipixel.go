package uipixel

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const WindowTitle = "L'abbaye de Thelema"

const MinFontSize = 9
const MaxFontSize = 25

const Res = 32

type Camera struct {
	X int32
	Y int32
}

type UI struct {
	WindowWidth  float64
	WindowHeight float64
	// window             *sdl.Window
	// renderer           *sdl.Renderer
	// textureAtlas       *sdl.Texture
	// textureIndex       map[game.Tile][]sdl.Rect
	// Cam                Camera
	// Game               *game.Game
	// playerTextures     map[string]*sdl.Texture
	// pnjTextures        map[string]*sdl.Texture
	// backgroundTextures map[string]*sdl.Texture
	// uiTextures         map[string]*sdl.Texture
	// mapTextures        map[string]*sdl.Texture
	// Fonts              map[int]*ttf.Font
	// Texts              map[int]*TextCache
	// Keymap             map[string]sdl.Keycode
	// LastKeyDown        sdl.Keycode
	// Mp                 *MusicPlayer
	// Event              *UIEvent
}

func (ui *UI) doRun() {
	cfg := pixelgl.WindowConfig{
		Title:  WindowTitle,
		Bounds: pixel.R(0, 0, ui.WindowWidth, ui.WindowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	// HERE: draw everything

	pic, err := loadPicture("uipixel/assets/arcanea.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// END

	for !win.Closed() {
		win.Update()
	}
}

func (ui *UI) Run() {

	// HERE : load everything

	pixelgl.Run(ui.doRun)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
