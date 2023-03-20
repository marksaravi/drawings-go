package main

import (
	"fmt"
	"log"
	"math"
	"time"

	devgpio "github.com/marksaravi/devices-go/hardware/gpio"
	"github.com/marksaravi/devices-go/hardware/ili9341"
	"github.com/marksaravi/drawings-go/colors"
	"github.com/marksaravi/drawings-go/drawings"
	"github.com/marksaravi/fonts-go/fonts"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

func ToRad(degree float64) float64 {
	return math.Pi / 180 * degree
}

func ToDeg(rad float64) float64 {
	return rad / math.Pi * 180
}

type gpioOut struct {
	pin gpio.PinOut
}

func (p *gpioOut) Out(level devgpio.Level) {
	p.pin.Out(gpio.Level(level))
}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Testing Sketcher...")
	host.Init()
	spiConn := createSPIConnection(1, 0)
	dataCommandSelect := createGpioOutPin("GPIO22")
	reset := createGpioOutPin("GPIO23")

	ili9341Dev, err := ili9341.NewILI9341(spiConn, dataCommandSelect, reset)
	sketcher := drawings.NewSketcher(ili9341Dev)
	checkFatalErr(err)
	tests := []func(drawings.Sketcher){
		drawPoints,
		// drawCrossLines,
		// drawLines,
		// drawArc,
		// draThickwArc,
		// drawCircle,
		// drawFillCircle,
		// drawThickCircle,
		// drawRectangle,
		// drawFillRectangle,
		// drawThickRectangle,
		// drawFontsArea,
		// drawDigits,
		drawCalibrationPoints,
	}

	for i := 0; i < len(tests); i++ {
		sketcher.SetBackgroundColor(colors.WHITE)
		sketcher.Clear()
		sketcher.Update()
		ts := time.Now()
		tests[i](sketcher)
		numsegs := sketcher.Update()
		fmt.Println("Update Duration(ms): ", time.Since(ts).Milliseconds(), ", Num of updated Segments: ", numsegs)
		time.Sleep(time.Second / 10)
	}
	ili9341Dev.Update()
}

func drawPoints(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_0)
	sketcher.SetColor(colors.RED)
	sketcher.Circle(0, 0, 5)
}

func drawCrossLines(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_0)
	xmax := float64(sketcher.ScreenWidth() - 1)
	ymax := float64(sketcher.ScreenHeight() - 1)
	fmt.Println(xmax, ymax)
	sketcher.SetBackgroundColor(colors.BLACK)
	sketcher.Clear()
	const OFFSET = 5
	sketcher.SetColor(colors.RED)
	sketcher.Line(OFFSET, OFFSET, xmax-OFFSET, OFFSET)
	sketcher.SetColor(colors.GREEN)
	sketcher.Line(xmax-OFFSET, OFFSET, xmax-OFFSET, ymax-OFFSET)
	sketcher.SetColor(colors.BLUE)
	sketcher.Line(xmax-OFFSET, ymax-OFFSET, OFFSET, ymax-OFFSET)
	sketcher.SetColor(colors.PINK)
	sketcher.Line(OFFSET, ymax-OFFSET, OFFSET, OFFSET)
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

	sketcher.SetBackgroundColor(colors.WHITE)
	sketcher.Clear()
	sketcher.SetColor(colors.BLUE)
	for angle := sAngle; angle < sAngle+rAngle; angle += dAngle {
		x := math.Cos(angle) * radius
		y := math.Sin(angle) * radius
		sketcher.Line(xc, yc, xc+x, yc+y)
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
		sketcher.SetColor(colorset[i])
		sketcher.Circle(xyc[i][0], xyc[i][1], xyc[i][2])
	}
}

func drawFillCircle(sketcher drawings.Sketcher) {
	const N int = 3
	xyc := [N][]float64{{30, 30, 45}, {160, 120, 115}, {400, 400, 250}}
	colorset := [N]colors.Color{colors.BLACK, colors.DARKBLUE, colors.DARKGREEN}
	for i := 0; i < N; i++ {
		sketcher.SetColor(colorset[i])
		sketcher.FillCircle(xyc[i][0], xyc[i][1], xyc[i][2])
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
		sketcher.SetColor(colorset[i])
		sketcher.ThickCircle(xyc[i][0], xyc[i][1], xyc[i][2], width, widhTypes[i])
		sketcher.SetColor(colors.RED)
		sketcher.Circle(xyc[i][0], xyc[i][1], xyc[i][2])
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
		sketcher.SetColor(colorset[i])
		sketcher.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4])
	}
	sketcher.SetColor(colors.RED)
	sketcher.Line(160, 0, 160, 239)
	sketcher.Line(0, 120, 319, 120)
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
		sketcher.SetColor(colorset[i])
		sketcher.ThickArc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4], 10, widhTypes[i])
		sketcher.SetColor(colors.RED)
		sketcher.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4])
	}
}

