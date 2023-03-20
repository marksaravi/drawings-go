package colors

import (
	"errors"
)

const (
	RGB565_RED   RGB565 = 0xF800
	RGB565_GREEN RGB565 = 0x07E0
	RGB565_BLUE  RGB565 = 0x001F
)

func RGB888ToRGB565(rgb888 RGB888) RGB565 {

	b := RGB565(rgb888 & 0xFF)
	g := RGB565((rgb888 >> 8) & 0xFF)
	r := RGB565((rgb888 >> 16) & 0xFF)

	return ((r & 0b11111000) << 8) | ((g & 0b11111100) << 3) | (b >> 3)
}

func ToRGB565(color Color) (RGB565, error) {
	c, ok := color.(RGB888)
	if !ok {
		return RGB888ToRGB565(BLACK), errors.New("rgb565 color type mistmatch")
	}
	return RGB888ToRGB565(c), nil
}
