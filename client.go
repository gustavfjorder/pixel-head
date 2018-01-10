package main

import (
	"github.com/gustavfjorder/pixel-head/model"
	"encoding/gob"
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"fmt"
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/client"
	"github.com/gustavfjorder/pixel-head/config"
)

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

func run() {
	//config.LoadJson("settings.json", &config.Conf)
	registerModels()

	var (
		//me               = model.NewPlayer(config.Conf.Id)
		frames           = 0
		second           = time.Tick(time.Second)
		fps              = time.Tick(time.Second / config.Conf.Fps)
		cfg              = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1600, 800),}
		r                = model.Request{PlayerId: config.Conf.Id}
		GameUri          string
		ClientUri        string
		state            = &model.State{}
		animations       = client.Load("client/sprites", "", client.ANIM)
		activeAnimations = make(map[string]*client.Animation)
		myspc            space.Space
		servspc          space.Space
		me                model.Player
	)
	bullet, err := client.LoadAnimation(config.Conf.BulletPath)
	if err != nil{
		panic(err)
	}
	bullet.Start(config.Conf.AnimationSpeed)
	animations["bullet"] = bullet

	for k := range animations {
		fmt.Print(k, " ")
	}

	if config.Conf.Online {
		servspc = space.NewRemoteSpace(config.Conf.LoungeUri)
		_, err := servspc.Put("request", config.Conf.Id)
		if err != nil {
			panic(err)
		}

		k, err := servspc.Get("join", config.Conf.Id, &GameUri, &ClientUri)
		fmt.Println(k)
		if err != nil {
			panic(err)
		}

		servspc = space.NewRemoteSpace(GameUri)
		myspc = space.NewRemoteSpace(ClientUri)
	} else {
		// todo: Implement when Game/Server is final
		//go server.StartGame(myuri, []string{config.Conf.Id})
		//servspc = space.NewRemoteSpace(myuri)
	}

	// Load map from server
	mapTuple, err := myspc.Get("map", &model.Map{})
	if err != nil {
		panic(err)
	}
	imd := client.LoadMap(mapTuple.GetFieldAt(1).(model.Map))

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	myspc.Put("joined")

	go client.HandleEvents(myspc, state, &me)

	win.SetSmooth(true)
	for !win.Closed() {
		//Handle controls -> send request
		client.HandleControls(*win, &r)
		myspc.Put(r)


		//Update visuals
		win.Clear(colornames.Darkolivegreen)
		imd.Draw(win)
		client.HandleAnimations(win, *state, animations, activeAnimations)
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

func main() {
	pixelgl.Run(run)
}
