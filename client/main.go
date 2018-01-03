package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"
)

func main() {
	pixelgl.Run(run)

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames     = 0
		second     = time.Tick(time.Second)
		fps        = time.Tick(time.Second / config.Fps)
		playerAnim = loadAnimations("sprites/survivor", "")
		zombieAnim = loadAnimations("sprites/zombie", "")
		zMove      = zombieAnim["walk"].start()
		idle       = playerAnim["rifle.idle"].start()
		move       = playerAnim["rifle.move"].start()
		shoot      = playerAnim["rifle.shoot"].start()
	)

	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		move.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		idle.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(200, 100))))
		shoot.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-200, -100))))
		zMove.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-400, -200))))
		zMove.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(400, -200))))
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
