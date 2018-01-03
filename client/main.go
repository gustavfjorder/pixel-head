package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
)

func main() {
	pixelgl.Run(run)

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	playerAnim := loadAnimations("sprites/survivor", "")
	zombieAnim := loadAnimations("sprites/zombie", "")
	zMove := zombieAnim["walk"].start()
	idle := playerAnim["rifle.idle"].start()
	move := playerAnim["rifle.move"].start()
	shoot := playerAnim["rifle.shoot"].start()

	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		move.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		idle.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(200, 100))))
		shoot.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-200, -100))))
		zMove.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-400, -200))))
		win.Update()
	}

}
