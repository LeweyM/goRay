package Vector

import (
	"math"
	"testing"
)

func TestVector_Minus(t *testing.T) {
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

func TestVector_Product(t *testing.T) {
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

func TestVector_Normalize(t *testing.T) {
	type fields struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name   string
		fields fields
		want   Vector
	}{
		{
			name: "1 1 1",
			fields: fields{
				x: 1,
				y: 0,
				z: 0,
			},
			want: Vector{
				x: 1.0,
				y: 0.0,
				z: 0.0,
			},
		},
		{
			name: "1 3 6",
			fields: fields{
				x: 1,
				y: 3,
				z: 6,
			},
			want: Vector{
				x: 0.14744195615489714,
				y: 0.4423258684646914,
				z: 0.8846517369293828,
			},
		},
		{
			name: "1 1 1",
			fields: fields{
				x: 1,
				y: 1,
				z: 1,
			},
			want: Vector{
				x: 0.5773502691896258,
				y: 0.5773502691896258,
				z: 0.5773502691896258,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				x: tt.fields.x,
				y: tt.fields.y,
				z: tt.fields.z,
			}
			if got := v.Normalize(); !VectorIsEqual(&got, &tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Rotate(t *testing.T) {
	type args struct {
		rotationAxis Vector
		degrees      float64
	}
	tests := []struct {
		name string
		vector Vector
		args
		expected Vector
	}{
		{
			name: "facing z and rotate 90 degrees",
			vector: *New(0,0,1),
			args: args{
				degrees:      math.Pi/2,
			},
			expected: *New(1,0,0),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.vector.RotateY(test.degrees)
			if !VectorIsEqual(&got, &test.expected) {
				t.Errorf("RotateY() = %v, want %v", got, test.expected)
			}
		})
	}
}


func Benchmark_Normalize(b *testing.B) {
	v := New(2,4,1)
	for n := 0; n < b.N; n++ {
		v.Normalize()
	}
}

func Benchmark_Minus(b *testing.B) {
	v := New(0,0,1)
	v2 := New(1,3,1)
	for n := 0; n < b.N; n++ {
		v.Minus(*v2)
	}
}

func Benchmark_Rotate(b *testing.B) {
	v := New(0,0,1)
	for n := 0; n < b.N; n++ {
		v.RotateY(math.Pi)
	}
}

func Benchmark_ClosureRotate(b *testing.B) {
	v := *New(0,0,1)
	r := RotateYBuilder(math.Pi)
	for n := 0; n < b.N; n++ {
		r(v)
	}
}

func VectorIsEqual(v1 *Vector, v2 *Vector) bool {
	return withinErrorMargin(v1.X(), v2.X()) && withinErrorMargin(v1.Y(), v2.Y()) && withinErrorMargin(v1.Z(), v2.Z())
}

func withinErrorMargin(f1 float64, f2 float64) bool {
	return math.Abs(f1-f2) < 0.0000001
}

func TestVector_DistanceBetween(t *testing.T) {
	type fields struct {
		x float64
		y float64
		z float64
	}
	type args struct {
		vector Vector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "test",
			fields: fields{
				x: 1,
				y: 2,
				z: 3,
			},
			args: args{
				vector: Vector{
					x: 5,
					y: 4,
					z: 3,
				},
			},
			want: 4.47213595,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector{
				x: tt.fields.x,
				y: tt.fields.y,
				z: tt.fields.z,
			}
			if got := v.DistanceBetween(tt.args.vector); !withinErrorMargin(got, tt.want) {
				t.Errorf("DistanceBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}