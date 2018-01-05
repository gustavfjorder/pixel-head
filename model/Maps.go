package model

import (
	"math"
	"github.com/faiface/pixel"
)

type Map struct {
	Walls []Wall
}

type Point struct {
	X float64
	Y float64
}

type Line struct {
	P Point
	Q Point
}

type Wall struct {
	Line
	Thickness float64
}

var MapTemplates = map[string]Map{
	"Test1" : {
		Walls:[]Wall{
			NewWall(NewPoint(100,100), NewPoint(1000,100), 30),
			NewWall(NewPoint(1000,100), NewPoint(1000,1000), 30),
			NewWall(NewPoint(1000,1000), NewPoint(100,1000), 30),
			NewWall(NewPoint(100,1000), NewPoint(100,100), 30),
		},
	},
}

func (w Wall) Intersect(l Line) bool {
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

func NewWall(p,q Point, thickness float64) Wall{
	return Wall{
		NewLine(p,q),
		thickness,
	}
}

func NewLine(p, q Point) Line {
	return Line{P: p, Q: q}
}

func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

//returns angle in range [-Pi;Pi]
func (this Point) Angle(other Point) float64 {
	return math.Atan2(other.Y-this.Y, other.X-this.X)
}

func (this Point) Add(other Point) (Point) {
	return Point{
		X: this.X + other.X,
		Y: this.Y + other.Y,
	}
}

func (this Point) Dist(other Point) float64 {
	dx := other.X - this.X
	dy := other.Y - this.Y
	return math.Sqrt( dx*dx + dy*dy)
}

//implemented from https://www.geeksforgeeks.org/check-if-two-given-line-segments-intersect/
func (this Line) Intersect(other Line) bool {
	var (
		p1 = this.P
		p2 = this.Q
		q1 = other.P
		q2 = other.Q
		o1 = orientation(p1, q1, p2)
		o2 = orientation(p1, q1, q2)
		o3 = orientation(p2, q2, p1)
		o4 = orientation(p2, q2, q1)
	)
	return (o1 != o2 && o3 != o4) ||
		(o1 == 0 && onSegment(p1, p2, q1)) ||
		(o2 == 0 && onSegment(p1, q2, q1)) ||
		(o3 == 0 && onSegment(p2, p1, q2)) ||
		(o4 == 0 && onSegment(p2, q1, q2))
}

func onSegment(p, q, r Point) bool {
	return (q.X <= max(p.X, r.X) && q.X >= min(p.X, r.X)) &&
		(q.Y <= max(p.Y, r.Y) && q.X >= min(p.Y, r.Y))
}

// 0 --> p, q and r are colinear
// 1 --> Clockwise
// 2 --> Counterclockwise
func orientation(p, q, r Point) int {
	val := (q.Y-p.Y)*(r.X-q.X) - ( q.X-p.X)*(r.Y-q.Y)
	switch {
	case val == 0:
		return 0
	case val > 0:
		return 1
	default:
		return 2
	}
}

func max(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1
	}
	return v2
}

func min(v1, v2 float64) float64 {
	if v1 < v2 {
		return v1
	}
	return v2
}
