package client

import (
	"errors"
	"github.com/faiface/pixel"
	_ "image/png"
	"time"
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
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