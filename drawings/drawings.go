package drawings

import (
	"errors"
	"math"

	"github.com/marksaravi/drawings-go/colors"
	"github.com/marksaravi/fonts-go/fonts"
)

type WidthType int
type FontType int

const (
	DEG90  = math.Pi / 2
	DEG180 = DEG90 * 2
	DEG270 = DEG90 * 3
	DEG360 = DEG90 * 4
)

const (
	INNER_WIDTH    WidthType = 0
	OUTER_WIDTH    WidthType = 1
	CENTER_WIDTH   WidthType = 2
	MAX_FONT_SCALE int       = 10
)

const (
	BITMAP_FONT FontType = 0

	ROTATION_0   = 0
	ROTATION_90  = 1
	ROTATION_180 = 2
	ROTATION_270 = 3
)

type arcSector struct {
	ok             bool
	xs, xe, ys, ye float64
}

type pixelDevice interface {
	Update() int
	Pixel(x, y int, color colors.Color)
	Clear(color colors.Color)
	ScreenWidth() int
	ScreenHeight() int
}

type Sketcher interface {
	Update() int
	SetRotation(rotation int)
	ScreenWidth() int
	ScreenHeight() int
	SetBackgroundColor(color colors.Color)
	SetColor(color colors.Color)
	ClearArea(x1, y1, x2, y2 float64)
	Clear()
	// Pixel(x, y float64)
	Line(x1, y1, x2, y2 float64)
	Arc(xc, yc, radius, startAngle, endAngle float64)
	ThickArc(xc, yc, radius, startAngle, endAngle float64, width int, widthType WidthType)
	Circle(x, y, radius float64)
	Rectangle(x1, y1, x2, y2 float64)
	FillCircle(x, y, radius float64)
	ThickCircle(x, y, radius float64, width int, widthType WidthType)
	FillRectangle(x1, y1, x2, y2 float64)
	ThickRectangle(x1, y1, x2, y2 float64, width int, widthType WidthType)
	SetFont(font interface{}) error
	WriteScaled(text string, xscale, yscale int)
	Write(text string)
	MoveCursor(x, y int)
	GetTextArea(text string) (x1, y1, x2, y2 int)
}

type sketcher struct {
	pixeldev        pixelDevice
	color           colors.Color
	bgColor         colors.Color
	font            interface{}
	bitmapFont      fonts.BitmapFont
	fontType        FontType
	cursorX         int
	cursorY         int
	charAdvanceX    int
	textLeftPadding int
	textTopPadding  int
	rotation        int
}

func NewSketcher(pixeldev pixelDevice) Sketcher {
	return &sketcher{
		pixeldev:        pixeldev,
		color:           colors.WHITE,
		bgColor:         colors.BLACK,
		fontType:        BITMAP_FONT,
		font:            fonts.FreeMono18pt7b,
		cursorX:         0,
		cursorY:         0,
		charAdvanceX:    0,
		textLeftPadding: 0,
		textTopPadding:  0,
		rotation:        ROTATION_0,
	}
}

func (d *sketcher) Update() int {
	return d.pixeldev.Update()
}

func (d *sketcher) SetRotation(rotation int) {
	d.rotation = rotation
}

func (d *sketcher) ScreenWidth() int {
	if d.rotation == ROTATION_90 || d.rotation == ROTATION_270 {
		return d.pixeldev.ScreenHeight()
	}
	return d.pixeldev.ScreenWidth()
}

func (d *sketcher) ScreenHeight() int {
	if d.rotation == ROTATION_90 || d.rotation == ROTATION_270 {
		return d.pixeldev.ScreenWidth()
	}
	return d.pixeldev.ScreenHeight()
}

func (d *sketcher) SetBackgroundColor(color colors.Color) {
	d.bgColor = color
}

func (d *sketcher) SetColor(color colors.Color) {
	d.color = color
}

