package gauge

import (
	"fmt"
	"io"
	"math"

	svg "github.com/ajstarks/svgo"
	visual "github.com/osraige/visualisations"
)

type GaugeOptions struct {
	canvas *svg.SVG

	// Size is determines the width and height of the gauge
	Size float64
	// Padding determines the padding around the gauge
	Padding float64
	// GapRadians sets the angle of the gap at the bottom of the gauge
	GapRadians float64
	// BackgroundColour sets the non filled portion of the gauge's colour
	BackgroundColour string
	// Colour is the colour to fill the gauge with
	Colour    string
	LineWidth float64
	// FillPorportion is the proportion of the gauge to fill between 0 and 1
	FillProportion float64

	// Label is the text to display in the center of the gauge
	Label       string
	LabelFont   string
	LabelColour string
	LabelSize   int
}

func (g *GaugeOptions) drawGauge() {
	c := g.Size / 2
	r := c - g.Padding
	//
	startAngle := math.Pi + (math.Pi-g.GapRadians)/2
	startAngle = math.Mod(startAngle, 2*math.Pi)
	startX := int(c + r*math.Cos(startAngle))
	startY := int(c - r*math.Sin(startAngle))
	//
	endAngle := startAngle + g.GapRadians
	endAngle = math.Mod(endAngle, 2*math.Pi)
	endX := int(c + r*math.Cos(endAngle))
	endY := int(c - r*math.Sin(endAngle))
	//
	midAngle := startAngle - g.FillProportion*
		(math.Pi*2-(endAngle-startAngle))
	midX := int(c + r*math.Cos(midAngle))
	midY := int(c - r*math.Sin(midAngle))
	//
	large := math.Pi < startAngle-midAngle
	g.canvas.Arc(startX, startY, int(r), int(r), 0, large, true, midX, midY,
		visual.ParseStyles(
			visual.ParseStroke(g.Colour),
			visual.ParseStrokeWidth(g.LineWidth),
			visual.ParseFill("none"),
		))
	//
	large = math.Pi > endAngle-midAngle
	g.canvas.Arc(midX, midY, int(r), int(r), 0, large, true, endX, endY,
		visual.ParseStyles(
			visual.ParseStroke(g.BackgroundColour),
			visual.ParseStrokeWidth(g.LineWidth),
			visual.ParseFill("none"),
		))

	g.canvas.Text(int(c), int(c), g.Label,
		visual.ParseStyles(
			visual.ParseFill(g.LabelColour),
			visual.ParseFontSize(g.LabelSize),
			visual.ParseDominantBaseline("central"),
			visual.ParseTextAnchor("middle"),
			visual.ParseFontFamily(g.LabelFont),
		))
}

// Gauge generates a gauge with the given options
func Gauge(out io.Writer, opts GaugeOptions) {
	canvas := svg.New(out)
	canvas.Start(int(opts.Size), int(opts.Size))
	defer canvas.End()
	opts.canvas = canvas
	canvas.Gid("root")
	opts.drawGauge()
	canvas.Gend()
}

// Gauges generates a single image with several gauges joined horizontaly
func Gauges(out io.Writer, opts []GaugeOptions) {
	var totalWidth float64
	var maxHeight float64
	for _, opt := range opts {
		totalWidth += opt.Size
		if maxHeight < opt.Size {
			maxHeight = opt.Size
		}
	}
	canvas := svg.New(out)
	canvas.Start(int(totalWidth), int(maxHeight))
	defer canvas.End()
	curWidth := 0
	canvas.Gid("root")
	for _, opt := range opts {
		opt.canvas = canvas
		canvas.Gtransform(fmt.Sprintf("translate(%v)", curWidth))
		curWidth += int(opt.Size)
		opt.drawGauge()
		canvas.Gend()
	}
	canvas.Gend()
}
