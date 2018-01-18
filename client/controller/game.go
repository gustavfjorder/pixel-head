package controller

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/framework"
	"golang.org/x/image/colornames"
	"github.com/gustavfjorder/pixel-head/client"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/pspaces/gospace/space"
	"fmt"
	"github.com/gustavfjorder/pixel-head/server"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/client/animation"
	"strconv"
)

type Game struct {
	framework.Controller
	ah client.AnimationHandler

	me         *model.Player
	state      *model.State
	GameDone   chan bool
	Ready      bool
	LoadChan   chan bool
	LoadedAnim bool
	usedMap    *imdraw.IMDraw
}

func (g *Game) Init() {
	g.me = &model.Player{Id: config.ID}
	g.state = &model.State{}
	g.LoadChan = make(chan bool)
	g.GameDone = make(chan bool, 2)
	if !g.LoadedAnim {
		go func() {
			g.ah = client.NewAnimationHandler()
			g.LoadedAnim = true
			g.LoadChan <- true
		}()
	}
}

func (g *Game) Run() {
	go func() {
		if !g.LoadedAnim{
			<-g.LoadChan
		}
		var (
			spc, gameMap = gotoLounge()
		)

		g.usedMap = animation.LoadMap(gameMap)
		updateChan := make(chan model.Updates, config.Conf.ServerHandleSpeed)

		win := g.Container.Get("window").(*pixelgl.Window)
		g.ah.SetWindow(win)
		g.ah.SetUpdateChan(updateChan)

		//Start state handler
		go client.HandleEvents(&spc, g.state, updateChan, g.GameDone)
		//Start control handler
		go client.HandleControls(&spc, win, g.GameDone, g.me)
		g.Ready = true
	}()
}

func (g *Game) Update() {
	win := g.Container.Get("window").(*pixelgl.Window)
	win.Clear(colornames.Darkolivegreen)
	if g.Ready {
		client.GetPlayer(g.state.Players, g.me)

		g.usedMap.Draw(win)

		win.SetMatrix(pixel.IM.Moved(win.Bounds().Center().Sub(g.me.Pos)))

		g.ah.Draw(*g.state)

		select {
		case <- g.GameDone:
			g.ah.Clear()
			win.SetMatrix(pixel.IM)
			g.App.ChangeTo("game_over")
		default:
		}
	}

	win.Update()
}

var port = 31415
func gotoLounge() (spc space.Space, m model.Map) {
	if config.Conf.Online {
		var myUri string
		servspc := space.NewRemoteSpace(config.Conf.LoungeUri)
		_, err := servspc.Put("request", config.ID)
		if err != nil {
			panic(err)
		}

		k, err := servspc.Get("join", config.ID, &myUri)
		fmt.Println(k)
		if err != nil {
			panic(err)
		}
		spc = space.NewRemoteSpace(myUri)
		// Load map from server
	} else {
		g := model.NewGame([]string{config.ID}, "Test1")
		m = model.MapTemplates["Test1"]
		uri := "tcp://localhost:" + strconv.Itoa(port) + "/game"
		port++
		clientSpace := server.ClientSpace{
			Id:    config.ID,
			Uri:   uri,
			Space: server.SetupSpace(uri),
		}
		go server.Start(&g, []server.ClientSpace{clientSpace}, make(chan bool, 1))
		spc = space.NewRemoteSpace(uri)
	}
	spc.Get("map", &m)
	spc.Put("joined")

	return
}
