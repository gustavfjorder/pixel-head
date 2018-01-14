package animation

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

type StillAnimation struct {
	prefix         string
	Sprite         *pixel.Sprite
	transformation Transformation
}

func (sa *StillAnimation) CurrentSprite() *pixel.Sprite {
	return sa.Sprite
}

func (sa *StillAnimation) Prefix() string {
	return sa.prefix
}

func (sa *StillAnimation) Draw(win *pixelgl.Window) {
	sa.Sprite.Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, sa.transformation.Rotation).
			Scaled(pixel.ZV, sa.transformation.Scale).
			Moved(sa.transformation.Pos))
}

func (sa *StillAnimation) Next() Animation {
	return sa
}

func (sa *StillAnimation) ChangeAnimation(animation Animation) Animation {
	animation.SetTransformation(sa.transformation)
	return animation
}

func (sa *StillAnimation) SetTransformation(transformation Transformation) {
	sa.transformation = transformation
}

func (sa *StillAnimation) SetAnimationSpeed(duration time.Duration){}

func (sa *StillAnimation) Copy() Animation{
	cpy := *sa
	return &cpy
}

func (sa *StillAnimation) SetDir(dir float64) {
	sa.transformation.Rotation = dir
}

func (sa *StillAnimation) SetPos(pos pixel.Vec) {
	sa.transformation.Pos = pos
}