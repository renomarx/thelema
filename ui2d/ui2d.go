package ui2d

import (
	"fmt"
	"path/filepath"
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
	WindowWidth    int
	WindowHeight   int
	window         *sdl.Window
	renderer       *sdl.Renderer
	textureAtlas   *sdl.Texture
	textureIndex   map[game.Tile][]sdl.Rect
	Cam            Camera
	Game           *game.Game
	playerTextures map[string]*sdl.Texture
	pnjTextures    map[string]*sdl.Texture
	Fonts          map[int]*ttf.Font
	Texts          map[int]*TextCache
	Keymap         map[string]sdl.Keycode
	LastKeyDown    sdl.Keycode
	Mp             *MusicPlayer
}

func init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	// defer sdl.Quit()
}

func NewUI(g *game.Game) *UI {
	ui := UI{
		WindowWidth:  1024,
		WindowHeight: 600,
		Game:         g}

	ui.LoadKeymap()

	window, err := sdl.CreateWindow(WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(ui.WindowWidth), int32(ui.WindowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	// defer window.Destroy()
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(ui.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.textureAtlas = ui.imgFileToTexture("ui2d/assets/tiles.png")
	ui.textureIndex = ui.loadTextureIndex("ui2d/assets/atlas-index.txt")

	ui.Cam.X = 0
	ui.Cam.Y = 0

	ui.playerTextures = make(map[string]*sdl.Texture)
	ui.playerTextures["sayan"] = ui.imgFileToTexture("ui2d/assets/player/sayan.png")
	ui.playerTextures["monk"] = ui.imgFileToTexture("ui2d/assets/player/monk.png")
	ui.playerTextures["namek"] = ui.imgFileToTexture("ui2d/assets/player/namek.png")
	ui.playerTextures["cyborg"] = ui.imgFileToTexture("ui2d/assets/player/cyborg.png")
	ui.pnjTextures = make(map[string]*sdl.Texture)
	ui.pnjTextures["jason"] = ui.imgFileToTexture("ui2d/assets/pnjs/jason.png")
	ui.pnjTextures["sarah"] = ui.imgFileToTexture("ui2d/assets/pnjs/sarah.png")

	if err := ttf.Init(); err != nil {
		panic(err)
	}
	ui.Fonts = make(map[int]*ttf.Font)
	ui.Texts = make(map[int]*TextCache)
	for i := MinFontSize; i < MaxFontSize; i++ {
		fontPath, _ := filepath.Abs("ui2d/assets/fonts/OpenSans-Regular.ttf")
		font, err := ttf.OpenFont(fontPath, i)
		if err != nil {
			panic(err)
		}
		ui.Fonts[i] = font
		ui.Texts[i] = NewTextCache()
	}

	ui.Mp = NewMusicPlayer()
	ui.Mp.LoadMusics()

	g.GetEventManager().Subscribe(&ui, game.PlayerEventsType)

	return &ui
}

func (ui *UI) On(e *game.Event) {
	fmt.Println(e.Message)
	ui.Mp.On(e)
}

func (ui *UI) Draw() {
	ui.DrawLevel()
	ui.DrawMenu()
	ui.DrawGameGeneratorScreen()
	ui.renderer.Present()
}

func (ui *UI) drawObject(pos game.Pos, tile game.Tile) {
	if len(ui.textureIndex[tile]) > 0 {
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][0],
			&sdl.Rect{X: int32(pos.X*Res) + ui.Cam.X, Y: int32(pos.Y*Res) + ui.Cam.Y, W: Res, H: Res})
	}
}

func (ui *UI) Run() {
	ui.Mp.PlayMusic()
	for {
		ui.Draw()
	}
}
