package component

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/pkg/errors"
	"github.com/gustavfjorder/pixel-head/helper"
)

type ComponentInterface interface {
	Pos(pos pixel.Vec) ComponentInterface
	ParentPos(pos pixel.Vec) ComponentInterface
	Center() ComponentInterface
	Render() ComponentInterface
	Draw(win *pixelgl.Window)
	Bounds() pixel.Rect
}

type ClickableInterface interface {
	AddListener(button pixelgl.Button, handler func())
	RunListeners(button pixelgl.Button)
	DetermineEvent(win *pixelgl.Window)

	OnLeftMouseClick(handler func())
	OnRightMouseClick(handler func())
}

type Component struct {
	ComponentInterface

	columns float64
	rows    float64

	center    bool
	pos       pixel.Vec
	bounds    pixel.Rect
	parentPos pixel.Vec

	children []ComponentInterface

	Batch *pixel.Batch
	Text  *text.Text

	pic   pixel.Picture
	Rects []pixel.Rect
	data  SpriteData
}

func (component *Component) Pos(pos pixel.Vec) ComponentInterface {
	component.pos = pos
	component.CalculateBounds()

	return component
}

func (component *Component) ParentPos(pos pixel.Vec) ComponentInterface {
	component.parentPos = pos
	component.CalculateBounds()

	return component
}

func (component *Component) Center() ComponentInterface {
	component.center = true
	component.CalculateBounds()

	return component
}

func (component *Component) CalculateBounds() {
	pos := component.parentPos.Add(component.pos)

	topRight := pixel.V(component.columns * component.data.Width, component.rows * component.data.Height)
	component.bounds = pixel.Rect{
		Min: pos,
		Max: pos.Add(topRight),
	}

	if component.center {
		move := component.bounds.Center().Sub(component.bounds.Max)
		component.bounds = component.bounds.Moved(move)
	}
}

func (component *Component) Bounds() pixel.Rect {
	return component.bounds
}

func (component *Component) Child(child ...ComponentInterface) {
	if len(component.children) == 0 {
		component.children = make([]ComponentInterface, 0)
	}

	component.children = append(component.children, child...)
}

func (component *Component) ClearChildren() {
	component.children = make([]ComponentInterface, 0)
}

func (component *Component) Draw(win *pixelgl.Window) {
	// todo: would be better to do this here. Investigate polymorphic behaviour
	//if component.center {
	//	move := component.bounds.Center().Sub(component.bounds.Max)
	//	component.bounds = component.bounds.Moved(move)
	//}
	//
	//component.Render()

	if component.Text != nil {
		pos := component.bounds.Min
		if component.center {
			pos = pos.Sub(component.Text.Bounds().Center())
		}
		component.Text.Draw(win, pixel.IM.Moved(pos))
		component.Text.Clear()
	}
	if component.Batch != nil {
		component.Batch.Draw(win)
		component.Batch.Clear()
	}

	for _, child := range component.children {
		child.ParentPos(component.bounds.Min)
		child.Render().Draw(win)

		var clickableInterface ClickableInterface
		if helper.TypeImplements(child, &clickableInterface) {
			child.(ClickableInterface).DetermineEvent(win)
		}
	}
}

func (component *Component) Render() ComponentInterface {
	panic(errors.New("COMPONENT HAS NO RENDER..."))
}

func (component *Component) loadSprite(file string) {
	component.pic, component.Rects, component.data = LoadSprite(file)
	component.Batch = pixel.NewBatch(&pixel.TrianglesData{}, component.pic)
}



type Clickable struct {
	Component

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

func (c *Clickable) DetermineEvent(win *pixelgl.Window) {
	panic(errors.New("CLICKABLE IS NOT HANDLING EVENTS..."))
}

func (c *Clickable) DistributeEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.Pressed = true
	}

	if win.JustReleased(pixelgl.MouseButtonLeft) {
		c.Pressed = false
		c.RunListeners(pixelgl.MouseButtonLeft)
	} else if win.JustReleased(pixelgl.MouseButtonRight) {
		c.RunListeners(pixelgl.MouseButtonRight)
	}
}
