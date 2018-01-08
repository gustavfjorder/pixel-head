package client

import (
	"errors"
	"github.com/faiface/pixel"
	_ "image/png"
	"time"
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
	"github.com/gustavfjorder/pixel-head/config"
)

type Animation struct {
	prefix   string
	Sprites  []*pixel.Sprite
	Cur      int
	Tick     *time.Ticker
	NextAnim *Animation
	Blocking bool
}

func (a *Animation) Start(s time.Duration) {
	a.Tick = time.NewTicker(time.Second / s)
}

func (a *Animation) Next() (s *pixel.Sprite) {
	s = a.Sprites[a.Cur]
	select {
	case <-a.Tick.C:
		a.Cur = (a.Cur + 1) % len(a.Sprites)
		if a.Cur <= 0 && a.NextAnim != nil && len(a.NextAnim.Sprites) > 0 {
			a.Blocking = a.NextAnim.Blocking
			a.Sprites = a.NextAnim.Sprites
			*a.NextAnim = Animation{}
		}
	default:
		break
	}
	return
}

func (a *Animation) ChangeAnimation(other Animation, blocking bool) (e error) {
	if len(other.Sprites) <= 0 {
		e = errors.New("need non empty animation")
		return
	}
	if a.Blocking {
		fmt.Println("Changing animation on animation end")
		a.NextAnim = &other
		a.NextAnim.Blocking = blocking
	} else {
		fmt.Println("Changing animation")
		a.Sprites = other.Sprites
		a.Blocking = blocking
		a.Cur = 0
	}
	return
}

func HandleAnimations(win *pixelgl.Window, state StateLock, anims map[string]Animation, currentAnims map[string]*Animation){
	center := pixel.ZV
	for _, player := range state.State.Players {
		transformation := pixel.IM.Rotated(center, player.Dir).Scaled(center, 0.5).Moved(player.Pos)
		movement := "idle"
		blocking := false
		switch {
		case player.Reload:
			movement = "reload"
			blocking = true
		case player.Shoot:
			movement = "shoot"
			blocking = true
		case player.Melee:
			movement = "melee"
			blocking = true
		case player.Moved:
			movement = "moved"
		default:
			movement = "idle"
		}
		prefix := Prefix("survivor", player.GetWeapon().Name, movement)

		anim, ok := currentAnims[player.Id]
		if !ok {
			newAnim, ok := anims[prefix]
			if ok {
				newAnim.Start(config.Conf.AnimationSpeed)
				currentAnims[player.Id] = &newAnim
				fmt.Println(newAnim.prefix, prefix)
				newAnim.prefix = prefix
			}
			anim = &newAnim
		}
		if anim.prefix != prefix {
			newAnim, found := anims[prefix]
			if found {
				fmt.Println(anim.prefix, prefix)
				anim.prefix = prefix
				anim.ChangeAnimation(newAnim, blocking)
			}
		}
		if len(anim.Sprites) > 0 {
			anim.Next().Draw(win, transformation)
		}
	}
	//for _, zombie := range state.State.Zombies {
	//	v, ok := currentAnims[zombie.Id]
	//	transformation := pixel.IM.Rotated(center, zombie.Dir).Moved(zombie.Pos)
	//	if !ok {
	//		v = anims[Prefix("zombie","idle")]
	//		currentAnims[zombie.Id] = v
	//	}
	//	v.Next().Draw(win, transformation)
	//}
	//for _, shoot := range state.State.Shoots {
	//	fmt.Println(shoot.Weapon)
	//}
}

func LoadMap(m model.Map) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	for _, w := range m.Walls {
		imd.Color = colornames.Black
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(w.P.X, w.P.Y), pixel.V(w.Q.X, w.Q.Y))
		imd.Line(w.Thickness)
	}
	return imd
}