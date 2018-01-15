package animation

import (
	"github.com/faiface/pixel"
	_ "image/png"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/errors"
	"time"
	"math"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
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
	SetDir(dir float64)
	SetPos(pos pixel.Vec)
	SetScale(scale float64)
	SetAnimationSpeed(duration time.Duration)
	Prefix() string
	CurrentSprite() *pixel.Sprite
	Copy() Animation
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
