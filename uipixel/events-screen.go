package uipixel

import (
	"time"
)

type UIEvent struct {
	Message   string
	CreatedAt time.Time
	Lifetime  int // Seconds
}

func NewUIEvent(message string) *UIEvent {
	return &UIEvent{Message: message, CreatedAt: time.Now(), Lifetime: 3}
}

func (ui *UI) DrawEvents() {
	e := ui.Event
	if e != nil {
		t := time.Now()
		deltaD := t.Sub(e.CreatedAt)
		if deltaD >= time.Duration(e.Lifetime)*time.Second {
			ui.Event = nil
		}
		if e.Message != "" {
			// TODO : draw text
		}
	}
	// if e != nil && e.Message != "" {
	// 	tex := ui.GetTexture(e.Message, TextSizeM, ColorActive)
	// 	_, _, w, h, _ := tex.Query()
	// 	offsetX := ui.WindowWidth - int(w) - 10
	// 	offsetY := ui.WindowHeight - Res
	// 	ui.renderer.Copy(tex, nil, &sdl.Rect{int32(offsetX), int32(offsetY), w, h})
	// }
}
