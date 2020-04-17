package clock

import (
	"fmt"
	"io"
	"math"

	svg "github.com/ajstarks/svgo/float"
	visual "github.com/osraige/visualisations"
)

type ClockOptions struct {
	canvas             *svg.SVG
	circumOut          float64
	circumIn           float64
	radiOut            float64
	radiIn             float64
	Size               float64
	CenterRadius       float64
	HandGap            float64
	Segments           int
	Colour             string
	ColourAccent       string
	ColourAverage      string
	AverageStrokeWidth float64
	AveragePointRadius float64
	DataHands          []int
	DataAverage        []int
	Debug              bool
	Animate            bool
}

func (o ClockOptions) drawHands(group string) {
	handTop := ((o.circumOut / float64(o.Segments)) - o.HandGap) / 2.0
	handBottom := ((o.circumIn / float64(o.Segments)) - o.HandGap) / 2.0
	o.canvas.Gid(group)
	defer o.canvas.Gend()
	o.iterDataOnSeg(o.DataHands, func(height int, a float64) {
		widthSc := visual.ScaleRange(float64(height), 0, 100, handBottom, handTop)
		heightSc := visual.ScaleRange(float64(height), 0, 100, o.radiIn, o.radiOut)
		o.canvas.TranslateRotate(o.radiOut, o.radiOut, float64(a-180))
		// draw the background of the hand
		o.canvas.Polyline(
			[]float64{-handBottom, handBottom, handTop, -handTop},
			[]float64{o.radiIn, o.radiIn, o.radiOut, o.radiOut},
			visual.ParseFill(o.ColourAccent),
		)
		// draw the hand itself. height depends on the current
		// data point. the width of the side of the trapezoid
		// closest to the edge is scaled to the new range
		o.canvas.Polyline(
			[]float64{-handBottom, handBottom, widthSc, -widthSc},
			[]float64{o.radiIn, o.radiIn, heightSc, heightSc},
			visual.ParseFill(o.Colour),
		)
		o.canvas.Gend()
	})
}

func (o ClockOptions) drawHourMarkings(group string) {
	textStyle := []string{
		visual.ParseFontFamily("monospace"),
		visual.ParseFontSize(24),
		visual.ParseTextAnchor("middle"),
		visual.ParseDominantBaseline("central"),
	}
	o.canvas.Gid(group)
	defer o.canvas.Gend()
	// increment on every 45 degress for positon of marking, and every
	// 3 hours for the text itself
	for x, t := -90, 0; x < 270; x, t = x+45, t+3 {
		px, py := visual.PointOnCircum(o.radiOut, o.radiOut, o.radiIn-20, float64(x))
		opts := append([]string{}, textStyle...)
		// grey out every other marking
		if t%6 != 0 {
			opts = append(opts, visual.ParseFill("#eee"))
		}
		o.canvas.Text(px, py, fmt.Sprintf("%02d", t), visual.ParseStyles(opts...))
	}
}

func (o ClockOptions) drawAverage(group string) {
	strokeStyle := visual.ParseStyles(
		visual.ParseFill("none"),
		visual.ParseStroke(o.ColourAverage),
		visual.ParseStrokeWidth(o.AverageStrokeWidth),
	)
	xs := []float64{}
	ys := []float64{}
	o.canvas.Gid(group)
	defer o.canvas.Gend()
	o.iterDataOnSeg(o.DataAverage, func(height int, a float64) {
		heightSc := visual.ScaleRange(float64(height), 0, 100, o.radiIn, o.radiOut)
		px, py := visual.PointOnCircum(o.radiOut, o.radiOut, heightSc, float64(a-90))
		xs = append(xs, px)
		ys = append(ys, py)
		o.canvas.Circle(px, py, o.AveragePointRadius, visual.ParseFill(o.ColourAverage))
	})
	// wrap the last trend data segment to the first to connect the dots
	xs = append(xs, xs[0])
	ys = append(ys, ys[0])
	o.canvas.Polyline(xs, ys, strokeStyle)
}

func (o ClockOptions) drawDebug(group string) {
	strokeStyle := visual.ParseStyles(
		visual.ParseFill("none"),
		visual.ParseStroke("black"),
	)
	o.canvas.Gid(group)
	defer o.canvas.Gend()
	o.canvas.Line(o.radiOut, 0, o.radiOut, o.Size, strokeStyle)
	o.canvas.Line(0, o.radiOut, o.Size, o.radiOut, strokeStyle)
}

func (o ClockOptions) iterDataOnSeg(data []int, cb func(int, float64)) {
	angleInc := 360.0 / float64(o.Segments)
	for x := 0.0; x < 360.0; x += angleInc {
		var point int
		if len(data) > 0 {
			point = data[int(x/angleInc)%len(data)]
		}
		cb(point, x)
	}
}

func Clock(out io.Writer, opts ClockOptions) {
	canvas := svg.New(out)
	canvas.Start(opts.Size, opts.Size)
	defer canvas.End()
	opts.canvas = canvas
	canvas.Gid("root")
	defer canvas.Gend()
	opts.radiOut = opts.Size / 2.0
	opts.radiIn = opts.CenterRadius
	opts.circumOut = 2.0 * math.Pi * opts.radiOut
	opts.circumIn = 2.0 * math.Pi * opts.radiIn
	opts.drawHands("hands")
	opts.drawHourMarkings("hour-markings")
	if opts.Debug {
		opts.drawDebug("debug")
	}
	opts.drawAverage("average")
	if opts.Animate {
		canvas.Animate("#average", "opacity", 0, 1, 0.75, 1)
	}
}
