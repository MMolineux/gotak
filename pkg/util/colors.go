package util

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// StandardColors maps common color names to their hex RGB values
var StandardColors = map[string]uint32{
	"black":   0xFF000000,
	"white":   0xFFFFFFFF,
	"red":     0xFFFF0000,
	"green":   0xFF00FF00,
	"blue":    0xFF0000FF,
	"yellow":  0xFFFFFF00,
	"cyan":    0xFF00FFFF,
	"magenta": 0xFFFF00FF,
	"gray":    0xFF808080,
	"grey":    0xFF808080,
	"orange":  0xFFFFA500,
	"purple":  0xFF800080,
	"brown":   0xFFA52A2A,
	"pink":    0xFFFFC0CB,
	"lime":    0xFF00FF00,
	"teal":    0xFF008080,
	"navy":    0xFF000080,
	"maroon":  0xFF800000,
	"olive":   0xFF808000,
}

// ColorConverter provides methods for converting colors between different formats
type ColorConverter struct{}

// NewColorConverter creates a new ColorConverter
func NewColorConverter() *ColorConverter {
	return &ColorConverter{}
}

// HexToInt converts a hex color string to uint32 integer format
// Accepts formats: "#RRGGBB", "#AARRGGBB", "RRGGBB", "AARRGGBB"
func (cc *ColorConverter) HexToInt(hexColor string) (uint32, error) {
	// Remove # prefix if present
	hexColor = strings.TrimPrefix(hexColor, "#")

	// Handle both RGB and ARGB formats
	var fullColor string
	switch len(hexColor) {
	case 6: // RRGGBB format
		fullColor = "FF" + hexColor // Add full opacity
	case 8: // AARRGGBB format
		fullColor = hexColor
	default:
		return 0, fmt.Errorf("invalid hex color format: %s", hexColor)
	}

	// Convert hex to uint32
	colorInt, err := strconv.ParseUint(fullColor, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid hex value: %v", err)
	}

	return uint32(colorInt), nil
}

// IntToHex converts a uint32 color value to hex string
// Returns format: AARRGGBB
func (cc *ColorConverter) IntToHex(colorInt uint32) string {
	return fmt.Sprintf("%08X", colorInt)
}

// RGBToInt converts RGB values to uint32 format with full opacity
func (cc *ColorConverter) RGBToInt(r, g, b uint8) uint32 {
	return 0xFF000000 | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b)
}

// RGBAToInt converts RGBA values to uint32 format
func (cc *ColorConverter) RGBAToInt(r, g, b, a uint8) uint32 {
	return (uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b)
}

// IntToRGBA converts a uint32 color to individual RGBA components
func (cc *ColorConverter) IntToRGBA(colorInt uint32) (r, g, b, a uint8) {
	r = uint8((colorInt >> 16) & 0xFF)
	g = uint8((colorInt >> 8) & 0xFF)
	b = uint8(colorInt & 0xFF)
	a = uint8((colorInt >> 24) & 0xFF)
	return
}

// NameToInt converts a color name to uint32 format
// If the color name isn't found, returns an error
func (cc *ColorConverter) NameToInt(colorName string) (uint32, error) {
	colorName = strings.ToLower(strings.TrimSpace(colorName))

	if val, ok := StandardColors[colorName]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("unknown color name: %s", colorName)
}

