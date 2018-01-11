package main

import (
	"github.com/gustavfjorder/pixel-head/model"
	"encoding/gob"
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/client"
	"github.com/gustavfjorder/pixel-head/config"
	"golang.org/x/image/colornames"
	"github.com/gustavfjorder/pixel-head/server"
	"sync"
)

func run() {
	//config.LoadJson("settings.json", &config.Conf)
	registerModels()
	animations := client.Load("client/sprites", "", client.ANIM)
	animations["bullet"], _ = client.LoadAnimation(config.Conf.BulletPath)
	var (
		frames             = 0
		second             = time.Tick(time.Second)
		fps                = time.Tick(config.Conf.Fps)
		me                 = model.Player{Id: config.Conf.Id}
		state              = &model.State{}
		activeAnimations   = make(map[string]*client.Animation)
		spc, gameMap, game = gotoLounge()
		imd                = client.LoadMap(gameMap)
		cfg                = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1600, 800),}
		lock               = &sync.Mutex{}
	)

	//Start state handler
	if config.Conf.Online {
		go client.HandleEvents(&spc, state, lock)
	} else {
		state = &game.State
	}

	//Make window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	//Start control handler
	go client.HandleControls(&spc, win)

	for !win.Closed() {

		client.GetPlayer(state.Players, &me)

		//Update visuals
		win.Clear(colornames.Darkolivegreen)

		imd.Draw(win)
		lock.Lock()
		fmt.Println(state)
		client.HandleAnimations(win, *state, animations, activeAnimations)
		lock.Unlock()
		client.DrawAbilities(win, me)
		client.DrawHealthbar(win, me)

		win.Update()

		//Count FPS
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

		//Don't exceed the fps limit
		<-fps
	}
	config.SaveConfig("settings.json")
}

func gotoLounge() (spc space.Space, m model.Map, game *model.Game) {
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
		game = &g
		m = model.MapTemplates["Test1"]
		uri := config.Conf.LoungeUri
		clientSpace := server.ClientSpace{
			Id:    config.Conf.Id,
			Uri:   uri,
			Space: server.SetupSpace(uri),
		}
		c := make(chan bool, 1)
		go server.Start(game, []server.ClientSpace{clientSpace}, c)
		spc = space.NewRemoteSpace(uri)
	}
	spc.Get("map", &m)
	spc.Put("joined")

	return
}

func main() {
	pixelgl.Run(run)
}

func registerModels() {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Zombie{})
	gob.Register([]model.Zombie{})
	gob.Register(model.Shoot{})
	gob.Register([]model.Shoot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.State{})
}
