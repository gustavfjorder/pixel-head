package animation

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"reflect"
	"fmt"
	"time"
)

type BLockingAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	transformation Transformation
	cur            int
	nextAnimation  Animation
	animationSpeed AnimationSpeed
}

func (ba *BLockingAnimation) CurrentSprite() *pixel.Sprite {
	return ba.Sprites[ba.cur]
}

func (ba *BLockingAnimation) Prefix() string {
	return ba.prefix
}

func (ba *BLockingAnimation) Draw(win *pixelgl.Window) {
	ba.Sprites[ba.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, ba.transformation.Rotation).
			Scaled(pixel.ZV, ba.transformation.Scale).
			Moved(ba.transformation.Pos))
}

func (ba *BLockingAnimation) Next() Animation {
	inc := ba.cur + ba.animationSpeed.IncFrames()
	if inc >= len(ba.Sprites) && ba.nextAnimation != nil {
		ba.nextAnimation.SetTransformation(ba.transformation)
		return ba.nextAnimation
	}
	ba.cur = inc % len(ba.Sprites)
	return ba
}

func (ba *BLockingAnimation) ChangeAnimation(animation Animation) Animation {
	if reflect.TypeOf(animation) == reflect.TypeOf(&TerminalAnimation{}) {
		fmt.Println("Made terminal animation")
		animation.SetTransformation(ba.transformation)
		return animation
	}
	if ba.cur+1 >= len(ba.Sprites) {
		return animation
	}
	ba.nextAnimation = animation
	return ba
}

func (ba *BLockingAnimation) SetTransformation(transformation Transformation) {
	ba.transformation = transformation
}

func (ba *BLockingAnimation) SetAnimationSpeed(duration time.Duration){
	ba.animationSpeed.Speed = duration
}
func (ba *BLockingAnimation) Copy() Animation{
	cpy := *ba
	return &cpy
}

func (ba *BLockingAnimation) SetDir(dir float64) {
	ba.transformation.Rotation = dir
}

func (ba *BLockingAnimation) SetPos(pos pixel.Vec) {
	ba.transformation.Pos = pos
}

func (ba *BLockingAnimation) SetScale(scale float64){
	ba.transformation.Scale = scale
}