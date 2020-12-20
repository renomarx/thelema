package main

import (
  "fmt"
  "os"
	"log"

  "github.com/lafriks/go-tiled"
  "github.com/lafriks/go-tiled/render"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{
	img *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

const mapPath = "data/tmx/example.tmx" // Path to your Tiled Map.

func main() {
    // Parse .tmx file.
    gameMap, err := tiled.LoadFromFile(mapPath)
    if err != nil {
        fmt.Printf("error parsing map: %s", err.Error())
        os.Exit(2)
    }

    fmt.Println(gameMap)

    // You can also render the map to an in-memory image for direct
    // use with the default Renderer, or by making your own.
    renderer, err := render.NewRenderer(gameMap)
    if err != nil {
        fmt.Printf("map unsupported for rendering: %s", err.Error())
        os.Exit(2)
    }

    // Render just layer 0 to the Renderer.
    err = renderer.RenderLayer(0)
    if err != nil {
        fmt.Printf("layer unsupported for rendering: %s", err.Error())
        os.Exit(2)
    }

    // Get a reference to the Renderer's output, an image.NRGBA struct.
    img := renderer.Result

    // And so on. You can also export the image to a file by using the
    // Renderer's Save functions.
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Poc ebiten with tmx file")

	g := Game{}
	g.img = ebiten.NewImageFromImage(img)
	if g.img == nil {
		panic("Error: img si nil")
	}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

	// Clear the render result after copying the output if separation of
	// layers is desired.
	renderer.Clear()
}
