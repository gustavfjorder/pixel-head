package animation

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
	"time"
)

type TerminalAnimation struct {
	prefix         string
	Sprites        []*pixel.Sprite
	transformation Transformation
	cur            int
	animationSpeed AnimationSpeed
}

func (ta *TerminalAnimation) CurrentSprite() *pixel.Sprite {
	return ta.Sprites[ta.cur]
}

func (ta *TerminalAnimation) Prefix() string {
	return ta.prefix
}

func (ta *TerminalAnimation) Draw(win *pixelgl.Window) {
	ta.Sprites[ta.cur].Draw(win,
		pixel.IM.
			Rotated(pixel.ZV, ta.transformation.Rotation).
			Scaled(pixel.ZV, ta.transformation.Scale).
			Moved(ta.transformation.Pos))
}

func (ta *TerminalAnimation) Next() Animation {
	inc := ta.cur + ta.animationSpeed.IncFrames()
	if inc >= len(ta.Sprites) {
		fmt.Println("Removed")
		return nil
	}
	ta.cur = inc % len(ta.Sprites)
	return ta
}

func (ta *TerminalAnimation) ChangeAnimation(animation Animation) Animation {
	return ta
}

func (ta *TerminalAnimation) SetTransformation(transformation Transformation) {
	ta.transformation = transformation
}

func (ta *TerminalAnimation) SetAnimationSpeed(duration time.Duration){
	ta.animationSpeed.Speed = duration
}

func (ta *TerminalAnimation) Copy() Animation{
	cpy := *ta
	return &cpy
}

func (ta *TerminalAnimation) SetDir(dir float64) {
	ta.transformation.Rotation = dir
}

func (ta *TerminalAnimation) SetPos(pos pixel.Vec) {
	ta.transformation.Pos = pos
}