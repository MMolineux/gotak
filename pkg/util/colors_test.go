package util

import (
	"image/color"
	"testing"
)

func TestHexToInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		hexColor string
		want     uint32
		wantErr  bool
	}{
		{"RGB format", "#FF0000", 0xFFFF0000, false},
		{"RGB without hash", "00FF00", 0xFF00FF00, false},
		{"ARGB format", "#7700FFFF", 0x7700FFFF, false},
		{"ARGB without hash", "880000FF", 0x880000FF, false},
		{"Invalid length", "#12345", 0, true},
		{"Invalid chars", "#GGGGGG", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.HexToInt(tt.hexColor)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToHex(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		colorInt uint32
		want     string
	}{
		{"Red", 0xFFFF0000, "FFFF0000"},
		{"Green", 0xFF00FF00, "FF00FF00"},
		{"Transparent blue", 0x800000FF, "800000FF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.IntToHex(tt.colorInt); got != tt.want {
				t.Errorf("IntToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGBToInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name string
		r    uint8
		g    uint8
		b    uint8
		want uint32
	}{
		{"Red", 255, 0, 0, 0xFFFF0000},
		{"Green", 0, 255, 0, 0xFF00FF00},
		{"Blue", 0, 0, 255, 0xFF0000FF},
		{"White", 255, 255, 255, 0xFFFFFFFF},
		{"Black", 0, 0, 0, 0xFF000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.RGBToInt(tt.r, tt.g, tt.b); got != tt.want {
				t.Errorf("RGBToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGBAToInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name string
		r    uint8
		g    uint8
		b    uint8
		a    uint8
		want uint32
	}{
		{"Opaque red", 255, 0, 0, 255, 0xFFFF0000},
		{"Semi-transparent green", 0, 255, 0, 128, 0x8000FF00},
		{"Transparent blue", 0, 0, 255, 0, 0x000000FF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.RGBAToInt(tt.r, tt.g, tt.b, tt.a); got != tt.want {
				t.Errorf("RGBAToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToRGBA(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		colorInt uint32
		wantR    uint8
		wantG    uint8
		wantB    uint8
		wantA    uint8
	}{
		{"Red", 0xFFFF0000, 255, 0, 0, 255},
		{"Green", 0x8000FF00, 0, 255, 0, 128},
		{"Blue", 0x000000FF, 0, 0, 255, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotG, gotB, gotA := cc.IntToRGBA(tt.colorInt)
			if gotR != tt.wantR || gotG != tt.wantG || gotB != tt.wantB || gotA != tt.wantA {
				t.Errorf("IntToRGBA() = (%v, %v, %v, %v), want (%v, %v, %v, %v)",
					gotR, gotG, gotB, gotA, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
			}
		})
	}
}

func TestNameToInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name      string
		colorName string
		want      uint32
		wantErr   bool
	}{
		{"Red", "red", 0xFFFF0000, false},
		{"RED uppercase", "RED", 0xFFFF0000, false},
		{"Blue with spaces", " blue ", 0xFF0000FF, false},
		{"Unknown color", "notacolor", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.NameToInt(tt.colorName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NameToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NameToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		colorStr string
		want     uint32
		wantErr  bool
	}{
		{"Named color", "red", 0xFFFF0000, false},
		{"Hex with hash", "#00FF00", 0xFF00FF00, false},
		{"Hex without hash", "0000FF", 0xFF0000FF, false},
		{"ARGB hex", "#80FF00FF", 0x80FF00FF, false},
		{"RGB format", "rgb(255,0,0)", 0xFFFF0000, false},
		{"RGB with spaces", "rgb(0, 255, 0)", 0xFF00FF00, false},
		{"RGBA format", "rgba(0,0,255,128)", 0x800000FF, false},
		{"Invalid format", "invalid", 0, true},
		{"Empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.ParseColor(tt.colorStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseColor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToCoTColor(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		colorVal interface{}
		want     string
		wantErr  bool
	}{
		{"String color name", "red", "FFFF0000", false},
		{"String hex", "#00FF00", "FF00FF00", false},
		{"Uint32 value", uint32(0xFF0000FF), "FF0000FF", false},
		{"Int value", int(0xFFFF00FF), "FFFF00FF", false},
		{"color.RGBA", color.RGBA{R: 255, G: 0, B: 0, A: 255}, "FFFF0000", false},
		{"Invalid type", 3.14, "", true},
		{"Invalid string", "notacolor", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.ConvertToCoTColor(tt.colorVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToCoTColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToCoTColor() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUintToInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name      string
		uintColor uint32
		want      int32
	}{
		{"Positive value", 0x7FFFFFFF, 2147483647},
		{"Negative value", 0xFFFFFFFF, -1},
		{"Standard red", 0xFFFF0000, -16777216},
		{"Blue with alpha", 0x800000FF, -2147483393},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.UintToInt(tt.uintColor); got != tt.want {
				t.Errorf("UintToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToUint(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		intColor int32
		want     uint32
	}{
		{"Positive value", 2147483647, 0x7FFFFFFF},
		{"Negative value", -1, 0xFFFFFFFF},
		{"Negative red", -16777216, 0xFFFF0000},
		{"Negative alpha blue", -2147483393, 0x800000FF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.IntToUint(tt.intColor); got != tt.want {
				t.Errorf("IntToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSignedInt(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		colorVal interface{}
		want     int32
		wantErr  bool
	}{
		{"String color name", "red", -16777216, false},         // 0xFFFF0000 as signed int
		{"String hex", "#00FF00", -16711936, false},            // 0xFF00FF00 as signed int
		{"Uint32 value", uint32(0xFF0000FF), -16776961, false}, // 0xFF0000FF as signed int
		{"Int32 value", int32(-16777216), -16777216, false},    // Already CoT format
		{"color.RGBA", color.RGBA{R: 255, G: 0, B: 0, A: 255}, -16777216, false},
		{"Invalid type", 3.14, 0, true},
		{"Invalid string", "notacolor", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.GetSignedInt(tt.colorVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSignedInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSignedInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCoTColor(t *testing.T) {
	cc := NewColorConverter()

	tests := []struct {
		name     string
		cotColor string
		want     uint32
		wantErr  bool
	}{
		{"Red as CoT int", "-16777216", 0xFFFF0000, false},
		{"Green as CoT int", "-16711936", 0xFF00FF00, false},
		{"Blue as CoT int", "-16776961", 0xFF0000FF, false},
		{"Alpha gray as CoT int", "-1761607681", 0x97AAAAAA, false}, // Example from prompt
		{"Invalid format", "not-a-number", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cc.ParseCoTColor(tt.cotColor)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCoTColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCoTColor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
