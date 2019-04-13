package ui2d

import (
	"thelema/game"
	"github.com/veandco/go-sdl2/sdl"
)

type Color sdl.Color

const DialogScreenOffsetX = 10
const DialogScreenOffsetY = 4

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

func (ui *UI) DrawDialog(p *game.Pnj) {
	if p.IsTalking {
		ui.drawDialogBox()
		node := p.Dialog.GetCurrentNode()

		text := node.Message

		offsetH := (ui.WindowHeight/Res - DialogScreenOffsetY) * Res
		tex := ui.GetTexture(text, TextSizeM, ColorDisabled)
		_, _, w, h, _ := tex.Query()
		ui.renderer.Copy(tex, nil, &sdl.Rect{DialogScreenOffsetX, int32(offsetH), w, h})

		offsetH += int(h)
		for _, choice := range node.Choices {
			tex := ui.GetTexture(choice.Cmd, TextSizeS, ColorActive)
			if choice.Highlighted {
				tex.SetColorMod(0, 255, 0)
			} else {
				tex.SetColorMod(255, 255, 255)
			}
			_, _, w, h, _ = tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{DialogScreenOffsetX + 10, int32(offsetH), w, h})
			offsetH += int(h)
		}
	}
}

func (ui *UI) drawDialogBox() {
	for x := 0; x <= ui.WindowWidth/Res; x++ {
		for y := ui.WindowHeight/Res - DialogScreenOffsetY; y <= ui.WindowHeight/Res; y++ {
			ui.renderer.Copy(ui.textureAtlas,
				&ui.textureIndex['Ʈ'][0],
				&sdl.Rect{X: int32(x * Res), Y: int32(y * Res), W: Res, H: Res})
		}
	}
}
