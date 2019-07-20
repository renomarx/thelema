package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"log"
	"path/filepath"
	"strings"
	"thelema/game"
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
	WindowWidth        int
	WindowHeight       int
	window             *sdl.Window
	renderer           *sdl.Renderer
	textureAtlas       *sdl.Texture
	textureIndex       map[game.Tile][]sdl.Rect
	Cam                Camera
	Game               *game.Game
	playerTextures     map[string]*sdl.Texture
	pnjTextures        map[string]*sdl.Texture
	backgroundTextures map[string]*sdl.Texture
	uiTextures         map[string]*sdl.Texture
	Fonts              map[int]*ttf.Font
	Texts              map[int]*TextCache
	Keymap             map[string]sdl.Keycode
	LastKeyDown        sdl.Keycode
	Mp                 *MusicPlayer
	Event              *UIEvent
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
	playerFiles := game.LoadFilenames("ui2d/assets/player/")
	for _, playerFile := range playerFiles {
		player := strings.Split(playerFile, ".")
		playerName := player[0]
		ui.playerTextures[playerName] = ui.imgFileToTexture("ui2d/assets/player/" + playerFile)
	}

	ui.pnjTextures = make(map[string]*sdl.Texture)
	pnjFiles := game.LoadFilenames("ui2d/assets/pnjs/")
	for _, pnjFile := range pnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = ui.imgFileToTexture("ui2d/assets/pnjs/" + pnjFile)
	}
	commonPnjFiles := game.LoadFilenames("ui2d/assets/pnjs/common/")
	for _, pnjFile := range commonPnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = ui.imgFileToTexture("ui2d/assets/pnjs/common/" + pnjFile)
	}
	enemyPnjFiles := game.LoadFilenames("ui2d/assets/pnjs/enemy/")
	for _, pnjFile := range enemyPnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/" + pnjFile)
	}

	ui.backgroundTextures = make(map[string]*sdl.Texture)
	ui.backgroundTextures["outdoor"] = ui.imgFileToTexture("ui2d/assets/backgrounds/battle-background-sunny-hillsx4.png")

	ui.uiTextures = make(map[string]*sdl.Texture)
	ui.uiTextures["downbox"] = ui.imgFileToTexture("ui2d/assets/ui/HUD.png")

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
	ui.Mp.LoadSounds()

	return &ui
}

func (ui *UI) On(e *game.Event) {
	if e.Message != "" {
		log.Println(e.Message)
		ui.Event = NewUIEvent(e.Message)
	}
	ui.Mp.On(e)
}

func (ui *UI) Draw() {
	ui.DrawLevel()
	ui.DrawFightingRing()
	ui.DrawMenu()
	ui.DrawGameGeneratorScreen()
	ui.DrawEvents()
	ui.renderer.Present()
}

func (ui *UI) drawObject(pos game.Pos, tile game.Tile) {
	if len(ui.textureIndex[tile]) > 0 {
		ui.renderer.Copy(ui.textureAtlas,
			&ui.textureIndex[tile][(pos.X+pos.Y)%len(ui.textureIndex[tile])],
			&sdl.Rect{X: int32(pos.X*Res) + ui.Cam.X, Y: int32(pos.Y*Res) + ui.Cam.Y, W: Res, H: Res})
	}
}

func (ui *UI) Run() {
	ui.Mp.PlayMusic(32)
	for {
		ui.Draw()
	}
}
