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
	"github.com/gustavfjorder/pixel-head/server"
)

func setupSpace(uri string) space.Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register(model.Zombie{})
	gob.Register(model.Shoot{})

	return space.NewRemoteSpace(uri)
}

func run() {
	Config.LoadJson("settings.json", &Config.Conf)
	var spc space.Space
	me := model.NewPlayer(Config.Conf.Id)
	if Config.Conf.Online {
		spc = setupSpace(Config.Conf.LoungeUri)
		_, err := spc.Put("client", Config.Conf.Id)
		panic(err)
	}
	var (
		frames     = 0
		second     = time.Tick(time.Second)
		fps        = time.Tick(time.Second / Config.Conf.Fps)
		cfg        = pixelgl.WindowConfig{Title: "Zombie Hunter 3000!", Bounds: pixel.R(0, 0, 1920, 1080),}
		r          = model.Request{}
		GameUri    = ""
		state      = client.StateLock{}
		animations = client.LoadAnimations("client/sprites", "")
		activeAnimations = make(map[string]client.Animation)
	)
	if Config.Conf.Online {
		_, err := spc.Get("ready", Config.Conf.Id, &GameUri, &me)
		if err != nil {
			panic(err)
		}
		spc = space.NewSpace(GameUri)
	} else {
		spc = setupSpace(Config.Conf.LocalUri)
		go server.StartGame(Config.Conf.LocalUri, []string{Config.Conf.Id})
	}
	go client.HandleEvents(spc, &state)
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	for !win.Closed() {
		client.HandleControls(*win, &r)
		spc.Put(Config.Conf.Id, r)
		win.Clear(colornames.Darkolivegreen)
		client.HandleAnimations(win, state, animations, activeAnimations)
		win.Update()

		//s, b := r.MovementArgs()
		//prefix := client.Prefix(r.WeaponName(), s)
		//if curAnimPath != prefix {
		//	if anim, ok := playerAnim[prefix]; ok {
		//		curAnimPath = prefix
		//		curAnim.ChangeAnimation(anim, b)
		//	}
		//}
		//if r.Move {
		//	me.Pos = me.Pos.Add(pixel.V(me.MoveSpeed, 0).Rotated(r.Dir))
		//}
		//transformation := r.GetRotation().Scaled(center, 0.5).Moved(center.Add(me.Pos))
		//curMap.Draw(win)
		//curAnim.Next().Draw(win, transformation)

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