func (d *sketcher) ClearArea(x1, y1, x2, y2 float64) {
	xs := int(math.Round(x1))
	xe := int(math.Round(x2))
	ys := int(math.Round(y1))
	ye := int(math.Round(y2))
	if x1 > x2 {
		t := xs
		xs = xe
		xe = t
	}
	if y1 > y2 {
		t := ys
		ys = ye
		ye = t
	}
	for x := xs; x <= xe; x += 1 {
		for y := ys; y <= ye; y += 1 {
			d.rotatedPixel(float64(x), float64(y), d.bgColor)
		}
	}
}

// Drawing methods
func (d *sketcher) Clear() {
	d.pixeldev.Clear(d.bgColor)
}

func (d *sketcher) rotatePoint(x, y float64) (float64, float64) {
	if d.rotation == ROTATION_0 {
		return x, y
	}
	if d.rotation == ROTATION_90 {
		return float64(d.pixeldev.ScreenWidth()) - y, x
	}
	if d.rotation == ROTATION_180 {
		return float64(d.pixeldev.ScreenWidth()) - x, float64(d.pixeldev.ScreenHeight()) - y
	}
	return y, float64(d.pixeldev.ScreenHeight()) - x
}

func (d *sketcher) rotatedPixel(x, y float64, color colors.Color) {
	rotatedX, rotatedY := d.rotatePoint(x, y)
	d.pixeldev.Pixel(int(math.Round(rotatedX)), int(math.Round(rotatedY)), color)
}

func (d *sketcher) Line(x1, y1, x2, y2 float64) {
	// Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	xs := int(math.Round(x1))
	ys := int(math.Round(y1))
	xe := int(math.Round(x2))
	ye := int(math.Round(y2))
	dx := int(math.Abs(x2 - x1))
	// sx := xs < xe ? 1 : -1
	sx := -1
	if xs < xe {
		sx = 1
	}
	dy := -int(math.Abs(y2 - y1))
	// sy := ys < ye ? 1 : -1
	sy := -1
	if ys < ye {
		sy = 1
	}
	err := dx + dy

	for {
		d.rotatedPixel(float64(xs), float64(ys), d.color)
		if xs == xe && ys == ye {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			if xs == xe {
				break
			}
			err = err + dy
			xs = xs + sx
		}
		if e2 <= dx {
			if ys == ye {
				break
			}
			err = err + dx
			ys = ys + sy
		}
	}
}

func isInsideSector0(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x > xe && y >= ys && y < ye
}

func isInsideSector1(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x > xe && y <= ys && y > ye
}

func isInsideSector2(x, y, xs, ys, xe, ye float64) bool {
	return x >= xs && x < xe && y <= ys && y > ye
}

func isInsideSector3(x, y, xs, ys, xe, ye float64) bool {
	return x >= xs && x < xe && y >= ys && y < ye
}

func findArcSectors(startAngle, endAngle, radius float64) map[int][]arcSector {
	var sectorsmap map[int][]arcSector = map[int][]arcSector{
		0: make([]arcSector, 0),
		1: make([]arcSector, 0),
		2: make([]arcSector, 0),
		3: make([]arcSector, 0),
	}
	PI2 := math.Pi / 2
	from := math.Mod(startAngle, math.Pi*2)
	to := math.Mod(endAngle, math.Pi*2)
	if to < from {
		to += math.Pi * 2
	}
	angle := from
	for sec := 0; angle < to; sec++ {
		sector := arcSector{
			ok: false,
			xs: 0,
			ys: 0,
			xe: 0,
			ye: 0,
		}
		s := float64(sec) * PI2
		e := float64(sec+1) * PI2
		if e >= to {
			e = to
		}
		if angle >= s && angle < e {
			sector.ok = true
			sector.xs = radius * math.Cos(angle)
			sector.ys = radius * math.Sin(angle)
			sector.xe = radius * math.Cos(e)
			sector.ye = radius * math.Sin(e)
			angle = e
			sectorsmap[sec%4] = append(sectorsmap[sec%4], sector)
		}
	}
	// showSectors(sectorsmap)
	return sectorsmap
}

func (dev *sketcher) arcPutPixel(sector int, xc, yc, x, y float64, s arcSector) {
	tests := []func(x, y, xs, ys, xe, ye float64) bool{
		isInsideSector0,
		isInsideSector1,
		isInsideSector2,
		isInsideSector3,
	}
	if tests[sector](x, y, s.xs, s.ys, s.xe, s.ye) {
		dev.rotatedPixel(x+xc, y+yc, dev.color)
	}
}

