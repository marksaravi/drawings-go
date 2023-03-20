package colors

import (
	"testing"
)

func TestConvertRBG888ToRGB565(t *testing.T) {
	colorset := []Color{RED, GREEN, BLUE, BLACK, CYAN, YELLOW, NAVY, ROYALBLUE, BROWN}
	want := []RGB565{0xF800, 0x07E0, 0x001F, 0x0000, 0x07FF, 0xFFE0, 0x0010, 0x435C, 0xA145}

	for i := 0; i < len(colorset); i++ {
		got, _ := ToRGB565(colorset[i])
		if got != want[i] {
			t.Errorf("at %d, wanted %x, got %x", i, want[i], got)
		}
	}
}
