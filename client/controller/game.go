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
	"sync"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel"
)

type Game struct {
	framework.Controller
	ah client.AnimationHandler

	me    *model.Player
	state *model.State
	lock  *sync.Mutex

	usedMap          *imdraw.IMDraw
}

func (g *Game) Init() {
	g.me = &model.Player{Id: config.Conf.Id}
	g.state = &model.State{}
	g.lock = &sync.Mutex{}
}

func (g *Game) Run() {
	var (
		spc, gameMap = gotoLounge()
	)

	g.usedMap = client.LoadMap(gameMap)

	//Start state handler
	updateChan := make(chan model.Updates)
	go client.HandleEvents(&spc, g.state, updateChan)
	g.ah = client.NewAnimationHandler(updateChan)

	win := g.Container.Get("window").(*pixelgl.Window)
	g.ah.SetWindow(win)

	//Start control handler
	go client.HandleControls(&spc, win)
}

func (g *Game) Update() {
	win := g.Container.Get("window").(*pixelgl.Window)

	client.GetPlayer(g.state.Players, g.me)

	win.Clear(colornames.Darkolivegreen)

	g.usedMap.Draw(win)

	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center().Sub(g.me.Pos)))

	g.ah.Draw(*g.state)

	win.Update()
}



func gotoLounge() (spc space.Space, m model.Map) {
	if config.Conf.Online {
		var myUri string
		servspc := space.NewRemoteSpace(config.Conf.LoungeUri)
		_, err := servspc.Put("request", config.Conf.Id)
		if err != nil {
			panic(err)
		}

		k, err := servspc.Get("join", config.Conf.Id, &myUri)
		fmt.Println(k)
		if err != nil {
			panic(err)
		}
		spc = space.NewRemoteSpace(myUri)
		// Load map from server
	} else {
		g := model.NewGame([]string{config.Conf.Id}, "Test1")
		m = model.MapTemplates["Test1"]
		uri := config.Conf.LoungeUri
		clientSpace := server.ClientSpace{
			Id:    config.Conf.Id,
			Uri:   uri,
			Space: server.SetupSpace(uri),
		}
		c := make(chan bool, 1)
		go server.Start(&g, []server.ClientSpace{clientSpace}, c)
		spc = space.NewRemoteSpace(uri)
	}
	spc.Get("map", &m)
	spc.Put("joined")

	return
}

