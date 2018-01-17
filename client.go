package main

import (
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/framework"
	"github.com/gustavfjorder/pixel-head/client/controller"
	"github.com/gustavfjorder/pixel-head/setup"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	setup.RegisterModels()

	var (
		frames             = 0
		second             = time.Tick(time.Second)
		fps                = time.Tick(config.Conf.Fps)
	)

	//Make window
	cfg := pixelgl.WindowConfig{
		Title:     "Zombie Hunter 3000!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		//VSync:     true,
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
	app.AddController("multiplayer", &controller.Multiplayer{})
	app.AddController("game", &controller.Game{})
	app.AddController("game_over", &controller.GameOver{})

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
