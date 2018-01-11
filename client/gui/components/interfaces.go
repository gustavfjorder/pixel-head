package components

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type component interface {
	Draw(target pixel.Target, pos pixel.Vec, center ...bool)
	//Target() *pixel.Batch
}

type clickable interface {
	OnLeftMouseClick(win *pixelgl.Window, handler func())
	OnRightMouseClick(win *pixelgl.Window, handler func())
	OnClick(win *pixelgl.Window, handler func(button pixelgl.Button))
}
