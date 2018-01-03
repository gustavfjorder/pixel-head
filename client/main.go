package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"
)

var (
	win *pixelgl.Window
	cfg pixelgl.WindowConfig
)

func run() {
	LoadConfiguration("conf.json", &config)
	makeWindow()
	var (
		frames     = 0
		second     = time.Tick(time.Second)
		fps        = time.Tick(time.Second / config.Fps)
		playerAnim = loadAnimations("sprites/survivor", "")
		zombieAnim = loadAnimations("sprites/zombie", "")
		zMove      = zombieAnim["walk"].start(config.AnimationSpeed)
		zMove2     = zombieAnim["walk"].start(config.AnimationSpeed)
		idle       = playerAnim["rifle.idle"].start(config.AnimationSpeed)
		move       = playerAnim["rifle.move"].start(config.AnimationSpeed)
		shoot      = playerAnim["rifle.shoot"].start(config.AnimationSpeed)
		curRot     = pixel.IM
	)
	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		curRot = handleDir(curRot)
		move.Next().Draw(win, curRot.Moved(win.Bounds().Center()))
		idle.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(200, 100))))
		shoot.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-200, -100))))
		zMove.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-400, -200))))
		zMove2.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(400, -200))))
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		<-fps
	}
}

func makeWindow() {
	cfg = pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
}

func main() {
	pixelgl.Run(run)
}
