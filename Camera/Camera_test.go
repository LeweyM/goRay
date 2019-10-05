package Camera

import (
	"goRay/Vector"
	"math"
	"testing"
)

func TestCamera_GetPixelHeadingVector(t *testing.T) {
	type args struct {
		pHeight float64
		pWidth  float64
	}
	tests := []struct {
		name string
		args args
		want *Vector.Vector
	}{
		{
			name: "looking straight along z axis",
			args: args{
				pHeight: 0.0,
				pWidth:  0.0,
			},
			want: Vector.New(0, 0, 1),
		},
		{
			name: "looking 90degrees right from z axis",
			args: args{
				pHeight: 0.0,
				pWidth:  1.0,
			},
			want: Vector.New(0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees left from z axis",
			args: args{
				pHeight: 0.0,
				pWidth:  -1.0,
			},
			want: Vector.New(-0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees up from z axis",
			args: args{
				pHeight: 1.0,
				pWidth:  0.0,
			},
			want: Vector.New(0, 0.7071067811, 0.7071067811),
		},
		{
			name: "looking 90degrees up and left from z axis",
			args: args{
				pHeight: 1.0,
				pWidth:  -1.0,
			},
			want: Vector.New(-0.5, 0.7071067811, 0.5), // BUT WHY THO???
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Camera{}
			got := c.GetPixelHeadingVector(tt.args.pHeight, tt.args.pWidth)

			if vectorIsEqual(got, tt.want) {
				t.Errorf("GetPixelHeadingVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func vectorIsEqual(v1 *Vector.Vector, v2 *Vector.Vector) bool {
	return withinErrorMargin(v1.X(), v2.X()) || withinErrorMargin(v1.Y(), v2.Y()) || withinErrorMargin(v1.Z(), v2.Z())
}

func withinErrorMargin(f1 float64, f2 float64) bool {
	return math.Abs(f1-f2) > 0.00000001
}
