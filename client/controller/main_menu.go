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

	menuContainer := component.NewBox(15, 22)
	menuContainer.Center()

	buttonSP := component.NewButton(8)
	buttonSP.Text = "Single Player"
	buttonSP.Pos(pixel.V(140, 200)).Center()
	buttonSP.OnLeftMouseClick(func() {
		c.App.ChangeTo("game")
	})

	buttonExit := component.NewButton(8)
	buttonExit.Text = "Exit"
	buttonExit.Pos(pixel.V(140, 175)).Center()
	buttonExit.OnLeftMouseClick(func() {
		os.Exit(0)
	})

	menuContainer.Child(buttonSP, buttonExit)

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
