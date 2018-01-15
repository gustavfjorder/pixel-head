package main

import (
	"github.com/gustavfjorder/pixel-head/model"
	"encoding/gob"
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/framework"
	"github.com/gustavfjorder/pixel-head/client/controller"
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

	//Make window
	cfg := pixelgl.WindowConfig{
		Title:  "Zombie Hunter 3000!",
		Bounds: pixel.R(0, 0, 600, 500),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	container := framework.NewContainer()
	container.SetService("window", win)

	app := framework.NewApplication(container)

	app.AddController("main", &controller.MainMenu{})
	gameCntrl := &controller.Game{}
	gameCntrl.Init()
	app.AddController("game", gameCntrl)
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
	gob.Register([]model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Shot{})
	gob.Register([]model.Shot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.State{})
	gob.Register(model.Updates{})
	gob.Register(model.Barrel{})
	z := model.NewZombie(pixel.ZV,model.ZOMBIE)
	gob.Register(z)
	gob.Register([]model.ZombieI{z})
	var t time.Duration
	gob.Register(t)
}
