package psystem

import "math"

// FloatPoint represents a point with x and y coordinates that are of float64 type
type FloatPoint struct {
	X, Y float64
}

// ZeroMagError represents a problem where a zero magnitude vector was encountered
type ZeroMagError struct{}

func (z ZeroMagError) Error() string {
	return "Zero magnitude vector"
}

// Add returns the vector p + q
func (p FloatPoint) Add(q FloatPoint) FloatPoint {
	return FloatPoint{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p - q
func (p FloatPoint) Sub(q FloatPoint) FloatPoint {
	return FloatPoint{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p * k
func (p FloatPoint) Mul(k float64) FloatPoint {
	return FloatPoint{p.X * k, p.Y * k}
}

// Div returns the vector p / k
func (p FloatPoint) Div(k float64) FloatPoint {
	return FloatPoint{p.X / k, p.Y / k}
}

// Mag returns the magnitude of p
func Mag(p FloatPoint) float64 {
	return math.Sqrt(math.Pow(p.X, 2.0) + math.Pow(p.Y, 2.0))
}

// Normalize returns a normalized vector of p
func Normalize(p FloatPoint) (FloatPoint, error) {
	mag := Mag(p)
	if mag != 0 {
		return FloatPoint{p.X / mag, p.Y / mag}, nil
	}
	return FloatPoint{0.0, 0.0}, ZeroMagError{}
}

// Dist returns distance from a to b
func Dist(a, b FloatPoint) float64 {
	return math.Sqrt(math.Pow(float64(a.X)-float64(b.X), 2.0) +
		math.Pow(float64(a.Y)-float64(b.Y), 2.0))
}

// DotProduct calculates dot product of points
func DotProduct(a, b FloatPoint) float64 {
	return a.X*b.X + a.Y*b.Y
}
