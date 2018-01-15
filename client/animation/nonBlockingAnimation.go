package animation

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

type NonBlockingAnimation struct {
	Animation
	prefix         string
	Sprites        []*pixel.Sprite
	transformation Transformation
	animationSpeed AnimationSpeed
	cur            int
}

func (nba *NonBlockingAnimation) CurrentSprite() *pixel.Sprite {
	return nba.Sprites[nba.cur]
}

func (nba *NonBlockingAnimation) Prefix() string {
	return nba.prefix
}

func (nba *NonBlockingAnimation) Draw(win *pixelgl.Window) {
	nba.Sprites[nba.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, nba.transformation.Rotation).
			Scaled(pixel.ZV, nba.transformation.Scale).
			Moved(nba.transformation.Pos))
}

func (nba *NonBlockingAnimation) Next() Animation {
	inc := nba.animationSpeed.IncFrames()
	nba.cur = (nba.cur + inc) % len(nba.Sprites)
	return nba
}

func (nba *NonBlockingAnimation) ChangeAnimation(animation Animation) Animation {
	animation.SetTransformation(nba.transformation)
	return animation
}

func (nba *NonBlockingAnimation) SetTransformation(transformation Transformation) {
	nba.transformation = transformation
}

func (nba *NonBlockingAnimation) SetAnimationSpeed(duration time.Duration){
	nba.animationSpeed.Speed = duration
}

func (nba *NonBlockingAnimation) Copy() Animation{
	cpy := *nba
	return &cpy
}
func (nba *NonBlockingAnimation) SetDir(dir float64) {
	nba.transformation.Rotation = dir
}

func (nba *NonBlockingAnimation) SetPos(pos pixel.Vec) {
	nba.transformation.Pos = pos
}

func (nba *NonBlockingAnimation) SetScale(scale float64){
	nba.transformation.Scale = scale
}

func (nba *NonBlockingAnimation) GetTransformation()*Transformation{
	return &nba.transformation
}