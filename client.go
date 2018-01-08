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
	gob.Register(model.Zombie{})
	gob.Register(model.Shoot{})
}

func run() {
	config.LoadJson("settings.json", &config.Conf)
	registerModels()

	var (
		me               = model.NewPlayer(config.Conf.Id)
		frames           = 0
		second           = time.Tick(time.Second)
		fps              = time.Tick(time.Second / config.Conf.Fps)
		cfg              = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1920, 1080),}
		r                = model.Request{}
		GameUri          string
		ClientUri        string
		state            = client.StateLock{}
		animations       = client.LoadAnimations("client/sprites", "")
		activeAnimations = make(map[string]client.Animation)
		myspc            space.Space
		servspc          space.Space
	)

	if config.Conf.Online {
		servspc = space.NewRemoteSpace(config.Conf.LoungeUri)
		_, err := servspc.Put("client", config.Conf.Id)
		if err != nil {
			panic(err)
		}

		_, err = servspc.Get("ready", config.Conf.Id, &GameUri, &ClientUri)
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

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load map from server
	mapTuple, err := myspc.Get("map", &model.Map{})
	if err != nil {
		panic(err)
	}
	imd := client.LoadMap(mapTuple.GetFieldAt(1).(model.Map))

	go client.HandleEvents(myspc, &state)

	win.SetSmooth(true)
	for !win.Closed() {
		//Handle controls -> send request
		client.HandleControls(*win, &r)
		servspc.Put(config.Conf.Id, r)

		//Update visuals
		win.Clear(colornames.Darkolivegreen)
		imd.Draw(win)
		client.HandleAnimations(win, state, animations, activeAnimations)
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
