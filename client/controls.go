package client

import (
	"github.com/faiface/pixel/pixelgl"
	"math"
	. "github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/model"
	"time"
	"github.com/pspaces/gospace/space"
)

func HandleControls(spc *space.Space, win *pixelgl.Window, done <-chan bool, me *model.Player) {
	t := time.Tick(Conf.ServerHandleSpeed)
	start := time.Now()
	r := model.Request{PlayerId: ID}
	for {
		angle, i := 0.0, 0

		if win.Pressed(Conf.UpKey) {
			angle += math.Pi / 2
			i++
		}
		if win.Pressed(Conf.RightKey) {
			if i <= 0 {
				angle += math.Pi * 2
			}
			i++
		}
		if win.Pressed(Conf.LeftKey) {
			angle += math.Pi
			i++
		}
		if win.Pressed(Conf.DownKey) {
			angle += math.Pi * 3 / 2
			i++
		}
		if i <= 0 {
			r.Dir = math.NaN()
			r.Action = model.IDLE
		} else {
			r.Dir = angle / float64(i)
			r.Action = model.MOVE
		}
		switch {
		case win.Pressed(Conf.KnifeKey):
			r.Weapon = model.KNIFE
		case win.Pressed(Conf.RifleKey):
			r.Weapon = model.RIFLE
		case win.Pressed(Conf.ShotgunKey):
			r.Weapon = model.SHOTGUN
		case win.Pressed(Conf.HandgunKey):
			r.Weapon = model.HANDGUN
		}

		if win.Pressed(Conf.ReloadKey) && r.Weapon != model.KNIFE {
			r.Action = model.RELOAD
		} else if win.Pressed(Conf.ShootKey) {
			r.Action = model.SHOOT
		} else if win.Pressed(Conf.MeleeKey) {
			r.Action = model.MELEE
		} else if win.Pressed(Conf.BarrelKey) {
			r.Action = model.BARREL
		}
		r.PlayerId = ID
		if r.Valid(*me){
			r.Timestamp = time.Since(start)
			spc.Put(r)
		}

		select {
		case <-done:
			return
		default:
		}
		<-t
	}
}
