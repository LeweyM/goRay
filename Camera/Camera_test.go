package Camera

import (
	"goRay/Object"
	"goRay/Ray"
	"goRay/Vector"
	"image/color"
	"math"
	"reflect"
	"testing"
)

func TestCamera_GetPixelHeadingVector(t *testing.T) {
	type args struct {
		pHeight  float64
		pWidth   float64
		unitZoom float64
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
				unitZoom: 1.0,
			},
			want: Vector.New(0, 0, 1),
		},
		{
			name: "looking 90degrees right from z axis",
			args: args{
				pHeight: 0.0,
				pWidth:  1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees left from z axis",
			args: args{
				pHeight: 0.0,
				pWidth:  -1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(-0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees up from z axis",
			args: args{
				pHeight: 1.0,
				pWidth:  0.0,
				unitZoom: 1.0,
			},
			want: Vector.New(0, 0.7071067811, 0.7071067811),
		},
		{
			name: "looking 90degrees up and left from z axis",
			args: args{
				pHeight: 1.0,
				pWidth:  -1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(-0.5, 0.7071067811, 0.5), // BUT WHY THO???
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Camera{}
			got := c.GetPixelHeadingVector(tt.args.pHeight, tt.args.pWidth, tt.args.unitZoom)

			if vectorIsEqual(got, tt.want) {
				t.Errorf("GetPixelHeadingVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getScreenMatrix(t *testing.T) {
	type args struct {
		height float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "scale 3",
			args: args{
				height: 3,
			},
			want: []float64{-1, 0, 1},
		},
		{
			name: "scale 5",
			args: args{
				height: 5,
			},
			want: []float64{-2, -1, 0, 1, 2},
		},
		{
			name: "scale 4",
			args: args{
				height: 4,
			},
			want: []float64{-1.5, -0.5, 0.5, 1.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getScreenMatrix(tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScreenMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCamera_GetPixel(t *testing.T) {
	zRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(0, 0, 1))
	type fields struct {
		widthMatrix  []float64
		heightMatrix []float64
		height       int
		width        int
		origin       Vector.Vector
		objectList   []Object.Object
		pixelList    []Pixel
	}
	type args struct {
		x   int
		y   int
		ray Ray.Ray
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Pixel
	}{
		{
			name: "without objects",
			fields: fields{
				origin: Vector.Vector{},
			},
			args: args{
				x:   0,
				y:   0,
				ray: zRay,
			},
			want: Pixel{
				color: color.Black,
				x:     0,
				y:     0,
			},
		},
		{
			name: "with an object in front",
			fields: fields{
				origin:     Vector.Vector{},
				objectList: []Object.Object{Object.NewSphere(*Vector.New(0, 0, 5), 2, )},
			},
			args: args{
				x:   0,
				y:   0,
				ray: zRay,
			},
			want: Pixel{
				color: color.White,
				x:     0,
				y:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Camera{
				widthMatrix:  tt.fields.widthMatrix,
				heightMatrix: tt.fields.heightMatrix,
				height:       tt.fields.height,
				width:        tt.fields.width,
				origin:       tt.fields.origin,
				objectList:   tt.fields.objectList,
				pixelList:    tt.fields.pixelList,
			}
			if got := c.GetPixel(tt.args.x, tt.args.y, tt.args.ray); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPixel() = %v, want %v", got, tt.want)
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