func (dev *sketcher) Arc(xc, yc, radius, startAngle, endAngle float64) {
	signs := [4][2]float64{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}
	iradius := math.Round(radius)
	sectormaps := findArcSectors(startAngle, endAngle, iradius)
	var iradius2 = iradius * iradius
	var l1 float64 = 0
	for l1 = 0; true; l1 += 1 {
		l2 := math.Sqrt(iradius2 - l1*l1)
		for sector := 0; sector < 4; sector++ {
			sectors := sectormaps[sector]
			for i := 0; i < len(sectors); i++ {
				if sectors[i].ok {
					dev.arcPutPixel(sector, xc, yc, signs[sector][0]*l1, signs[sector][1]*l2, sectors[i])
					dev.arcPutPixel(sector, xc, yc, signs[sector][0]*l2, signs[sector][1]*l1, sectors[i])
				}
			}
		}

		if l1 >= l2 {
			break
		}
	}
}

func (dev *sketcher) ThickArc(xc, yc, radius, startAngle, endAngle float64, width int, widthType WidthType) {
	rs := calcThicknessStart(radius, width, widthType)
	for dr := 0; dr < width; dr++ {
		dev.Arc(xc, yc, rs-float64(dr), startAngle, endAngle)
	}
}

func (dev *sketcher) Circle(x, y, radius float64) {
	// Midpoint circle algorithm https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	putpixels := func(xc, yc, dr, d float64) {
		dev.rotatedPixel(xc+d, yc+dr, dev.color)
		dev.rotatedPixel(xc+d, yc-dr, dev.color)
		dev.rotatedPixel(xc+dr, yc+d, dev.color)
		dev.rotatedPixel(xc+dr, yc-d, dev.color)

		dev.rotatedPixel(xc-d, yc+dr, dev.color)
		dev.rotatedPixel(xc-d, yc-dr, dev.color)
		dev.rotatedPixel(xc-dr, yc+d, dev.color)
		dev.rotatedPixel(xc-dr, yc-d, dev.color)
	}

	var dy float64 = radius
	for dx := float64(0); dx < dy; dx += 1 {
		dy = math.Sqrt(radius*radius - dx*dx)
		putpixels(x, y, dx, dy)
	}
}

func (dev *sketcher) FillCircle(x, y, radius float64) {
	// Midpoint circle algorithm https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	putpixels := func(xc, yc, dr, d float64) {
		dev.Line(xc+d, yc+dr, xc-d, yc+dr)
		dev.Line(xc+d, yc-dr, xc-d, yc-dr)

		dev.Line(xc+dr, yc+d, xc-dr, yc+d)
		dev.Line(xc+dr, yc-d, xc-dr, yc-d)
	}
	for dr := float64(0); dr <= math.Ceil(radius*0.707); dr += 1 {
		d := math.Sqrt(radius*radius - dr*dr)
		putpixels(x, y, dr, d)
	}
}

func calcThicknessStart(mid float64, width int, widthType WidthType) float64 {
	from := mid
	switch widthType {
	case OUTER_WIDTH:
		from = mid + float64(width)
	case CENTER_WIDTH:
		from = mid + float64(width)/2
	}
	return from
}

func (dev *sketcher) ThickCircle(x, y, radius float64, width int, widthType WidthType) {
	rs := calcThicknessStart(radius, width, widthType)
	for dr := 0; dr < width; dr++ {
		dev.Circle(x, y, rs-float64(dr))
	}
}

func (dev *sketcher) Rectangle(x1, y1, x2, y2 float64) {
	dev.Line(x1, y1, x2, y1)
	dev.Line(x2, y1, x2, y2)
	dev.Line(x2, y2, x1, y2)
	dev.Line(x1, y2, x1, y1)
}

func (dev *sketcher) FillRectangle(x1, y1, x2, y2 float64) {
	l := math.Round(y2 - y1)
	dy := float64(1)
	if l < 0 {
		dy = -1
	}

	for y := float64(0); y != l; y += dy {
		dev.Line(x1, y1+y, x2, y1+y)
	}
}

