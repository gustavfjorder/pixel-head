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

type AnimationType int

const(
	NonBlocking AnimationType = iota
	Blocking
	Terminal
	Image
)

type Animation struct {
	Prefix   string
	Sprites  []*pixel.Sprite
	Cur      int
	NextAnim *Animation
	Pos      pixel.Vec
	Rotation float64
	Scale    float64
	Blocking bool
	Terminal bool
	Finished bool
}

func NewAnimation(prefix string, sprites []*pixel.Sprite) Animation{
	return Animation{Prefix:prefix,Sprites:sprites}
}

func (a *Animation) Draw(win *pixelgl.Window) {
	a.Sprites[a.Cur].Draw(win, pixel.IM.Rotated(pixel.ZV, a.Rotation).Scaled(pixel.ZV, a.Scale).Moved(a.Pos))
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
	if a.Terminal || (a.NextAnim != nil && a.NextAnim.Terminal) {
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
		a.NextAnim.Scale = a.Scale
		a.NextAnim.Pos = a.Pos
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
