package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/sdl"
)

type Color sdl.Color

const DialogScreenOffsetX = 30
const DialogScreenOffsetY = 20

type TextCache struct {
	Textures map[string]*sdl.Texture
}

func NewTextCache() *TextCache {
	tc := &TextCache{}
	tc.Textures = make(map[string]*sdl.Texture)
	return tc
}

func (ui *UI) GetTexture(text string, size int, color sdl.Color) *sdl.Texture {
	tex, exists := ui.Texts[size].Textures[text]
	if exists {
		return tex
	}
	fontSurface, _ := ui.Fonts[size].RenderUTF8Solid(text, color)
	tex, _ = ui.renderer.CreateTextureFromSurface(fontSurface)
	ui.Texts[size].Textures[text] = tex

	return tex
}

func (ui *UI) DrawDialog(c *game.Npc) {
	if c.TalkingTo != nil {
		ui.drawDialogBox()
		node := c.Dialog.GetCurrentNode()

		offsetH := (3 * ui.WindowHeight / 4) + DialogScreenOffsetY
		for _, text := range node.Messages {
			tex := ui.GetTexture(text, TextSizeL, ColorDisabled)
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{DialogScreenOffsetX, int32(offsetH), w, h})
			offsetH += int(h)
		}
		for _, choice := range node.Choices {
			tex := ui.GetTexture(choice.Cmd, TextSizeM, ColorActive)
			if choice.Highlighted {
				tex.SetColorMod(0, 255, 0)
			} else {
				tex.SetColorMod(255, 255, 255)
			}
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{DialogScreenOffsetX + 10, int32(offsetH), w, h})
			offsetH += int(h)
		}
	}
}

func (ui *UI) drawDialogBox() {
	ui.renderer.Copy(ui.uiTextures["downbox"],
		&sdl.Rect{X: 0, Y: 0, W: 320, H: 64},
		&sdl.Rect{X: 0, Y: int32(3 * ui.WindowHeight / 4), W: int32(ui.WindowWidth), H: int32(ui.WindowHeight / 4)})
}
