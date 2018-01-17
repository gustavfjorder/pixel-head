package controller

import (
	"github.com/gustavfjorder/pixel-head/framework"
	"github.com/gustavfjorder/pixel-head/client/gui/component"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/helper"
	"net"
	"fmt"
	"time"
)

type Multiplayer struct {
	framework.Controller

	viewItems []component.ComponentInterface
}

func (c *Multiplayer) Init() {
	c.viewItems = make([]component.ComponentInterface, 0)
}

func (c *Multiplayer) Run() {
	broadcasting := false
	endListen := make(chan bool)
	endBrodcast := make(chan bool)


	headLine := component.NewTextWithContent("Select multiplayer game")
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
	buttonSP.Text("Create game")
	buttonSP.OnLeftMouseClick(func() {
		if broadcasting {
			endBrodcast <- true
			buttonSP.Text("Create game")
			broadcasting = false
		} else {
			go broadCastServer(endBrodcast)
			buttonSP.Text("Close game")
			broadcasting = true
		}
	})

	buttonExit := component.NewButton(8)
	buttonExit.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) - 12,
	)).Center()
	buttonExit.Text("Back")
	buttonExit.OnLeftMouseClick(func() {
		if broadcasting {
			endBrodcast <- true
		}
		endListen <- true
		c.App.ChangeTo("main")
	})

	menuContainer.Child(buttonSP, buttonExit)

	itemList := component.NewList()
	itemList.Pos(pixel.V(-250, -250))

	broadcasts := make(map[string]bool)
	go listenForBroadCast(endListen, func(addr net.Addr, msg string) {
		_, found := broadcasts[addr.String()]
		if found && msg != "open" {
			delete(broadcasts, addr.String())
		} else {
			broadcasts[addr.String()] = true
			broadcasts["12.12.12.12:3000"] = true
		}

		itemList.ClearChildren()

		i := 0
		for addr := range broadcasts {
			a := addr
			item := component.NewButton(12)
			item.Pos(pixel.V(0, -float64(20 * i)))
			item.Text(addr)
			item.OnLeftMouseClick(func() {
				fmt.Println("Clicked: " + a)
			})

			itemList.Child(item)
			i++
		}
	})

	c.addViewItem(component.NewContainer(menuContainer, headLine, itemList))
}

func (c *Multiplayer) Update() {
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

func (c *Multiplayer) addViewItem(viewItem component.ComponentInterface) {
	c.viewItems = append(c.viewItems, viewItem)
}

func listenForBroadCast(c chan bool, handler func(addr net.Addr, msg string)) {
	listener, _ := net.ListenMulticastUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(225, 0, 0, 1),
		Port: 9999,
	})

	defer listener.Close()

	for {
		data := make([]byte, 4096)
		n, addr, err := listener.ReadFromUDP(data) // blocking

		if err == nil {
			msg := string(data[:n])
			handler(addr, msg)
		}

		select {
		case <-c:
			goto endListen
		default:

		}
	}

	endListen:
}

func broadCastServer(c chan bool) {
	localAddr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:25001")

	conn, _ := net.DialUDP("udp4", localAddr, &net.UDPAddr{
		IP:   net.IPv4(225, 0, 0, 1),
		Port: 9999,
	})

	defer conn.Close()

	ticker := time.Tick(time.Second * 5)
	for {
		conn.Write([]byte("open"))

		select {
		case <-c:
			conn.Write([]byte("closed"))
			goto endBroadCast
		default:

		}

		<- ticker
	}

	endBroadCast:
}
