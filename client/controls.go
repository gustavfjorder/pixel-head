package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type Control struct {
	Left  pixelgl.Button
	Right pixelgl.Button
	Up    pixelgl.Button
	Down  pixelgl.Button
}

func handleDir(prev pixel.Matrix) pixel.Matrix {
	angle, i := 0.0, 0
	if win.Pressed(config.Control.Up) {
		angle += math.Pi / 2
		i++
	}
	if win.Pressed(config.Control.Right) {
		if i <= 0 {
			angle += math.Pi * 2
		}
		i++
	}
	if win.Pressed(config.Control.Left) {
		angle += math.Pi
		i++
	}
	if win.Pressed(config.Control.Down) {
		angle += math.Pi * 3 / 2
		i++
	}
	if i <= 0 {
		return prev
	}
	return pixel.IM.Rotated(pixel.V(0, 0), angle/float64(i))
}
