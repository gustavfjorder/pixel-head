package controller

import (
	"github.com/gustavfjorder/pixel-head/framework"
	"github.com/gustavfjorder/pixel-head/client/gui/component"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"os"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/helper"
	"github.com/gustavfjorder/pixel-head/config"
)

type MainMenu struct {
	framework.Controller

	viewItems []component.ComponentInterface
}

func (c *MainMenu) Init() {
	c.viewItems = make([]component.ComponentInterface, 0)
}

func (c *MainMenu) Run() {
	headLine := component.NewTextWithContent("Zombie Hunter 3000!")
	headLine.SetSize(40)
	headLine.Color = colornames.Chocolate
	headLine.Pos(pixel.V(0, 250))

	menuContainer := component.NewBox(14, 6)
	menuContainer.Center()

	buttonSP := component.NewButton(8)
	buttonSP.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) + 25,
	)).Center()
	buttonSP.Text("Single Player")
	buttonSP.OnLeftMouseClick(func() {
		config.Conf.Online = false
		c.App.ChangeTo("game")
	})

	buttonMP := component.NewButton(8)
	buttonMP.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		menuContainer.Bounds().H() / 2,
	)).Center()
	buttonMP.Text("Multi Player")
	buttonMP.OnLeftMouseClick(func() {
		c.App.ChangeTo("multiplayer")
	})

	buttonExit := component.NewButton(8)
	buttonExit.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) - 25,
	)).Center()
	buttonExit.Text("Exit")
	buttonExit.OnLeftMouseClick(func() {
		os.Exit(0)
	})

	menuContainer.Child(buttonSP, buttonMP, buttonExit)

	c.addViewItem(component.NewContainer(menuContainer, headLine))
}

func (c *MainMenu) Update() {
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

func (c *MainMenu) addViewItem(viewItem component.ComponentInterface) {
	c.viewItems = append(c.viewItems, viewItem)
}
