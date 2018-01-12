package component

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type component interface {
	Draw(target pixel.Target, pos pixel.Vec, center ...bool)
	//Target() *pixel.Batch
}

type ClickableInterface interface {
	AddListener(button pixelgl.Button, handler func())
	RunListeners(button pixelgl.Button)
	HandleEvents(win *pixelgl.Window)

	OnLeftMouseClick(handler func())
	OnRightMouseClick(handler func())
}

type Clickable struct {
	ClickableInterface

	Pressed  bool
	Handlers map[pixelgl.Button][]func()
}

func (c *Clickable) AddListener(button pixelgl.Button, handler func()) {
	if len(c.Handlers) == 0 {
		c.Handlers = make(map[pixelgl.Button][]func())
	}

	handlers, found := c.Handlers[button]
	if ! found {
		handlers = make([]func(), 0)
	}

	handlers = append(handlers, handler)

	c.Handlers[button] = handlers
}

func (c *Clickable) RunListeners(button pixelgl.Button) {
	if len(c.Handlers) == 0 {
		return
	}

	handlers, found := c.Handlers[button]
	if ! found {
		return
	}

	for _, handler := range handlers {
		handler()
	}
}

func (c *Clickable) OnLeftMouseClick(handler func()) {
	c.AddListener(pixelgl.MouseButtonLeft, handler)
}

func (c *Clickable) OnRightMouseClick(handler func()) {
	c.AddListener(pixelgl.MouseButtonRight, handler)
}

func (c *Clickable) HandleEvents(win *pixelgl.Window) {
	//mouse := win.MousePosition()
	//if ! c.bounds.Contains(mouse) {
	//	return
	//}

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.Pressed = true
		c.RunListeners(pixelgl.MouseButtonLeft)
	} else if win.JustPressed(pixelgl.MouseButtonRight) {
		c.RunListeners(pixelgl.MouseButtonRight)
	}

	if win.JustReleased(pixelgl.MouseButtonLeft) {
		c.Pressed = false
	}
}
