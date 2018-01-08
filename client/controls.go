package client

import (
	"github.com/faiface/pixel/pixelgl"
	"math"
	. "github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/model"
)

func HandleControls(win pixelgl.Window, r *model.Request) {
	angle, i := 0.0, 0
	r.Shoot = false; r.Reload = false; r.Melee=false

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
		r.Move = false
	}else{
		r.Move = true
		angle/= float64(i)
		r.Dir = angle
	}
	switch {
	case win.JustPressed(Conf.KnifeKey):
		r.CurrentWep = model.Knife
	case win.JustPressed(Conf.RifleKey):
		r.CurrentWep = model.Rifle
	case win.JustPressed(Conf.ShotgunKey):
		r.CurrentWep = model.Shotgun
	case win.JustPressed(Conf.HandgunKey):
		r.CurrentWep = model.Handgun
	}


	if win.JustPressed(Conf.ReloadKey) && r.CurrentWep != model.Knife{
		r.Reload = true
	} else if win.JustPressed(Conf.ShootKey) {
		r.Shoot = true
	} else if win.JustPressed(Conf.MeleeKey){
		r.Melee = true
	}

	return
}
