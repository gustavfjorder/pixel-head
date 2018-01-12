package client

import (
	"errors"
	"github.com/faiface/pixel"
	_ "image/png"
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
)

type Animation struct {
	Prefix         string
	Sprites        []*pixel.Sprite
	Cur            int
	NextAnim       *Animation
	Transformation pixel.Matrix
	Blocking       bool
	Terminal       bool
	Finished       bool
}

func (a *Animation) Draw(win *pixelgl.Window){
	if a.Cur >= len(a.Sprites) {
		a.Cur = 0
	}
	a.Sprites[a.Cur].Draw(win, a.Transformation)
}

//inc is one element controlling how many frames to move in the animation (will only use first argument)
func (a *Animation) Next() {
	a.Finished = a.Terminal && a.Cur+1 >= len(a.Sprites)
	if a.Finished {
		return
	}
	if !a.Terminal && a.Cur+1 >= len(a.Sprites) && a.NextAnim != nil && len(a.NextAnim.Sprites) > 0 {
		a.Blocking = a.NextAnim.Blocking
		a.Sprites = a.NextAnim.Sprites
		*a.NextAnim = Animation{}
	}
	a.Cur = (a.Cur + 1) % len(a.Sprites)
}

func (a *Animation) ChangeAnimation(other Animation, blocking, terminal bool) (e error) {
	if a.Terminal || (a.NextAnim == nil && a.NextAnim.Terminal) {
		e = errors.New("cannot change terminal animation")
		return
	}
	if len(other.Sprites) <= 0 {
		e = errors.New("need non empty animation")
		return
	}
	if a.Blocking {
		a.NextAnim = &other
		a.NextAnim.Blocking = blocking
		a.NextAnim.Terminal = terminal
		a.NextAnim.Transformation = a.Transformation
	} else {
		a.Sprites = other.Sprites
		a.Blocking = blocking
		a.Terminal = terminal
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
