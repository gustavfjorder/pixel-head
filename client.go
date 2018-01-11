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
)

func run() {
	//config.LoadJson("settings.json", &config.Conf)
	registerModels()
	animations := client.Load("client/sprites", "", client.ANIM)
	animations["bullet"], _ = client.LoadAnimation(config.Conf.BulletPath)
	var (
		frames           = 0
		second           = time.Tick(time.Second)
		fps              = time.Tick(config.Conf.Fps)
		me               model.Player
		state            model.State
		activeAnimations = make(map[string]*client.Animation)
		spc, gameMap     = gotoLounge()
		imd              = client.LoadMap(gameMap)
		cfg              = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1600, 800),}
	)

	//Make window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	//Start handlers
	go client.HandleEvents(&spc, &state, &me)
	go client.HandleControls(&spc, win)

	for !win.Closed() {
		//Update visuals
		win.Clear(colornames.Darkolivegreen)

		imd.Draw(win)
		client.HandleAnimations(win, state, animations, activeAnimations)
		client.DrawAbilities(win, &me)
		client.DrawHealthbar(win, &me)

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
	var myUri string
	if config.Conf.Online {
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
	} else {
		// todo: Implement when Game/Server is final
		//go server.StartGame(myuri, []string{config.Conf.Id})
		//servspc = space.NewRemoteSpace(myuri)
	}
	spc = space.NewRemoteSpace(myUri)
	// Load map from server
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
