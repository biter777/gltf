package gltf

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestBuffer_IsEmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		b    *Buffer
		want bool
	}{
		{"embedded", &Buffer{URI: "data:application/octet-stream;base64,dsjdsaGGUDXGA"}, true},
		{"external", &Buffer{URI: "https://web.com/a"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IsEmbeddedResource(); got != tt.want {
				t.Errorf("Buffer.IsEmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_EmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		b    *Buffer
		want string
	}{
		{"base", &Buffer{Data: []byte("any + old & data")}, "data:application/octet-stream;base64,YW55ICsgb2xkICYgZGF0YQ=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.EmbeddedResource()
			if got := tt.b.URI; got != tt.want {
				t.Errorf("Buffer.EmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_IsEmbeddedResource(t *testing.T) {
	tests := []struct {
		name string
		im   *Image
		want bool
	}{
		{"png", &Image{URI: "data:image/png;base64,dsjdsaGGUDXGA"}, true},
		{"jpg", &Image{URI: "data:image/png;base64,dsjdsaGGUDXGA"}, true},
		{"external", &Image{URI: "https://web.com/a"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.im.IsEmbeddedResource(); got != tt.want {
				t.Errorf("Image.IsEmbeddedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_marshalData(t *testing.T) {
	tests := []struct {
		name    string
		b       *Buffer
		want    []byte
		wantErr bool
	}{
		{"error", &Buffer{URI: "data:application/octet-stream;base64,_"}, nil, true},
		{"external", &Buffer{URI: "http://web.com"}, nil, false},
		{"empty", &Buffer{URI: "data:application/octet-stream;base64,"}, nil, false},
		{"test", &Buffer{URI: "data:application/octet-stream;base64,TEST"}, []byte{76, 68, 147}, false},
		{"complex", &Buffer{URI: "data:application/octet-stream;base64,YW55IGNhcm5hbCBwbGVhcw=="}, []byte{97, 110, 121, 32, 99, 97, 114, 110, 97, 108, 32, 112, 108, 101, 97, 115}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.marshalData()
			if (err != nil) != tt.wantErr {
				t.Errorf("Buffer.marshalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buffer.marshalData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_MatrixOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [16]float64
	}{
		{"default", &Node{Matrix: DefaultMatrix}, DefaultMatrix},
		{"zeros", &Node{Matrix: emptyMatrix}, DefaultMatrix},
		{"other", &Node{Matrix: [16]float64{2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2}}, [16]float64{2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.MatrixOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.MatrixOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_RotationOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [4]float64
	}{
		{"default", &Node{Rotation: DefaultRotation}, DefaultRotation},
		{"zeros", &Node{Rotation: emptyRotation}, DefaultRotation},
		{"other", &Node{Rotation: [4]float64{1, 2, 3, 4}}, [4]float64{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.RotationOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.RotationOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_ScaleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [3]float64
	}{
		{"default", &Node{Scale: DefaultScale}, DefaultScale},
		{"zeros", &Node{Scale: emptyScale}, DefaultScale},
		{"other", &Node{Scale: [3]float64{1, 2, 3}}, [3]float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ScaleOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.ScaleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_TranslationOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *Node
		want [3]float64
	}{
		{"default", &Node{Translation: DefaultTranslation}, DefaultTranslation},
		{"other", &Node{Translation: [3]float64{1, 2, 3}}, [3]float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.TranslationOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.TranslationOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_BaseColorFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want RGBA
	}{
		{"empty", &PBRMetallicRoughness{}, *NewRGBA()},
		{"other", &PBRMetallicRoughness{BaseColorFactor: &RGBA{0.8, 0.8, 0.8, 0.5}}, RGBA{0.8, 0.8, 0.8, 0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.BaseColorFactorOrDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PBRMetallicRoughness.BaseColorFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOcclusionTexture_StrengthOrDefault(t *testing.T) {
	tests := []struct {
		name string
		o    *OcclusionTexture
		want float64
	}{
		{"empty", &OcclusionTexture{}, 1},
		{"other", &OcclusionTexture{Strength: Float64(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.StrengthOrDefault(); got != tt.want {
				t.Errorf("OcclusionTexture.StrengthOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalTexture_ScaleOrDefault(t *testing.T) {
	tests := []struct {
		name string
		n    *NormalTexture
		want float64
	}{
		{"empty", &NormalTexture{}, 1},
		{"other", &NormalTexture{Scale: Float64(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.ScaleOrDefault(); got != tt.want {
				t.Errorf("NormalTexture.ScaleOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaterial_AlphaCutoffOrDefault(t *testing.T) {
	tests := []struct {
		name string
		m    *Material
		want float64
	}{
		{"empty", &Material{}, 0.5},
		{"other", &Material{AlphaCutoff: Float64(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.AlphaCutoffOrDefault(); got != tt.want {
				t.Errorf("Material.AlphaCutoffOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_MetallicFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want float64
	}{
		{"empty", &PBRMetallicRoughness{}, 1},
		{"other", &PBRMetallicRoughness{MetallicFactor: Float64(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.MetallicFactorOrDefault(); got != tt.want {
				t.Errorf("PBRMetallicRoughness.MetallicFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBRMetallicRoughness_RoughnessFactorOrDefault(t *testing.T) {
	tests := []struct {
		name string
		p    *PBRMetallicRoughness
		want float64
	}{
		{"empty", &PBRMetallicRoughness{}, 1},
		{"other", &PBRMetallicRoughness{RoughnessFactor: Float64(2)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.RoughnessFactorOrDefault(); got != tt.want {
				t.Errorf("PBRMetallicRoughness.RoughnessFactorOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRGBColor(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		c color.RGBA
	}
	tests := []struct {
		name string
		args args
		want *RGB
	}{
		{"empty", args{color.RGBA{}}, &RGB{}},
		{"base", args{color.RGBA{R: 1, G: 1, B: 1}}, &RGB{R: 0.0003035, G: 0.0003035, B: 0.0003035}},
		{"max", args{color.RGBA{R: 255, G: 255, B: 255}}, &RGB{R: 1, G: 1, B: 1}},
		{"other", args{color.RGBA{R: 60, G: 120, B: 180}}, &RGB{R: 0.045186, G: 0.1878207, B: 0.4564110}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRGBColor(tt.args.c)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewRGBColor() = %v", diff)
			}
		})
	}
}

func TestNewRGB(t *testing.T) {
	tests := []struct {
		name string
		want *RGB
	}{
		{"base", &RGB{1, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRGB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRGB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRGBAColor(t *testing.T) {
	deep.FloatPrecision = 6
	type args struct {
		c color.RGBA
	}
	tests := []struct {
		name string
		args args
		want *RGBA
	}{
		{"empty", args{color.RGBA{}}, &RGBA{}},
		{"base", args{color.RGBA{R: 1, G: 1, B: 1, A: 1}}, &RGBA{R: 0.0003035, G: 0.0003035, B: 0.0003035, A: 0.00392156}},
		{"max", args{color.RGBA{R: 255, G: 255, B: 255, A: 255}}, &RGBA{R: 1, G: 1, B: 1, A: 1}},
		{"other", args{color.RGBA{R: 60, G: 120, B: 180, A: 220}}, &RGBA{R: 0.045186, G: 0.1878207, B: 0.4564110, A: 0.86274509}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRGBAColor(tt.args.c)
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewRGBAColor() = %v", diff)
			}
		})
	}
}
