package Vector

import "testing"

func TestMinus(t *testing.T) {
	tests := []struct {
		v1     Vector
		v2     Vector
		result Vector
	}{
		{
			v1:     *New(10, 10, 10),
			v2:     *New(2, 2, 2),
			result: *New(8, 8, 8),
		}, {
			v1:     *New(10, 5, 2),
			v2:     *New(3, 4, 1),
			result: *New(7, 1, 1),
		},
	}
	for i, test := range tests {
		if test.v1.Minus(test.v2) != test.result {
			t.Errorf("test %d: expected %v, got %v", i, test.result, test.v1.Minus(test.v2))
		}
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		v1     Vector
		v2     Vector
		result float64
	}{
		{
			v1:     *New(10, 10, 10),
			v2:     *New(2, 2, 2),
			result: 60,
		}, {
			v1:     *New(10, 5, 2),
			v2:     *New(3, 4, 1),
			result: 52,
		},
	}
	for i, test := range tests {
		if test.v1.Dot(test.v2) != test.result {
			t.Errorf("test %d: expected %v, got %v", i, test.result, test.v1.Dot(test.v2))
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		v      Vector
		result Vector
	}{
		{
			v:      *New(0, 0, 0),
			result: *New(0.0, 0.0, 0.0),
		},
		{
			v:      *New(1, 1, 1),
			result: *New(0.5773502691896258, 0.5773502691896258, 0.5773502691896258),
		},
		{
			v:      *New(5, 3, 2),
			result: *New(0.8111071056538127, 0.48666426339228763, 0.3244428422615251),
		},
	}
	for i, test := range tests {
		if test.v.Normalize() != test.result {
			t.Errorf("test %d: expected %v, got %v", i, test.result, test.v.Normalize())
		}
	}
}
