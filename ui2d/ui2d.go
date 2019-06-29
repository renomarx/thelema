package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"log"
	"path/filepath"
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
	Event          *UIEvent
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
	playerNames := []string{
		"asmodeus",
		"kali",
		"baal",
		"lilith",
		"dagon",
		"lucifer",
		"hecate",
	}
	for _, n := range playerNames {
		ui.playerTextures[n] = ui.imgFileToTexture("ui2d/assets/player/" + n + ".png")
		ui.playerTextures[n+"_with_knife"] = ui.imgFileToTexture("ui2d/assets/player/" + n + "_with_knife.png")
		ui.playerTextures[n+"_with_bow"] = ui.imgFileToTexture("ui2d/assets/player/" + n + "_with_bow.png")
		ui.playerTextures[n+"_with_spear"] = ui.imgFileToTexture("ui2d/assets/player/" + n + "_with_spear.png")
		ui.playerTextures[n+"_with_wand"] = ui.imgFileToTexture("ui2d/assets/player/" + n + "_with_wand.png")
	}

	ui.pnjTextures = make(map[string]*sdl.Texture)
	// VIPs
	ui.pnjTextures["jason"] = ui.imgFileToTexture("ui2d/assets/pnjs/jason.png")
	ui.pnjTextures["sarah"] = ui.imgFileToTexture("ui2d/assets/pnjs/sarah.png")
	ui.pnjTextures["nathaniel"] = ui.imgFileToTexture("ui2d/assets/pnjs/nathaniel.png")
	// Common
	ui.pnjTextures["monk"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/monk.png")
	ui.pnjTextures["lord"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/lord.png")
	ui.pnjTextures["warrior"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/warrior.png")
	ui.pnjTextures["policeman"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/policeman.png")
	ui.pnjTextures["doctor"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/doctor.png")
	ui.pnjTextures["artist"] = ui.imgFileToTexture("ui2d/assets/pnjs/common/artist.png")
	// Enemies
	ui.pnjTextures["orc_thief"] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/orc_thief.png")
	ui.pnjTextures["orc_guard"] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/orc_guard.png")
	ui.pnjTextures["skeleton_warrior"] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/skeleton_warrior.png")
	ui.pnjTextures["skeleton_sorcerer"] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/skeleton_sorcerer.png")
	ui.pnjTextures["skeleton_lord"] = ui.imgFileToTexture("ui2d/assets/pnjs/enemy/skeleton_lord.png")

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

	g.GetEventManager().Subscribe(&ui)

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
