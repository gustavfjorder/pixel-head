package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"
	. "github.com/gustavfjorder/pixel-head/client"
)

func run() {
	LoadJson("settings.json", &Conf)
	var (
		frames     = 0
		second     = time.Tick(time.Second)
		fps        = time.Tick(time.Second / Conf.Fps)
		playerAnim = LoadAnimations("client/sprites/survivor", "")
		cfg = pixelgl.WindowConfig{	Title:  "Pixel Rocks!",	Bounds: pixel.R(0, 0, 1024, 768),}
		weapon = HANDGUN
		movement = IDLE
		curRot = pixel.IM
		curAnimPath = Prefix(weapon, movement)
		curAnim = playerAnim[curAnimPath].Start(Conf.AnimationSpeed)
	)
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		HandleDir(*win, &curRot, &weapon, &movement)
		if curAnimPath != Prefix(weapon, movement){
			if anim, ok := playerAnim[Prefix(weapon, movement)]; ok{
				curAnimPath = Prefix(weapon, movement)
				curAnim.ChangeAnimation(anim)
			}
		}
		curAnim.Next().Draw(win, curRot.Moved(win.Bounds().Center()))
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

	SaveConfig()
}

func main() {
	pixelgl.Run(run)
}
