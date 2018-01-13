package main

import (
	"github.com/gustavfjorder/pixel-head/model"
	"encoding/gob"
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/gustavfjorder/pixel-head/client"
	"github.com/gustavfjorder/pixel-head/config"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	registerModels()

	var (
		frames             = 0
		second             = time.Tick(time.Second)
		fps                = time.Tick(config.Conf.Fps)
	)

	//Start state handler
	go client.HandleEvents(&spc, state, updateChan)

	//Make window
	cfg := pixelgl.WindowConfig{
		Title:  "Zombie Hunter 3000!",
		Bounds: pixel.R(0, 0, 1024, 768),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	animationHandler.SetWindow(win)

	container := framework.NewContainer()
	container.SetService("window", win)

	app := framework.NewApplication(container)

	app.AddController("main", &controller.MainMenu{})
	app.AddController("game", &controller.Game{})

	app.SetController("main")
	app.Run()

	for ! win.Closed() {
		app.Update()

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