// ParseColor attempts to parse a color from various formats:
// - Named colors: "red", "blue", etc.
// - Hex format: "#FF0000", "FF0000", "#FFFF0000", "FFFF0000"
// - RGB format: "rgb(255,0,0)", "rgb(255, 0, 0)"
// - RGBA format: "rgba(255,0,0,255)", "rgba(255, 0, 0, 255)"
func (cc *ColorConverter) ParseColor(colorStr string) (uint32, error) {
	colorStr = strings.TrimSpace(colorStr)

	// Empty string
	if colorStr == "" {
		return 0, fmt.Errorf("empty color string")
	}

	// Try as a named color
	if val, err := cc.NameToInt(colorStr); err == nil {
		return val, nil
	}

	// Try as hex
	if strings.HasPrefix(colorStr, "#") || isHexString(colorStr) {
		return cc.HexToInt(colorStr)
	}

	// Try as RGB format
	if strings.HasPrefix(strings.ToLower(colorStr), "rgb(") && strings.HasSuffix(colorStr, ")") {
		// Extract the rgb values
		rgbPart := colorStr[4 : len(colorStr)-1]
		parts := strings.Split(rgbPart, ",")

		if len(parts) == 3 {
			r, errR := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, errG := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, errB := strconv.Atoi(strings.TrimSpace(parts[2]))

			if errR == nil && errG == nil && errB == nil {
				return cc.RGBToInt(uint8(r), uint8(g), uint8(b)), nil
			}
		}
	}

	// Try as RGBA format
	if strings.HasPrefix(strings.ToLower(colorStr), "rgba(") && strings.HasSuffix(colorStr, ")") {
		// Extract the rgba values
		rgbaPart := colorStr[5 : len(colorStr)-1]
		parts := strings.Split(rgbaPart, ",")

		if len(parts) == 4 {
			r, errR := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, errG := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, errB := strconv.Atoi(strings.TrimSpace(parts[2]))
			a, errA := strconv.Atoi(strings.TrimSpace(parts[3]))

			if errR == nil && errG == nil && errB == nil && errA == nil {
				return cc.RGBAToInt(uint8(r), uint8(g), uint8(b), uint8(a)), nil
			}
		}
	}

	// If we get here, we couldn't parse the color
	return 0, fmt.Errorf("unable to parse color format: %s", colorStr)
}

// UintToInt converts an unsigned 32-bit integer color to a signed 32-bit integer
// This is needed because CoT expects colors as signed integers
func (cc *ColorConverter) UintToInt(uintColor uint32) int32 {
	return int32(uintColor)
}

// IntToUint converts a signed 32-bit integer color to an unsigned 32-bit integer
func (cc *ColorConverter) IntToUint(intColor int32) uint32 {
	return uint32(intColor)
}

// ConvertToCoTColor converts a color to the format expected by CoT
// Returns the color as a string representation of an integer in CoT format
func (cc *ColorConverter) ConvertToCoTColor(colorVal interface{}) (string, error) {
	var colorInt uint32
	var err error

	switch v := colorVal.(type) {
	case string:
		colorInt, err = cc.ParseColor(v)
		if err != nil {
			return "", err
		}
	case uint32:
		colorInt = v
	case int:
		colorInt = uint32(v)
	case int32:
		colorInt = uint32(v)
	case color.RGBA:
		colorInt = cc.RGBAToInt(v.R, v.G, v.B, v.A)
	case color.NRGBA:
		colorInt = cc.RGBAToInt(v.R, v.G, v.B, v.A)
	default:
		return "", fmt.Errorf("unsupported color type: %T", colorVal)
	}

	// Convert to the signed integer format expected by CoT
	cotColorInt := cc.UintToInt(colorInt)

	// Return as a string representation of the integer
	return fmt.Sprintf("%d", cotColorInt), nil
}

// GetSignedInt returns the color as a signed integer (the format used by CoT)
func (cc *ColorConverter) GetSignedInt(colorVal interface{}) (int32, error) {
	var colorInt uint32
	var err error

	switch v := colorVal.(type) {
	case string:
		colorInt, err = cc.ParseColor(v)
		if err != nil {
			return 0, err
		}
	case uint32:
		colorInt = v
	case int:
		colorInt = uint32(v)
	case int32:
		return v, nil
	case color.RGBA:
		colorInt = cc.RGBAToInt(v.R, v.G, v.B, v.A)
	case color.NRGBA:
		colorInt = cc.RGBAToInt(v.R, v.G, v.B, v.A)
	default:
		return 0, fmt.Errorf("unsupported color type: %T", colorVal)
	}

	return cc.UintToInt(colorInt), nil
}

// ParseCoTColor parses a CoT format color (signed integer string) to uint32
func (cc *ColorConverter) ParseCoTColor(cotColor string) (uint32, error) {
	// Try to parse the color as a signed integer
	intVal, err := strconv.ParseInt(cotColor, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid CoT color format: %v", err)
	}

	// Convert the signed int32 to uint32
	return cc.IntToUint(int32(intVal)), nil
}

// Helper function to check if a string is a valid hex color
func isHexString(s string) bool {
	if len(s) != 6 && len(s) != 8 {
		return false
	}

	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}

	return true
}
