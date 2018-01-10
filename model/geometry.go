package model

import (
	"math"
	"github.com/faiface/pixel"
)

type Point struct {
	X float64
	Y float64
}

type Segment struct {
	P Point
	Q Point
}

type Line struct {
	Slope     float64
	Intercept float64
	X         float64
	Vertical  bool
}

func NewLine(p, q Point) Segment {
	return Segment{P: p, Q: q}
}

func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

func PointFrom(v pixel.Vec) Point {
	return Point{X: v.X, Y: v.Y}
}

func (s Segment) Line() (line Line) {
	if s.Q.X == s.P.X {
		line.X = s.Q.X
		line.Vertical = true
		return line
	}
	left := s.Q
	right := s.P
	if s.P.X < s.Q.X {
		left = s.P
		right = s.Q
	}
	dx := right.X - left.X
	dy := right.Y - left.Y
	slope := dy / dx
	intercept := (-left.X)*slope + left.Y
	line = Line{
		Intercept: intercept,
		Slope:     slope,
		Vertical:  false,
	}
	return line
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
	return math.Sqrt(dx*dx + dy*dy)
}

func (this Segment) Intersect(other Segment) bool {
	var x, y float64
	thisLine := this.Line()
	otherLine := other.Line()
	switch {
	case !thisLine.Vertical:
		if !otherLine.Vertical {
			if thisLine.Slope == otherLine.Slope{
				return thisLine.Intercept == otherLine.Intercept
			}
			x = -(thisLine.Intercept - otherLine.Intercept) / (thisLine.Slope - otherLine.Slope)
			y = thisLine.Intercept + thisLine.Slope * x
		} else {
			x = otherLine.X
			y = thisLine.Intercept + thisLine.Slope * x
		}
	case !otherLine.Vertical:
		x = thisLine.X
		y = otherLine.Intercept + otherLine.Slope * x
	default:
		return otherLine.X == thisLine.X && (this.hasInRange(other.Q) || this.hasInRange(other.P))
	}
	return this.hasInRange(Point{x,y}) && other.hasInRange(Point{x,y})
}

func (s Segment) hasInRange(p Point) bool {
	minX := min(s.P.X, s.Q.X)
	maxX := max(s.P.X, s.Q.X)
	minY := min(s.P.Y, s.Q.Y)
	maxY := max(s.P.Y, s.Q.Y)
	return minX <= p.X && maxX >= p.X && minY <= p.Y && maxY >= p.Y
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

func minInt(v1,v2 int) int {
	if v1 < v2{
		return v1
	}
	return v2
}
