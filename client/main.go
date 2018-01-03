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

	playerAnim := loadAnimations("sprites/survivor", "player")
	idle := playerAnim["player.rifle.idle"]
	move := playerAnim["player.rifle.move"]
	shoot := playerAnim["player.rifle.shoot"]

	count := 0
	shootcount := 0

	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		move.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		idle.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(200, 100))))
		shoot.Next().Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.V(-200, -100))))
		count = (count + 1) % 20
		shootcount = (shootcount + 1) % 3
		win.Update()
	}

}