func (dev *sketcher) ThickRectangle(x1, y1, x2, y2 float64, width int, widthType WidthType) {
	xs := x1
	xe := x2
	if x2 < x1 {
		xs = x2
		xe = x1
	}
	ys := y1
	ye := y2
	if y2 < y1 {
		ys = y2
		ye = y1
	}
	s := calcThicknessStart(0, width, widthType)
	for dxy := float64(0); dxy < float64(width); dxy++ {
		dev.Rectangle(xs-s+dxy, ys-s+dxy, xe+s-dxy, ye+s-dxy)
	}
}

func (dev *sketcher) SetFont(font interface{}) error {
	dev.font = font
	if bitmapfont, ok := font.(fonts.BitmapFont); ok {
		dev.fontType = BITMAP_FONT
		dev.bitmapFont = bitmapfont
		return nil
	}
	return errors.New("font format is not implemented")
}

func (dev *sketcher) writeChar(char byte, xscale, yscale int) error {
	if char < ' ' || char > '~' {
		return errors.New("charCode code out of range")
	}

	switch dev.fontType {
	case BITMAP_FONT:
		dev.drawBitmapChar(char, xscale, yscale)
	default:
		return errors.New("font is not defined")
	}
	return nil
}

func (dev *sketcher) WriteScaled(text string, xscale, yscale int) {
	if xscale < 1 {
		xscale = 1
	}
	if yscale < 1 {
		yscale = 1
	}
	if xscale > MAX_FONT_SCALE {
		xscale = MAX_FONT_SCALE
	}
	if yscale > MAX_FONT_SCALE {
		yscale = MAX_FONT_SCALE
	}
	for i := 0; i < len(text); i++ {
		dev.writeChar(text[i], xscale, yscale)
	}
}

func (dev *sketcher) Write(text string) {
	for i := 0; i < len(text); i++ {
		dev.writeChar(text[i], 1, 1)
	}
}

func (dev *sketcher) MoveCursor(x, y int) {
	dev.cursorX = x
	dev.cursorY = y
}

func (dev *sketcher) GetTextArea(text string) (x1, y1, x2, y2 int) {
	x1 = 0
	y1 = 0
	x2 = 0
	y2 = 0
	switch dev.fontType {
	case BITMAP_FONT:
		x1, y1, x2, y2 = dev.getBitmapFontTextArea(text)
	}
	return
}

func (dev *sketcher) drawBitmapChar(char byte, xscale, yscale int) {
	glyph := dev.bitmapFont.Glyphs[char-0x20]
	for h := 0; h < glyph.Height; h++ {
		for w := 0; w < glyph.Width; w++ {
			bitIndex := h*glyph.Width + w
			shift := byte(bitIndex) % 8
			d := dev.bitmapFont.Bitmap[glyph.BitmapOffset+bitIndex/8]
			mask := byte(0b10000000) >> shift
			bit := d & mask
			color := dev.bgColor
			if bit != 0 {
				color = dev.color
			}
			x := dev.cursorX + (w+glyph.XOffset)*xscale
			y := dev.cursorY + (h+glyph.YOffset)*yscale
			for dx := 0; dx < xscale; dx++ {
				for dy := 0; dy < yscale; dy++ {
					dev.rotatedPixel(float64(x+dx), float64(y+dy), color)
				}
			}
		}
	}
	dev.cursorX += glyph.XAdvance * xscale
}

func (dev *sketcher) getBitmapFontTextArea(text string) (int, int, int, int) {
	bytes := []byte(text)
	ymax := 0
	ymin := 0
	x := 0
	for i := 0; i < len(bytes); i++ {
		glyph := dev.bitmapFont.Glyphs[bytes[i]-0x20]
		x += glyph.XAdvance
		y := glyph.YOffset + glyph.Height
		if y > ymax {
			ymax = y
		}
		if glyph.YOffset < ymin {
			ymin = glyph.YOffset
		}
	}
	return 0, ymin, x, ymax
}
