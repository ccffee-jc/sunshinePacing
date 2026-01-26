package main

import (
	"image"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	chartColorVideo   = color.NRGBA{R: 0x4c, G: 0xaf, B: 0x50, A: 0xff}
	chartColorGrid    = color.NRGBA{R: 0x2a, G: 0x2a, B: 0x2a, A: 0xff}
	chartColorBg      = color.NRGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xff}
)

type seriesBuffer struct {
	data []float64
	max  int
	head int
	size int
}

func newSeriesBuffer(max int) seriesBuffer {
	if max < 1 {
		max = 1
	}
	return seriesBuffer{
		data: make([]float64, max),
		max:  max,
	}
}

func (s *seriesBuffer) push(v float64) {
	if s.max == 0 {
		return
	}
	if s.size < s.max {
		s.data[s.size] = v
		s.size++
		return
	}
	s.data[s.head] = v
	s.head = (s.head + 1) % s.max
}

func (s *seriesBuffer) values() []float64 {
	if s.size == 0 {
		return nil
	}
	out := make([]float64, s.size)
	for i := 0; i < s.size; i++ {
		idx := (s.head + i) % s.max
		out[i] = s.data[idx]
	}
	return out
}

type BurstChart struct {
	widget.BaseWidget

	mu      sync.Mutex
	size    fyne.Size
	image   *image.RGBA
	video   seriesBuffer
	dirty   bool
	render  bool
}

func NewBurstChart(points int) *BurstChart {
	chart := &BurstChart{
		video:   newSeriesBuffer(points),
	}
	chart.ExtendBaseWidget(chart)
	return chart
}

func (c *BurstChart) MinSize() fyne.Size {
	return fyne.NewSize(520, 180)
}

func (c *BurstChart) Push(videoQueue int) {
	c.mu.Lock()
	c.video.push(float64(videoQueue))
	c.scheduleRenderLocked()
	c.mu.Unlock()
}

func (c *BurstChart) updateSize(size fyne.Size) {
	c.mu.Lock()
	if size == c.size {
		c.mu.Unlock()
		return
	}
	c.size = size
	c.scheduleRenderLocked()
	c.mu.Unlock()
}

func (c *BurstChart) scheduleRenderLocked() {
	c.dirty = true
	if c.render {
		return
	}
	c.render = true
	go c.renderAsync()
}

func (c *BurstChart) renderAsync() {
	for {
		c.mu.Lock()
		if !c.dirty {
			c.render = false
			c.mu.Unlock()
			return
		}
		c.dirty = false
		size := c.size
		video := c.video.values()
		c.mu.Unlock()

		img := buildChartImage(size, video)
		if img == nil {
			continue
		}
		fyne.Do(func() {
			c.mu.Lock()
			c.image = img
			c.mu.Unlock()
			c.Refresh()
		})
	}
}

func (c *BurstChart) CreateRenderer() fyne.WidgetRenderer {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	raster := canvas.NewImageFromImage(img)
	raster.FillMode = canvas.ImageFillStretch
	raster.ScaleMode = canvas.ImageScaleSmooth
	renderer := &burstChartRenderer{
		chart:   c,
		raster:  raster,
		objects: []fyne.CanvasObject{raster},
	}
	return renderer
}

type burstChartRenderer struct {
	chart   *BurstChart
	raster  *canvas.Image
	objects []fyne.CanvasObject
}

func (r *burstChartRenderer) Layout(size fyne.Size) {
	r.raster.Resize(size)
	r.chart.updateSize(size)
}

func (r *burstChartRenderer) MinSize() fyne.Size {
	return r.chart.MinSize()
}

func (r *burstChartRenderer) Refresh() {
	r.chart.mu.Lock()
	img := r.chart.image
	r.chart.mu.Unlock()
	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	r.raster.Image = img
	canvas.Refresh(r.raster)
}

func (r *burstChartRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *burstChartRenderer) Destroy() {}

func newBurstLegend() fyne.CanvasObject {
	return container.NewHBox(
		newLegendItem("video", chartColorVideo),
	)
}

func newLegendItem(label string, c color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(c)
	rect.SetMinSize(fyne.NewSize(12, 12))
	return container.NewHBox(rect, widget.NewLabel(label))
}

func buildChartImage(size fyne.Size, video []float64) *image.RGBA {
	if size.Width <= 0 || size.Height <= 0 {
		return nil
	}
	width := int(size.Width)
	height := int(size.Height)
	if width < 2 || height < 2 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fillRect(img, chartColorBg)

	padding := 8
	left := padding
	top := padding
	right := width - padding - 1
	bottom := height - padding - 1
	if right <= left || bottom <= top {
		return img
	}
	drawGrid(img, left, top, right, bottom)

	maxVal := maxSeriesValue(video)
	if maxVal < 1 {
		maxVal = 1
	}
	drawSeries(img, video, chartColorVideo, left, top, right, bottom, maxVal)
	return img
}

func maxSeriesValue(series ...[]float64) float64 {
	maxVal := float64(0)
	for _, s := range series {
		for _, v := range s {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	return maxVal
}

func fillRect(img *image.RGBA, c color.Color) {
	b := img.Bounds()
	r, g, bVal, a := c.RGBA()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(bVal >> 8),
				A: uint8(a >> 8),
			})
		}
	}
}

func drawGrid(img *image.RGBA, left, top, right, bottom int) {
	width := right - left
	height := bottom - top
	if width <= 0 || height <= 0 {
		return
	}
	for i := 0; i <= 4; i++ {
		y := top + int(float64(height)*float64(i)/4.0)
		drawLine(img, left, y, right, y, chartColorGrid)
	}
	drawLine(img, left, top, right, top, chartColorGrid)
	drawLine(img, left, bottom, right, bottom, chartColorGrid)
	drawLine(img, left, top, left, bottom, chartColorGrid)
	drawLine(img, right, top, right, bottom, chartColorGrid)
}

func drawSeries(img *image.RGBA, values []float64, c color.Color, left, top, right, bottom int, maxVal float64) {
	if len(values) == 0 {
		return
	}
	width := right - left
	height := bottom - top
	if width <= 0 || height <= 0 {
		return
	}
	points := len(values)
	if points == 1 {
		x := left
		y := bottom - int((values[0]/maxVal)*float64(height))
		drawLine(img, x, y, x, y, c)
		return
	}
	var prevX, prevY int
	for i, v := range values {
		x := left + int(float64(i)*float64(width)/float64(points-1))
		y := bottom - int((v/maxVal)*float64(height))
		if y < top {
			y = top
		}
		if y > bottom {
			y = bottom
		}
		if i > 0 {
			drawLine(img, prevX, prevY, x, y, c)
		}
		prevX = x
		prevY = y
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := absInt(x1 - x0)
	dy := -absInt(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		setPixel(img, x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func setPixel(img *image.RGBA, x, y int, c color.Color) {
	if !image.Pt(x, y).In(img.Bounds()) {
		return
	}
	img.Set(x, y, c)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
