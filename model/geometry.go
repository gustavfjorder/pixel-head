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

func (segment Segment) Line() (line Line) {
	if segment.Q.X == segment.P.X {
		line.X = segment.Q.X
		line.Vertical = true
		return line
	}
	left := segment.Q
	right := segment.P
	if segment.P.X < segment.Q.X {
		left = segment.P
		right = segment.Q
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
func (point Point) Angle(other Point) float64 {
	return math.Atan2(other.Y-point.Y, other.X-point.X)
}

func (point Point) Add(other Point) (Point) {
	return Point{
		X: point.X + other.X,
		Y: point.Y + other.Y,
	}
}

func (point Point) Dist(other Point) float64 {
	dx := other.X - point.X
	dy := other.Y - point.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (segment Segment) Intersect(other Segment) bool {
	var x, y float64
	thisLine := segment.Line()
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
		return otherLine.X == thisLine.X && (segment.hasInRange(other.Q) || segment.hasInRange(other.P))
	}
	return segment.hasInRange(Point{x,y}) && other.hasInRange(Point{x,y})
}

func (segment Segment) hasInRange(p Point) bool {
	minX := min(segment.P.X, segment.Q.X)
	maxX := max(segment.P.X, segment.Q.X)
	minY := min(segment.P.Y, segment.Q.Y)
	maxY := max(segment.P.Y, segment.Q.Y)
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
