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
	"runtime"
	"log"
)

func run() {
	//config.LoadJson("settings.json", &config.Conf)
	registerModels()
	var (
		campos           pixel.Vec
		frames           = 0
		second           = time.Tick(time.Second)
		fps              = time.Tick(config.Conf.Fps)
		me               = model.NewPlayer(config.Conf.Id)
		spc, gameMap     = gotoLounge()
		walls            = client.LoadMap(gameMap)
		cfg              = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1600, 800),}
		updateChan       = make(chan model.Updates, config.Conf.ServerHandleSpeed)
		state = &model.State{}
		animationHandler = client.NewAnimationHandler(updateChan)
	)

	//Start state handler
	go client.HandleEvents(&spc, state, updateChan)

	//Make window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	animationHandler.SetWindow(win)

	//Start control handler
	go client.HandleControls(&spc, win)

	for !win.Closed() {
		client.GetPlayer(state.Players, &me)
		campos = pixel.V(0, 0).Sub(me.Pos).Add(win.Bounds().Center())
		//Update visuals
		win.Clear(colornames.Darkolivegreen)

		walls.Draw(win)
		win.SetMatrix(pixel.IM.Moved(campos))
		animationHandler.Draw(*state)

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
		spc = space.NewRemoteSpace(uri)
		c := make(chan bool, 1)
		go server.Start(&g, []server.ClientSpace{clientSpace}, c)
	}
	spc.Get("map", &m)
	spc.Put("joined")

	return
}

func main() {
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
			time.Sleep(5 * time.Second)
		}
	}()
	pixelgl.Run(run)
}

func registerModels() {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Zombie{})
	gob.Register([]model.Zombie{})
	gob.Register(model.Shot{})
	gob.Register([]model.Shot{})
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
