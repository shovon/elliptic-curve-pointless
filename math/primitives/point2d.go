package primitives

type Point2D[T any] [2]T

func (p Point2D[T]) X() T {
	return p[0]
}

func (p Point2D[T]) Y() T {
	return p[1]
}
