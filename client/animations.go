package client

import (
	"github.com/faiface/pixel"
	_ "image/png"
	"github.com/faiface/pixel/imdraw"
	"github.com/gustavfjorder/pixel-head/model"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"reflect"
	"fmt"
	"time"
	"math"
	"github.com/gustavfjorder/pixel-head/config"
)

type AnimationType int

const (
	NonBlocking AnimationType = iota
	Blocking
	Terminal
	Still
)

type Animation interface {
	Draw(win *pixelgl.Window)
	Next() Animation
	ChangeAnimation(animation Animation) Animation
	SetTransformation(transformation Transformation)
	SetAnimationSpeed(duration time.Duration)
	Prefix() string
	CurrentSprite() *pixel.Sprite
}

func NewAnimation(prefix string, sprites []*pixel.Sprite, animationType AnimationType, speeds ...time.Duration) Animation {
	if len(sprites) <= 0 {
		panic(errors.New("unable to make animation from no sprites for:" + prefix))
	}
	speed := config.Conf.AnimationSpeed
	if len(speeds) > 0{
		speed = speeds[0]
	}
	as := AnimationSpeed{Speed:speed}
	switch animationType {
	case NonBlocking:
		return &NonBlockingAnimation{
			prefix:  prefix,
			sprites: sprites,
			cur:     0,
			animationSpeed:as,
		}
	case Blocking:
		return &BLockingAnimation{
			prefix:        prefix,
			Sprites:       sprites,
			cur:           0,
			nextAnimation: nil,
			animationSpeed:as,
		}
	case Still:
		return &StillAnimation{
			prefix: prefix,
			Sprite: sprites[0],
		}
	case Terminal:
		return &TerminalAnimation{
			prefix:  prefix,
			cur:     0,
			Sprites: sprites,
			animationSpeed:as,
		}

	default:
		return &NonBlockingAnimation{
			prefix:  prefix,
			sprites: sprites,
			cur:     0,
		}
	}
}

type Transformation struct {
	Pos      pixel.Vec
	Scale    float64
	Rotation float64
}

type AnimationSpeed struct {
	Speed     time.Duration
	LastFrame time.Time
	diff      float64
}

func (as *AnimationSpeed) IncFrames() int {
	zerotime := time.Time{}
	if as.LastFrame == zerotime {
		as.LastFrame = time.Now()
		return 0
	}
	duration := time.Since(as.LastFrame)
	diff := duration.Seconds() / as.Speed.Seconds() + as.diff
	as.LastFrame = time.Now()
	frames := math.Floor(diff)
	as.diff = diff - frames
	return int(frames)
}

type NonBlockingAnimation struct {
	prefix         string
	sprites        []*pixel.Sprite
	transformation Transformation
	animationSpeed AnimationSpeed
	cur            int
}

func (nba NonBlockingAnimation) CurrentSprite() *pixel.Sprite {
	return nba.sprites[nba.cur]
}

func (nba NonBlockingAnimation) Prefix() string {
	return nba.prefix
}

func (nba NonBlockingAnimation) Draw(win *pixelgl.Window) {
	nba.sprites[nba.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, nba.transformation.Rotation).
			Scaled(pixel.ZV, nba.transformation.Scale).
			Moved(nba.transformation.Pos))
}

func (nba NonBlockingAnimation) Next() Animation {
	inc := nba.animationSpeed.IncFrames()
	nba.cur = (nba.cur + inc) % len(nba.sprites)
	return &nba
}

func (nba NonBlockingAnimation) ChangeAnimation(animation Animation) Animation {
	animation.SetTransformation(nba.transformation)
	return animation
}

func (nba *NonBlockingAnimation) SetTransformation(transformation Transformation) {
	nba.transformation = transformation
}

func (nba *NonBlockingAnimation) SetAnimationSpeed(duration time.Duration){
	nba.animationSpeed.Speed = duration
}

type BLockingAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	transformation Transformation
	cur            int
	nextAnimation  Animation
	animationSpeed AnimationSpeed
}

func (ba BLockingAnimation) CurrentSprite() *pixel.Sprite {
	return ba.Sprites[ba.cur]
}

func (ba BLockingAnimation) Prefix() string {
	return ba.prefix
}

func (ba BLockingAnimation) Draw(win *pixelgl.Window) {
	ba.Sprites[ba.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, ba.transformation.Rotation).
			Scaled(pixel.ZV, ba.transformation.Scale).
			Moved(ba.transformation.Pos))
}

func (ba BLockingAnimation) Next() Animation {
	inc := ba.cur + ba.animationSpeed.IncFrames()
	if inc >= len(ba.Sprites) && ba.nextAnimation != nil {
		ba.nextAnimation.SetTransformation(ba.transformation)
		return ba.nextAnimation
	}
	ba.cur = inc % len(ba.Sprites)
	return &ba
}

func (ba BLockingAnimation) ChangeAnimation(animation Animation) Animation {
	if reflect.TypeOf(animation) == reflect.TypeOf(&TerminalAnimation{}) {
		fmt.Println("Made terminal animation")
		animation.SetTransformation(ba.transformation)
		return animation
	}
	if ba.cur+1 >= len(ba.Sprites) {
		return animation
	}
	ba.nextAnimation = animation
	return &ba
}

func (ba *BLockingAnimation) SetTransformation(transformation Transformation) {
	ba.transformation = transformation
}

func (ba *BLockingAnimation) SetAnimationSpeed(duration time.Duration){
	ba.animationSpeed.Speed = duration
}

type StillAnimation struct {
	prefix         string
	Sprite         *pixel.Sprite
	transformation Transformation
}

func (sa StillAnimation) CurrentSprite() *pixel.Sprite {
	return sa.Sprite
}

func (sa StillAnimation) Prefix() string {
	return sa.prefix
}

func (sa StillAnimation) Draw(win *pixelgl.Window) {
	sa.Sprite.Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, sa.transformation.Rotation).
			Scaled(pixel.ZV, sa.transformation.Scale).
			Moved(sa.transformation.Pos))
}

func (sa StillAnimation) Next() Animation {
	return &sa
}

func (sa StillAnimation) ChangeAnimation(animation Animation) Animation {
	animation.SetTransformation(sa.transformation)
	return animation
}

func (sa *StillAnimation) SetTransformation(transformation Transformation) {
	sa.transformation = transformation
}

func (sa *StillAnimation) SetAnimationSpeed(duration time.Duration){}

type TerminalAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	transformation Transformation
	cur            int
	animationSpeed AnimationSpeed
}

func (ta TerminalAnimation) CurrentSprite() *pixel.Sprite {
	return ta.Sprites[ta.cur]
}

func (ta TerminalAnimation) Prefix() string {
	return ta.prefix
}

func (ta TerminalAnimation) Draw(win *pixelgl.Window) {
	ta.Sprites[ta.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, ta.transformation.Rotation).
			Scaled(pixel.ZV, ta.transformation.Scale).
			Moved(ta.transformation.Pos))
}

func (ta TerminalAnimation) Next() Animation {
	inc := ta.cur + ta.animationSpeed.IncFrames()
	if inc >= len(ta.Sprites) {
		fmt.Println("Removed")
		return nil
	}
	ta.cur = inc % len(ta.Sprites)
	return &ta
}

func (ta TerminalAnimation) ChangeAnimation(animation Animation) Animation {
	return &ta
}

func (ta *TerminalAnimation) SetTransformation(transformation Transformation) {
	ta.transformation = transformation
}

func (ta *TerminalAnimation) SetAnimationSpeed(duration time.Duration){
	ta.animationSpeed.Speed = duration
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
