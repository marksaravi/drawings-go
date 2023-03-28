package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/drawings-go/drawings"
	"github.com/marksaravi/drivers-go/colors"
	"github.com/marksaravi/drivers-go/hardware/gpio"
	"github.com/marksaravi/drivers-go/hardware/ili9341"
	"github.com/marksaravi/drivers-go/hardware/spi"

	"github.com/marksaravi/fonts-go/fonts"
	"periph.io/x/host/v3"
)

func ToRad(degree float64) float64 {
	return math.Pi / 180 * degree
}

func ToDeg(rad float64) float64 {
	return rad / math.Pi * 180
}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Testing Sketcher...")
	host.Init()
	spiConn := spi.NewSPI(1, 0, spi.Mode2, 64, 8)
	dataCommandSelect := gpio.NewGPIOOut("GPIO22")
	reset := gpio.NewGPIOOut("GPIO23")

	ili9341Dev, err := ili9341.NewILI9341(ili9341.LCD_320x200, spiConn, dataCommandSelect, reset)
	sketcher := drawings.NewSketcher(ili9341Dev, colors.BLACK)
	checkFatalErr(err)
	tests := []func(drawings.Sketcher){
		drawPoints,
		drawCrossLines,
		drawLines,
		drawArc,
		draThickwArc,
		drawCircle,
		drawFillCircle,
		drawThickCircle,
		drawRectangle,
		drawFillRectangle,
		drawThickRectangle,
		drawFontsArea,
		drawDigits,
		drawCalibrationPoints,
	}

	for i := 0; i < len(tests); i++ {
		sketcher.Clear(colors.WHITE)
		ts := time.Now()
		tests[i](sketcher)
		numsegs := sketcher.Update()
		fmt.Println("Update Duration(ms): ", time.Since(ts).Milliseconds(), ", Num of updated Segments: ", numsegs)
		time.Sleep(time.Second / 2)
	}
	time.Sleep(time.Second)
	fmt.Println("end")
}

func drawPoints(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_0)
	sketcher.Circle(0, 0, 5, colors.RED)
}

func drawCrossLines(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_0)
	xmax := float64(sketcher.ScreenWidth() - 1)
	ymax := float64(sketcher.ScreenHeight() - 1)
	fmt.Println(xmax, ymax)
	sketcher.Clear(colors.BLACK)
	const OFFSET = 5
	sketcher.Line(OFFSET, OFFSET, xmax-OFFSET, OFFSET, colors.RED)
	sketcher.Line(xmax-OFFSET, OFFSET, xmax-OFFSET, ymax-OFFSET, colors.GREEN)
	sketcher.Line(xmax-OFFSET, ymax-OFFSET, OFFSET, ymax-OFFSET, colors.BLUE)
	sketcher.Line(OFFSET, ymax-OFFSET, OFFSET, OFFSET, colors.PINK)
}

func drawLines(sketcher drawings.Sketcher) {
	xmax := float64(sketcher.ScreenWidth() - 1)
	ymax := float64(sketcher.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2
	sAngle := math.Pi / 180 * 0
	rAngle := 2 * math.Pi
	dAngle := math.Pi / 180 * 5

	sketcher.Clear(colors.WHITE)
	for angle := sAngle; angle < sAngle+rAngle; angle += dAngle {
		x := math.Cos(angle) * radius
		y := math.Sin(angle) * radius
		sketcher.Line(xc, yc, xc+x, yc+y, colors.BLUE)
	}
}

func drawCircle(sketcher drawings.Sketcher) {
	const N int = 3
	xmax := float64(sketcher.ScreenWidth() - 1)
	ymax := float64(sketcher.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2.1
	xyc := [N][]float64{{xc, yc, radius}, {xc, yc, radius * .75}, {xc, yc, radius * .45}}
	colorset := [N]colors.Color{colors.BLACK, colors.DARKBLUE, colors.DARKGREEN}
	for i := 0; i < N; i++ {
		sketcher.Circle(xyc[i][0], xyc[i][1], xyc[i][2], colorset[i])
	}
}

func drawFillCircle(sketcher drawings.Sketcher) {
	const N int = 3
	xyc := [N][]float64{{30, 30, 45}, {160, 120, 115}, {400, 400, 250}}
	colorset := [N]colors.Color{colors.BLACK, colors.DARKBLUE, colors.DARKGREEN}
	for i := 0; i < N; i++ {
		sketcher.FillCircle(xyc[i][0], xyc[i][1], xyc[i][2], colorset[i])
	}
}

func drawThickCircle(sketcher drawings.Sketcher) {
	const N int = 3
	xmax := float64(sketcher.ScreenWidth() - 1)
	ymax := float64(sketcher.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2.1
	xyc := [N][]float64{{xc, yc, radius}, {xc, yc, radius * .75}, {xc, yc, radius * .45}}
	colorset := [N]colors.Color{colors.ROYALBLUE, colors.SILVER, colors.MEDIUMSPRINGGREEN}
	widhTypes := [N]drawings.WidthType{drawings.INNER_WIDTH, drawings.CENTER_WIDTH, drawings.OUTER_WIDTH}
	const width = 10
	for i := 0; i < N; i++ {
		sketcher.ThickCircle(xyc[i][0], xyc[i][1], xyc[i][2], width, widhTypes[i], colorset[i])
		sketcher.Circle(xyc[i][0], xyc[i][1], xyc[i][2], colors.RED)
	}
}

func drawArc(sketcher drawings.Sketcher) {
	const N int = 12
	colorset := [N]colors.Color{
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
	}

	xyc := [N][]float64{
		{160, 120, 50, ToRad(0), ToRad(90)},
		{160, 120, 55, ToRad(90), ToRad(180)},
		{160, 120, 60, ToRad(180), ToRad(270)},
		{160, 120, 65, ToRad(270), ToRad(360)},
		{160, 120, 70, ToRad(15), ToRad(45)},
		{160, 120, 75, ToRad(45), ToRad(15)},
		{160, 120, 80, ToRad(105), ToRad(135)},
		{160, 120, 85, ToRad(135), ToRad(105)},
		{160, 120, 90, ToRad(195), ToRad(225)},
		{160, 120, 95, ToRad(225), ToRad(195)},
		{160, 120, 100, ToRad(285), ToRad(315)},
		{160, 120, 105, ToRad(315), ToRad(285)},
	}
	for i := 0; i < N; i++ {
		sketcher.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4], colorset[i])
	}
	sketcher.Line(160, 0, 160, 239, colors.RED)
	sketcher.Line(0, 120, 319, 120, colors.RED)
}

func draThickwArc(sketcher drawings.Sketcher) {
	const N int = 3
	colorset := [N]colors.Color{
		colors.CYAN,
		colors.GREEN,
		colors.LIGHTBLUE,
	}

	widhTypes := [N]drawings.WidthType{drawings.OUTER_WIDTH, drawings.CENTER_WIDTH, drawings.INNER_WIDTH}
	xyc := [N][]float64{
		{160, 120, 70, ToRad(45), ToRad(175)},
		{160, 120, 90, ToRad(15), ToRad(300)},
		{160, 120, 115, ToRad(300), ToRad(15)},
	}

	for i := 0; i < N; i++ {
		sketcher.ThickArc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4], 10, widhTypes[i], colorset[i])
		sketcher.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4], colors.RED)
	}
}

