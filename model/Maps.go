package model

import (
	"math"
	"github.com/faiface/pixel"
	"math/rand"
)

type Map struct {
	Walls []Wall
	SpawnPoint []pixel.Vec
	LootPoints []Point
	Bounds pixel.Rect
}

type Wall struct {
	Segment
	Thickness float64
}



var MapTemplates = map[string]Map{
	"Test1": {
		Walls: NewWallSeries(30, NewPoint(100, 100), NewPoint(100, 1000), NewPoint(1000, 1000), NewPoint(1000, 100)),
		SpawnPoint: randomSpawnPoints(pixel.R(100,100,1000,1000), 10),
		LootPoints: []Point{
			NewPoint(150, 150),
			NewPoint(150, 250),
			NewPoint(150, 350),
		},
		Bounds : pixel.R(100,100,1000,1000),
	},
}

func (w Wall) Intersect(l Segment) bool {
	var (
		angle = w.P.Angle(w.Q)
		v1    = pixel.V(w.Thickness/2, 0).Rotated(angle + math.Pi/2)
		v2    = pixel.V(w.Thickness/2, 0).Rotated(angle - math.Pi/2)
		pBR   = w.P.Add(NewPoint(v1.X, v1.Y))
		pBL   = w.P.Add(NewPoint(v2.X, v2.Y))
		pTR   = w.Q.Add(NewPoint(v1.X, v1.Y))
		pTL   = w.Q.Add(NewPoint(v2.X, v2.Y))
		lB    = NewLine(pBR, pBL)
		lT    = NewLine(pTR, pTL)
		lR    = NewLine(pBR, pTR)
		lL    = NewLine(pBL, pTL)
	)
	return l.Intersect(lB) ||
		l.Intersect(lT) ||
		l.Intersect(lR) ||
		l.Intersect(lL)
}

func NewWallSeries(thickness float64, points ...Point) []Wall {
	walls := make([]Wall, len(points))
	for i := 0; i < len(points); i++ {
		walls[i] = NewWall(points[i], points[(i+1)%len(points)], thickness)
	}
	return walls
}

func NewWall(p, q Point, thickness float64) Wall {
	return Wall{
		NewLine(p, q),
		thickness,
	}
}

func randomSpawnPoints(bounds pixel.Rect, n int) []pixel.Vec{
	res := make([]pixel.Vec,n)
	diffX := bounds.Max.X - bounds.Min.X
	diffY := bounds.Max.Y - bounds.Min.Y
	for i := range res {
		x := rand.Float64()*diffX + bounds.Min.X
		y := rand.Float64()*diffY + bounds.Min.Y
		res[i] = pixel.V(x,y)
	}
	return res
}
