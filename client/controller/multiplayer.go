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
	"sync"
	"sort"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/server"
)

type Multiplayer struct {
	framework.Controller

	viewItems []component.ComponentInterface
}

var updateView func()

func (c *Multiplayer) Init() {
	c.viewItems = make([]component.ComponentInterface, 0)
}

func (c *Multiplayer) Run() {
	win := c.Container.Get("window").(*pixelgl.Window)

	broadcasting := false
	endListen := make(chan bool, 1)
	endBroadcast := make(chan bool, 1)


	headLine := component.NewTextWithContent("Select multiplayer game")
	headLine.SetSize(40)
	headLine.Color = colornames.Chocolate
	headLine.Pos(pixel.V(0, 250))

	menuContainer := component.NewBox(14, 6)
	menuContainer.Pos(pixel.V(
		win.Bounds().W() / 7,
		0,
	)).Center()

	buttonSP := component.NewButton(8)
	buttonSP.Pos(pixel.V(
		menuContainer.Bounds().W() / 2,
		(menuContainer.Bounds().H() / 2) + 12,
	)).Center()
	buttonSP.Text("Create game")
	buttonSP.OnLeftMouseClick(func() {
		if broadcasting {
			endBroadcast <- true
			buttonSP.Text("Create game")
			broadcasting = false
		} else {
			go server.NewLounge(1)
			go broadCastServer(endBroadcast)
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
			endBroadcast <- true
		}
		endListen <- true
		c.App.ChangeTo("main")
	})

	menuContainer.Child(buttonSP, buttonExit)

	gamesContainer := component.NewBox(14, 6)
	gamesContainer.Pos(pixel.V(
		- win.Bounds().W() / 7,
		0,
	)).Center()

	c.addViewItem(component.NewContainer(menuContainer, gamesContainer, headLine))


	type lanServer struct {
		addr    net.Addr
		ip      string
		selfDestructer *time.Timer
	}

	lock := &sync.RWMutex{}
	broadcasts := make(map[string]*lanServer)
	go listenForBroadCast(endListen, func(addr net.Addr, msg string) {
		fmt.Println("Found: ", addr.String(), msg)

		server, found := broadcasts[addr.String()]
		if found {
			server.selfDestructer.Reset(time.Second * 2)
		} else {
			broadcasts[addr.String()] = &lanServer{
				addr:    addr,
				ip:      msg,
				selfDestructer: time.AfterFunc(time.Second * 2, func() {
					lock.Lock()
					delete(broadcasts, addr.String())
					lock.Unlock()
				}),
			}
		}
	})

	updateView = func() {
		gamesContainer.ClearChildren()

		j := 0
		sorted := make([]string, len(broadcasts))
		for _, server := range broadcasts {
			sorted[j] = server.ip
			j++
		}
		sort.Strings(sorted)

		i := 0
		for _, ip := range sorted {
			uri := ip
			item := component.NewButton(10)
			item.Pos(pixel.V(
				gamesContainer.Bounds().W() / 2,
				gamesContainer.Bounds().H() - float64(20 * (i + 1)),
			)).Center()
			item.Text(ip)
			item.OnLeftMouseClick(func() {
				fmt.Println("Clicked: " + uri)
				config.Conf.LoungeUri = uri
				config.Conf.Online = true
				c.App.ChangeTo("game")
			})

			gamesContainer.Child(item)
			i++
		}
	}
}

func (c *Multiplayer) Update() {
	win := c.Container.Get("window").(*pixelgl.Window)

	win.Clear(colornames.Lightgoldenrodyellow)

	updateView()

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

	receivedUDP := make(chan bool)

	listen2 := func() {
		data := make([]byte, 4096)
		n, addr, err := listener.ReadFromUDP(data) // blocking

		if err == nil {
			msg := string(data[:n])
			handler(addr, msg)
			receivedUDP <- true
		}
	}

	go listen2()
	for {
		select {
		case <-c:
			goto endListen
		case <-receivedUDP:
			go listen2()
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

	ticker := time.Tick(time.Second / 10)
	for {
		ip, err := config.GetIp()
		if err != nil {
			continue
		}

		conn.Write([]byte(ip))

		select {
		case <-c:
			goto endBroadCast
		case <- ticker:

		}
	}

	endBroadCast:
}