func drawRectangle(sketcher drawings.Sketcher) {
	const N int = 2
	xy := [N][]float64{{10, 10, 100, 100}, {50, 50, 200, 200}}
	colorset := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		sketcher.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], colorset[i])
	}

}

func drawFillRectangle(sketcher drawings.Sketcher) {
	const N int = 2
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}}
	colors := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		sketcher.FillRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], colors[i])
	}

}

func drawThickRectangle(sketcher drawings.Sketcher) {
	const N int = 3
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}, {100, 100, 300, 220}}
	colorset := [N]colors.Color{colors.ROYALBLUE, colors.NAVY, colors.FORESTGREEN}
	widhTypes := [N]drawings.WidthType{drawings.INNER_WIDTH, drawings.CENTER_WIDTH, drawings.OUTER_WIDTH}
	const width = 10
	for i := 0; i < N; i++ {
		sketcher.ThickRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], width, widhTypes[i], colorset[i])
		sketcher.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], colors.RED)
	}

}

func drawFontsArea(sketcher drawings.Sketcher) {
	sketcher.SetFont(fonts.FreeSerif18pt7b)
	const LEN = 12
	const FROM byte = 0x20 + 20
	const TO byte = 0x7E
	var c byte = FROM
	yline := 32

	for c <= TO {
		s := make([]byte, 0)
		for i := 0; i < LEN && c <= TO; i++ {
			s = append(s, c)
			c++
		}
		text := string(s)
		xoffset := 8
		x1, y1, x2, y2 := sketcher.GetTextArea(float64(xoffset), float64(yline), text, 1, 1)

		sketcher.Rectangle(float64(x1), float64(y1), float64(x2), float64(y2), colors.RED)
		sketcher.Line(0, float64(yline), 319, float64(yline), colors.BLUE)
		sketcher.MoveCursor(int(xoffset), yline)
		sketcher.Write(string(s), colors.BLACK)
		yline += 48
	}
}

func drawGrids(sketcher drawings.Sketcher) {
	for x := float64(0); x < 320; x += 32 {
		sketcher.Line(x, 0, x, 239, colors.RED)
	}
	for y := float64(0); y < 240; y += 24 {
		sketcher.Line(0, y, 319, y, colors.RED)
	}
}

func drawDigits(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_270)
	sketcher.SetFont(fonts.FreeSans24pt7b)
	X := 30
	Y := 120
	xScale := float64(1)
	yScale := float64(1)
	value := 23.2
	text := fmt.Sprintf("%4.1f", value)
	x1, y1, x2, y2 := sketcher.GetTextArea(float64(X), float64(Y), text, xScale, yScale)
	fmt.Println(x1, y1, x2, y2)
	sketcher.Rectangle(float64(x1), float64(y1), float64(x2), float64(y2), colors.RED)
	sketcher.MoveCursor(X, Y)
	sketcher.Write(text, colors.BLACK)
}

func drawCalibrationPoints(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_180)
	const PADDING float64 = 40
	const RADIUS float64 = 5
	const N_SEGMENTS int = 2
	var X_OFFSET float64 = (float64(sketcher.ScreenWidth()) - PADDING*2) / float64(N_SEGMENTS)
	var Y_OFFSET float64 = (float64(sketcher.ScreenHeight()) - PADDING*2) / float64(N_SEGMENTS)

	for xseg := 0; xseg <= N_SEGMENTS; xseg++ {
		x := float64(xseg) * X_OFFSET
		for yseg := 0; yseg <= N_SEGMENTS; yseg++ {
			y := float64(yseg) * Y_OFFSET
			sketcher.FillCircle(x+PADDING, y+PADDING, RADIUS, colors.RED)
		}
	}

	for i := float64(0); i <= 0; i++ {
		sketcher.Line(0+i, 0+i, 319-i, 0+i, colors.RED)
		sketcher.Line(319-i, 0+i, 319-i, 239-i, colors.RED)
		sketcher.Line(319-i, 239-i, 0+i, 239-i, colors.RED)
		sketcher.Line(0+i, 239-i, 0+i, 0+i, colors.RED)
	}
}
