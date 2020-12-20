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
	layer1 *ebiten.Image
	layer2 *ebiten.Image
}

func (g *Game) loadTmx(tmxFile string) {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFromFile(tmxFile)
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
	g.layer1 = ebiten.NewImageFromImage(img)

	// Clear the render result after copying the output if separation of
	// layers is desired.
	renderer.Clear()

	// Render just layer 0 to the Renderer.
	err = renderer.RenderLayer(1)
	if err != nil {
			fmt.Printf("layer unsupported for rendering: %s", err.Error())
			os.Exit(2)
	}

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	img = renderer.Result
	g.layer2 = ebiten.NewImageFromImage(img)

	// Clear the render result after copying the output if separation of
	// layers is desired.
	renderer.Clear()
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-200, -200)
	//op.GeoM.Scale(0.5, 0.5)
	screen.DrawImage(g.layer1, op)
	screen.DrawImage(g.layer2, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {


  // And so on. You can also export the image to a file by using the
  // Renderer's Save functions.

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Poc ebiten with tmx file")

	g := Game{}
	g.loadTmx("data/tmx/example.tmx")

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
