package client

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"fmt"
)

type Control struct {
	Left  pixelgl.Button
	Right pixelgl.Button
	Up    pixelgl.Button
	Down  pixelgl.Button
	Shoot pixelgl.Button
	Melee pixelgl.Button
	Knife pixelgl.Button
	Rifle pixelgl.Button
	Shotgun pixelgl.Button
	Handgun pixelgl.Button
}

type Movement string
type Weapon string

const(
	IDLE = Movement("idle")
	MOVE = Movement("move")
	SHOOT = Movement("shoot")
	MELEE = Movement("meleeattack")
	RIFLE = Weapon("rifle")
	HANDGUN = Weapon("handgun")
	KNIFE = Weapon("knife")
	SHOTGUN = Weapon("shotgun")
)

func (m Movement) String() string{
	return string(m)
}

func (w Weapon) String() string{
	return string(w)
}

func HandleDir(win pixelgl.Window, mat *pixel.Matrix, w *Weapon, m *Movement)  {
	angle, i := 0.0, 0
	if win.Pressed(Conf.Control.Up) {
		angle += math.Pi / 2
		i++
	}
	if win.Pressed(Conf.Control.Right) {
		if i <= 0 {
			angle += math.Pi * 2
		}
		i++
	}
	if win.Pressed(Conf.Control.Left) {
		angle += math.Pi
		i++
	}
	if win.Pressed(Conf.Control.Down) {
		angle += math.Pi * 3 / 2
		i++
	}
	if i <= 0 {
		*m = IDLE
	}else{
		*mat = pixel.IM.Rotated(pixel.V(0, 0), angle/float64(i))
		*m = MOVE
	}
	if win.Pressed(Conf.Control.Shoot) {
		*m = SHOOT
	}else if win.Pressed(Conf.Control.Melee){
		*m = MELEE
	}
	switch {
	case win.JustPressed(Conf.Control.Knife):
		*w = KNIFE
	case win.JustPressed(Conf.Control.Rifle):
		*w = RIFLE
	case win.JustPressed(Conf.Control.Shotgun):
		*w = SHOTGUN
	case win.JustPressed(Conf.Control.Handgun):
		*w = HANDGUN
	}
}

func Prefix(ps ...fmt.Stringer) (res string){
	if len(ps) > 0 {
		res = ps[0].String()
	}

	for _, s := range ps[1:] {
		res += "." + s.String()
	}
	return
}
