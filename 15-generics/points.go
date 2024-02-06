package main

import (
	"fmt"
	"math"
)

type Pair[T fmt.Stringer] struct {
	Val1 T
	Val2 T
}

type Differ[T any] interface {
	fmt.Stringer
	Diff(T) float64
}

func FindCloserPair[T Differ[T]](pair1, pair2 Pair[T]) Pair[T] {
	d1 := pair1.Val1.Diff(pair1.Val2)
	d2 := pair2.Val1.Diff(pair1.Val2)
	if d1 < d2 {
		return pair1
	}
	return pair2
}

type Point2D struct {
	X, Y float64
}

func (p Point2D) String() string {
	return fmt.Sprintf("{%f,%f}", p.X, p.Y)
}

func (p Point2D) Diff(from Point2D) float64 {
	x := p.X - from.X
	y := p.Y - from.Y
	return math.Sqrt(x*x + y*y)
}

type Point3D struct {
	X, Y, Z float64
}

func (p Point3D) String() string {
	return fmt.Sprintf("{%f,%f,%f}", p.X, p.Y, p.Z)
}

func (p Point3D) Diff(from Point3D) float64 {
	x := p.X - from.X
	y := p.Y - from.Y
	z := p.Z - from.Z
	return math.Sqrt(x*x + y*y + z*z)
}
