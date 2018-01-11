package psystem

import (
	"math"
	"runtime/debug"
	"testing"
)

func assertEqual(a interface{}, b interface{}, t *testing.T) {
	if a != b {
		debug.PrintStack()
		t.Errorf("Got: %v, Expected: %v", a, b)
	}
}

func TestError(t *testing.T) {
	assertEqual(ZeroMagError{}.Error(), "Zero magnitude vector", t)
}

var values = []FloatPoint{FloatPoint{0, 0}, FloatPoint{3, 4},
	FloatPoint{7, 10}, FloatPoint{0, 0},
	FloatPoint{-5, 5}, FloatPoint{5, 5},
	FloatPoint{5, 5}, FloatPoint{-5, 5},
	FloatPoint{6, 7}, FloatPoint{12, 10}}

func TestDist(t *testing.T) {
	answers := []float64{5.0, 12.2, 10.0, 10.0, 6.7}

	aIndex := 0
	for i := 0; i < len(values)-1; i += 2 {
		response := float64(math.Trunc(Dist(values[i], values[i+1])*100) / 100)
		assertEqual(response, answers[aIndex], t)
		aIndex++
	}
}

func TestDotProduct(t *testing.T) {
	answers := []float64{0, 0, 0, 0, 142}

	aIndex := 0
	for i := 0; i < len(values)-1; i += 2 {
		response := DotProduct(values[i], values[i+1])
		assertEqual(response, answers[aIndex], t)
		aIndex++
	}
}

func TestAdd(t *testing.T) {
	for i := -50.0; i <= 50.0; i += 0.1 {
		for j := -50.0; j >= -50.0; j -= 0.1 {
			assertEqual(FloatPoint{i, j}.Add(FloatPoint{j, i}), FloatPoint{i + j, i + j}, t)
		}
	}
}

func TestSub(t *testing.T) {
	for i := -50.0; i <= 50.0; i += 0.1 {
		for j := -50.0; j >= -50.0; j -= 0.1 {
			assertEqual(FloatPoint{i, j}.Sub(FloatPoint{j, i}), FloatPoint{i - j, j - i}, t)
		}
	}
}

func TestMul(t *testing.T) {
	for i := -50.0; i <= 50.0; i += 0.1 {
		for j := -50.0; j >= -50.0; j -= 0.1 {
			assertEqual(FloatPoint{i, j}.Mul(i), FloatPoint{i * i, j * i}, t)
		}
	}
}

func TestDiv(t *testing.T) {
	for i := -50.0; i <= 50.0; i += 0.1 {
		for j := -50.0; j >= -50.0; j -= 0.1 {
			assertEqual(FloatPoint{i, j}.Div(i), FloatPoint{i / i, j / i}, t)
		}
	}
}

func TestMag(t *testing.T) {
	assertEqual(Mag(FloatPoint{0.0, 0.0}), 0.0, t)
	assertEqual(Mag(FloatPoint{1.0, 0.0}), 1.0, t)
	assertEqual(Mag(FloatPoint{0.0, 1.0}), 1.0, t)
	assertEqual(Mag(FloatPoint{5.0, 20.0}), math.Sqrt(425.0), t)
	assertEqual(math.Trunc(Mag(FloatPoint{0.1, 0.2})*1000),
		math.Trunc(math.Sqrt(0.05)*1000), t)
	assertEqual(Mag(FloatPoint{-3.0, 4.0}), 5.0, t)
	assertEqual(Mag(FloatPoint{-3.0, -4.0}), 5.0, t)
}

func TestNormalize(t *testing.T) {
	val, err := Normalize(FloatPoint{0.0, 0.0})
	assertEqual(err, ZeroMagError{}, t)
	assertEqual(val, FloatPoint{0.0, 0.0}, t)

	val, err = Normalize(FloatPoint{1.0, 0.0})
	assertEqual(err, nil, t)
	assertEqual(val, FloatPoint{1.0, 0.0}, t)

	val, err = Normalize(FloatPoint{3.0, 6.0})
	assertEqual(err, nil, t)
	assertEqual(val, FloatPoint{3.0 / math.Sqrt(45.0), 6.0 / math.Sqrt(45.0)}, t)

	val, err = Normalize(FloatPoint{-3.0, -6.0})
	assertEqual(err, nil, t)
	assertEqual(val, FloatPoint{-3.0 / math.Sqrt(45.0), -6.0 / math.Sqrt(45.0)}, t)

	val, err = Normalize(FloatPoint{-3.0, 6.0})
	assertEqual(err, nil, t)
	assertEqual(val, FloatPoint{-3.0 / math.Sqrt(45.0), 6.0 / math.Sqrt(45.0)}, t)
}
