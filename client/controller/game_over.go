package controller

import (
	"github.com/gustavfjorder/pixel-head/framework"
	"github.com/gustavfjorder/pixel-head/client/gui/component"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/helper"
)

type GameOver struct {
	framework.Controller

	viewItems []component.ComponentInterface
}

func (c *GameOver) Init() {
	c.viewItems = make([]component.ComponentInterface, 0)
}

func (c *GameOver) Run() {
	headLine := component.NewTextWithContent("Game over...")
	headLine.SetSize(40)
	headLine.Color = colornames.Chocolate
	headLine.Pos(pixel.V(0, 250))

	menuContainer := component.NewBox(14, 6)
	menuContainer.Center()

	buttonSP := component.NewButton(8)
	buttonSP.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) + 12,
	)).Center()
	buttonSP.Text("Go to menu")
	buttonSP.OnLeftMouseClick(func() {
		c.App.ChangeTo("main")
	})

	buttonExit := component.NewButton(8)
	buttonExit.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) - 12,
	)).Center()
	buttonExit.Text("Exit")
	buttonExit.OnLeftMouseClick(func() {
		os.Exit(0)
	})

	menuContainer.Child(buttonSP, buttonExit)

	c.addViewItem(component.NewContainer(menuContainer, headLine))

	// Reset window position
	c.Container.Get("window").(*pixelgl.Window).SetMatrix(pixel.IM)
}

func (c *GameOver) Update() {
	win := c.Container.Get("window").(*pixelgl.Window)

	win.Clear(colornames.Lightgoldenrodyellow)

	for _, view := range c.viewItems {
		view.Pos(win.Bounds().Center())
		view.Render().Draw(win)

		var clickableInterface component.ClickableInterface
		if helper.TypeImplements(view, &clickableInterface) {
			view.(component.ClickableInterface).DetermineEvent(win)
		}
	}

	win.Update()
}

func (c *GameOver) addViewItem(viewItem component.ComponentInterface) {
	c.viewItems = append(c.viewItems, viewItem)
}
