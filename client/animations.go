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
	"math/rand"
	"strconv"
	"math"
)

type Animation struct {
	Prefix   string
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
	if len(a.Sprites) > 1 {
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
	}
	return
}

func (a *Animation) ChangeAnimation(other Animation, blocking bool) (e error) {
	if len(other.Sprites) <= 0 {
		e = errors.New("need non empty animation")
		return
	}
	if a.Blocking {
		a.NextAnim = &other
		a.NextAnim.Blocking = blocking
	} else {
		a.Sprites = other.Sprites
		a.Blocking = blocking
		a.Cur = 0
	}
	return
}

func HandleAnimations(win *pixelgl.Window, state model.State, anims map[string]Animation, currentAnims map[string]*Animation){
	center := pixel.ZV

	bullet := anims["bullet"]
	for _, shot := range state.Shots {
		p := shot.GetPos(state.Timestamp)
		transformation := pixel.IM.Scaled(pixel.ZV, config.BulletScale).Rotated(pixel.ZV,shot.Angle - math.Pi/2).Moved(p)
		bullet.Next().Draw(win, transformation)
	}
	barrel:=anims["barrel"]
	for _, b := range state.Barrels{
		barrelx:=barrel.Next().Picture().Bounds().Max.X

		//b.GetHitBox()*2/barrelx

		transformation :=pixel.IM.ScaledXY(pixel.ZV,pixel.V(b.GetHitBox()*2/barrelx,b.GetHitBox()*2/barrelx)).Moved(b.Pos)
		barrel.Next().Draw(win, transformation)


	}
	for _, zombie := range state.Zombies {
		v, ok := currentAnims[zombie.Id]
		prefix := Prefix("zombie", "walk")
		if !ok{
			newanim, ok := anims[prefix]
			if ok {
				currentAnims[zombie.Id] = &newanim
				newanim.Start(config.Conf.AnimationSpeed)
				v = &newanim
			}else {
				fmt.Println("Did not find animation:",prefix)
				continue
			}
		}
		if zombie.Attacking {
			n := rand.Int()%3 + 1
			prefix = Prefix("zombie", "attack0"+strconv.Itoa(n))
		}

		if prefix != v.Prefix {
			if newanim, ok := anims[prefix]; ok{
				currentAnims[zombie.Id].ChangeAnimation(newanim, true)
				currentAnims[zombie.Id].Prefix = prefix
			}

		}
		if len(v.Sprites) > 0 {
			transformation := pixel.IM.Scaled(center, config.ZombieScale).Rotated(center, zombie.Dir).Moved(zombie.Pos)
			v.Next().Draw(win, transformation)
		}
	}
	//todo draw barrels
	/*for -,barrel := range state.Barrels{

	}*/
	for _, player := range state.Players {
		movement := "idle"
		blocking := false
		switch player.Action{
		case model.RELOAD:
			movement = "reload"
			blocking = true
		case model.SHOOT:
			movement = "shoot"
			blocking = true
		case model.MELEE:
			movement = "melee"
			blocking = true
		case model.IDLE:
			movement = "move"
		}

		prefix := Prefix("survivor", player.WeaponType.Name(), movement)
		anim, ok := currentAnims[player.Id]
		if !ok {
			newAnim, ok := anims[prefix]
			if ok {
				newAnim.Start(config.Conf.AnimationSpeed)
				currentAnims[player.Id] = &newAnim
				newAnim.Prefix = prefix
			}else{
				continue
			}
			anim = &newAnim
		}
		if anim.Prefix != prefix {
			newAnim, found := anims[prefix]
			if found {
				anim.Prefix = prefix
				anim.ChangeAnimation(newAnim, blocking)
			}
		}
		if len(anim.Sprites) > 0 {
			transformation := pixel.IM.Rotated(center, player.Dir).Scaled(center, config.HumanScale).Moved(player.Pos)
			anim.Next().Draw(win, transformation)
		}
	}
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