func drawRectangle(sketcher drawings.Sketcher) {
	const N int = 2
	xy := [N][]float64{{10, 10, 100, 100}, {50, 50, 200, 200}}
	colorset := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		sketcher.SetColor(colorset[i])
		sketcher.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

func drawFillRectangle(sketcher drawings.Sketcher) {
	const N int = 2
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}}
	colors := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		sketcher.SetColor(colors[i])
		sketcher.FillRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

func drawThickRectangle(sketcher drawings.Sketcher) {
	const N int = 3
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}, {100, 100, 300, 220}}
	colorset := [N]colors.Color{colors.ROYALBLUE, colors.NAVY, colors.FORESTGREEN}
	widhTypes := [N]drawings.WidthType{drawings.INNER_WIDTH, drawings.CENTER_WIDTH, drawings.OUTER_WIDTH}
	const width = 10
	for i := 0; i < N; i++ {
		sketcher.SetColor(colorset[i])
		sketcher.ThickRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], width, widhTypes[i])
		sketcher.SetColor(colors.RED)
		sketcher.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

func drawFontsArea(sketcher drawings.Sketcher) {
	sketcher.SetColor(colors.BLACK)
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

		x1, y1, x2, y2 := sketcher.GetTextArea(text)
		xoffset := 8
		sketcher.SetColor(colors.RED)
		sketcher.Rectangle(float64(xoffset+x1), float64(yline+y1), float64(xoffset+x2), float64(yline+y2))
		sketcher.SetColor(colors.BLUE)
		sketcher.Line(0, float64(yline), 319, float64(yline))
		sketcher.SetColor(colors.BLACK)
		sketcher.MoveCursor(xoffset, yline)
		sketcher.Write(string(s))
		yline += 48
	}
}

func drawGrids(sketcher drawings.Sketcher) {
	for x := float64(0); x < 320; x += 32 {
		sketcher.Line(x, 0, x, 239)
	}
	for y := float64(0); y < 240; y += 24 {
		sketcher.Line(0, y, 319, y)
	}
}

func drawDigits(sketcher drawings.Sketcher) {
	sketcher.SetRotation(drawings.ROTATION_90)
	sketcher.SetFont(fonts.FreeSans24pt7b)
	X := 30
	Y := 200
	value := 23.2
	text := fmt.Sprintf("%4.1f", value)
	sketcher.MoveCursor(X, Y)
	sketcher.SetColor(colors.BLACK)
	sketcher.WriteScaled(text, 2, 5)
	// drawGrids(sketcher)
}

func drawCalibrationPoints(sketcher drawings.Sketcher) {
	const PADDING float64 = 40
	const RADIUS float64 = 5
	const N_SEGMENTS int = 2
	const X_OFFSET float64 = (320 - PADDING*2) / float64(N_SEGMENTS)
	const Y_OFFSET float64 = (240 - PADDING*2) / float64(N_SEGMENTS)
	sketcher.SetColor(colors.RED)
	for xseg := 0; xseg <= N_SEGMENTS; xseg++ {
		x := float64(xseg) * X_OFFSET
		for yseg := 0; yseg <= N_SEGMENTS; yseg++ {
			y := float64(yseg) * Y_OFFSET
			sketcher.FillCircle(x+PADDING, y+PADDING, RADIUS)
		}
	}

	sketcher.SetColor(colors.RED)
	for i := float64(0); i <= 0; i++ {
		sketcher.Line(0+i, 0+i, 319-i, 0+i)
		sketcher.Line(319-i, 0+i, 319-i, 239-i)
		sketcher.Line(319-i, 239-i, 0+i, 239-i)
		sketcher.Line(0+i, 239-i, 0+i, 0+i)
	}
}

func createGpioOutPin(gpioPinNum string) devgpio.GPIOPinOut {
	var pin gpio.PinOut = gpioreg.ByName(gpioPinNum)
	if pin == nil {
		checkFatalErr(fmt.Errorf("failed to create GPIO pin %s", gpioPinNum))
	}
	pin.Out(gpio.Low)
	return &gpioOut{
		pin: pin,
	}
}

func createSPIConnection(busNumber int, chipSelect int) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(64)*physic.MegaHertz,
		spi.Mode2,
		8,
	)
	checkFatalErr(err)
	return spiConn
}
