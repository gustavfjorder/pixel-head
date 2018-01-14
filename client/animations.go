package client

import (
	"github.com/faiface/pixel"
	_ "image/png"
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
)

type AnimationType int

const (
	NonBlocking AnimationType = iota
	Blocking
	Terminal
	Image
)

type Animation interface {
	Draw(win *pixelgl.Window)
	Next() Animation
	ChangeAnimation(animation Animation) Animation
	SetTransformation(transformation Transformation)
	Prefix() string
	CurrentSprite() *pixel.Sprite
}

func NewAnimation(prefix string, sprites []*pixel.Sprite, animationType AnimationType) Animation {
	switch animationType {
	case NonBlocking:
		return &NonBlockingAnimation{
			prefix:prefix,
			Sprites:sprites,
			Cur : 0,
		}
	case Blocking:
		return &BLockingAnimation{
			prefix:prefix,
			Sprites:sprites,
			Cur:0,
		}
	default:
		return &NonBlockingAnimation{
			prefix:prefix,
			Sprites:sprites,
			Cur : 0,
		}
	}
}

type Transformation struct {
	Pos      pixel.Vec
	Scale    float64
	Rotation float64
}

type NonBlockingAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	Transformation Transformation
	Cur            int
}

func (nba NonBlockingAnimation) CurrentSprite() *pixel.Sprite{
	return nba.Sprites[nba.Cur]
}

func (nba NonBlockingAnimation) Prefix() string{
	return nba.prefix
}

func (nba NonBlockingAnimation) Draw(win *pixelgl.Window) {
	fmt.Println(nba.Transformation)
	nba.Sprites[nba.Cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, nba.Transformation.Rotation).
				Scaled(pixel.ZV, nba.Transformation.Scale).
					Moved(nba.Transformation.Pos))
}

func (nba NonBlockingAnimation) Next() Animation{
	nba.Cur = (nba.Cur + 1)%len(nba.Sprites)
	return &nba
}

func (nba NonBlockingAnimation) ChangeAnimation(animation Animation) Animation {
	return animation
}

func (nba *NonBlockingAnimation) SetTransformation(transformation Transformation){
 nba.Transformation = transformation
}

type BLockingAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	Transformation Transformation
	Cur            int
	NextAnimation  *Animation
}

func (ba BLockingAnimation) CurrentSprite() *pixel.Sprite{
	return ba.Sprites[ba.Cur]
}

func (ba BLockingAnimation) Prefix() string{
	return ba.prefix
}

func (ba BLockingAnimation) Draw(win *pixelgl.Window) {
	ba.Sprites[ba.Cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, ba.Transformation.Rotation).
			Scaled(pixel.ZV, ba.Transformation.Scale).
			Moved(ba.Transformation.Pos))
}

func (ba BLockingAnimation) Next() Animation{
	next := (ba.Cur + 1)%len(ba.Sprites)
	if next <= ba.Cur && ba.NextAnimation != nil{
		return *ba.NextAnimation
	}
	ba.Cur = next
	return &ba
}

func (ba BLockingAnimation) ChangeAnimation(animation Animation) Animation {
	if ba.Cur + 1 >= len(ba.Sprites){
		return animation
	}
	ba.NextAnimation = &animation
	return &ba
}

func (ba *BLockingAnimation) SetTransformation(transformation Transformation){
	ba.Transformation = transformation
}

//func (a *Animation) Draw(win *pixelgl.Window) {
//	a.Sprites[a.Cur].Draw(win, pixel.IM.Rotated(pixel.ZV, a.Rotation).Scaled(pixel.ZV, a.Scale).Moved(a.Pos))
//}
//
////inc is one element controlling how many frames to move in the animation (will only use first argument)
//func (a *Animation) Next() {
//	a.Finished = a.Terminal && a.Cur+1 >= len(a.Sprites)
//	if a.Finished {
//		return
//	}
//	if !a.Terminal && a.Cur+1 >= len(a.Sprites) && a.NextAnim != nil && len(a.NextAnim.Sprites) > 0 {
//		a.Blocking = a.NextAnim.Blocking
//		a.Sprites = a.NextAnim.Sprites
//		*a.NextAnim = Animation{}
//	}
//	a.Cur = (a.Cur + 1) % len(a.Sprites)
//}
//
//func (a *Animation) ChangeAnimation(other Animation, blocking, terminal bool) (e error) {
//	if a.Terminal || (a.NextAnim != nil && a.NextAnim.Terminal) {
//		e = errors.New("cannot change terminal animation")
//		return
//	}
//	if len(other.Sprites) <= 0 {
//		e = errors.New("need non empty animation")
//		return
//	}
//	if a.Blocking {
//		a.NextAnim = &other
//		a.NextAnim.Blocking = blocking
//		a.NextAnim.Terminal = terminal
//		a.NextAnim.Scale = a.Scale
//		a.NextAnim.Pos = a.Pos
//	} else {
//		a.Sprites = other.Sprites
//		a.Blocking = blocking
//		a.Terminal = terminal
//		a.Cur = 0
//	}
//	return
//}

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
