package animation

import (
	"reflect"
	"fmt"
)

//Animation that will block until it has shown all frames
//And upon finishing it will change to the next animation
//If there has been a ChangeAnimation call.
//(will change immediately if animation is terminal
type BLockingAnimation struct {
	NonBlockingAnimation
	nextAnimation  Animation
}

//Will go to the next frame of the animation and return self if
// not finished, otherwise it will return the pending animation
func (ba *BLockingAnimation) Next() Animation {
	inc := ba.cur + ba.animationSpeed.IncFrames()
	if inc >= len(ba.NonBlockingAnimation.Sprites) && ba.nextAnimation != nil {
		ba.nextAnimation.SetTransformation(ba.transformation)
		return ba.nextAnimation
	}
	ba.cur = inc % len(ba.Sprites)
	return ba
}

//Will change animation immediately if of type terminal or at last
// frame Otherwise will be returned upon last frame in Next()

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

func (ba *BLockingAnimation) Copy() Animation{
	cpy := *ba
	return &cpy
}