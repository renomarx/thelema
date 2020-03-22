package uipixel

import (
	"image"
	"log"
	"os"
	"strings"
	"thelema/game"
	"time"

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
	textureAtlas       pixel.Picture
	textureIndex       map[game.Tile][]pixel.Rect
	Cam                Camera
	Game               *game.Game
	playerTextures     map[string]pixel.Picture
	pnjTextures        map[string]pixel.Picture
	backgroundTextures map[string]pixel.Picture
	uiTextures         map[string]pixel.Picture
	mapTextures        map[string]pixel.Picture
	// Fonts              map[int]*ttf.Font
	// Texts              map[int]*TextCache
	// Keymap             map[string]sdl.Keycode
	// LastKeyDown        sdl.Keycode
	// Mp                 *MusicPlayer
	Event *UIEvent
}

func NewUI(g *game.Game) *UI {
	ui := UI{
		WindowWidth:  800.0,
		WindowHeight: 600.0,
		Game:         g,
	}

	//ui.LoadKeymap()

	ui.textureAtlas = loadPicture("assets/tiles.png")
	ui.textureIndex = ui.loadTextureIndex("assets/atlas-index.txt")

	ui.Cam.X = 0
	ui.Cam.Y = 0

	ui.playerTextures = make(map[string]pixel.Picture)
	playerFiles := game.LoadFilenames("assets/player/")
	for _, playerFile := range playerFiles {
		player := strings.Split(playerFile, ".")
		playerName := player[0]
		ui.playerTextures[playerName] = loadPicture("assets/player/" + playerFile)
	}

	ui.pnjTextures = make(map[string]pixel.Picture)
	pnjFiles := game.LoadFilenames("assets/pnjs/")
	for _, pnjFile := range pnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = loadPicture("assets/pnjs/" + pnjFile)
	}
	commonPnjFiles := game.LoadFilenames("assets/pnjs/common/")
	for _, pnjFile := range commonPnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = loadPicture("assets/pnjs/common/" + pnjFile)
	}
	enemyPnjFiles := game.LoadFilenames("assets/pnjs/enemy/")
	for _, pnjFile := range enemyPnjFiles {
		pnj := strings.Split(pnjFile, ".")
		pnjName := pnj[0]
		ui.pnjTextures[pnjName] = loadPicture("assets/pnjs/enemy/" + pnjFile)
	}

	ui.backgroundTextures = make(map[string]pixel.Picture)
	ui.backgroundTextures["outdoor"] = loadPicture("assets/backgrounds/battle-background-sunny-hillsx4.png")

	ui.uiTextures = make(map[string]pixel.Picture)
	ui.uiTextures["downbox"] = loadPicture("assets/ui/HUD.png")

	ui.mapTextures = make(map[string]pixel.Picture)
	mapFiles := game.LoadFilenames("assets/maps")
	for _, mapFile := range mapFiles {
		ma := strings.Split(mapFile, ".")
		mapName := strings.ReplaceAll(ma[0], "-", "/")
		ui.mapTextures[mapName] = loadPicture("assets/maps/" + mapFile)
	}
	//
	// if err := ttf.Init(); err != nil {
	// 	panic(err)
	// }
	// ui.Fonts = make(map[int]*ttf.Font)
	// ui.Texts = make(map[int]*TextCache)
	// for i := MinFontSize; i < MaxFontSize; i++ {
	// 	fontPath, _ := filepath.Abs("assets/fonts/OpenSans-Regular.ttf")
	// 	font, err := ttf.OpenFont(fontPath, i)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ui.Fonts[i] = font
	// 	ui.Texts[i] = NewTextCache()
	// }
	//
	// ui.Mp = NewMusicPlayer()
	// ui.Mp.LoadMusics()
	// ui.Mp.LoadSounds()

	return &ui
}

func (ui *UI) Run() {

	// HERE : load everything

	pixelgl.Run(ui.doRun)
}

func (ui *UI) On(e *game.Event) {
	if e.Message != "" {
		log.Println(e.Message)
		ui.Event = NewUIEvent(e.Message)
	}
	//ui.Mp.On(e)
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

	pic := loadPicture("uipixel/assets/arcanea.png")
	sprite := pixel.NewSprite(pic, pic.Bounds())

	angle := 0.0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		angle += 3 * dt

		win.Clear(colornames.Skyblue)

		// HERE: draw everything
		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.Rotated(win.Bounds().Center(), angle)
		sprite.Draw(win, mat)
		// END

		win.Update()
	}
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