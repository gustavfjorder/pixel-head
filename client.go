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
	"github.com/gustavfjorder/pixel-head/Config"
)

func run() {
	Config.LoadJson("settings.json", &Config.Conf)
	var spc space.Space
	me := model.NewPlayer(Config.Conf.Id)
	gob.Register(model.Player{})
	if Config.Conf.Online {
		spc = space.NewSpace(Config.Conf.Uri)
		_, err := spc.Put("client", Config.Conf.Id)
		panic(err)
	} else {
		spc = space.NewSpace("game")
	}
	var (
		frames      = 0
		second      = time.Tick(time.Second)
		fps         = time.Tick(time.Second / Config.Conf.Fps)
		playerAnim  = client.LoadAnimations("client/sprites/survivor", "")
		cfg         = pixelgl.WindowConfig{Title: "Pixel Rocks!", Bounds: pixel.R(0, 0, 1024, 768),}
		r           = model.Request{}
		curAnimPath = client.Prefix(r.WeaponName(), r.MovementName())
		curAnim     = playerAnim[curAnimPath].Start(Config.Conf.AnimationSpeed)
	)
	if Config.Conf.Online {
		_, err := spc.Get("ready", Config.Conf.Id, &me)
		if err != nil {
			panic(err)
		}
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	for !win.Closed() {
		win.Clear(colornames.Darkolivegreen)
		client.HandleControls(*win, &r)
		prefix := client.Prefix(r.WeaponName(), r.MovementName())
		if curAnimPath != prefix {
			if anim, ok := playerAnim[prefix]; ok {
				curAnimPath = prefix
				curAnim.ChangeAnimation(anim)
			}
		}
		curAnim.Next().Draw(win, r.GetRotation().Moved(win.Bounds().Center()))
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

	Config.SaveConfig("settings.json")
}

func main() {
	pixelgl.Run(run)
}
